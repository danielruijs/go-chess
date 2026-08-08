# Backend Architecture & Conceptual Flows

This document explains the architecture of the go-chess backend, defines its core concepts, and visualizes the lifecycles of connections and games.

## Core Concepts

The backend is built around a few central entities that manage user identity, network connections, and gameplay state:

*   **User (`auth.User`)**: A persistent user identity registered in the system (stored in the PostgreSQL database). Contains credentials (hashed via bcrypt) and a display name. Stored and queried via the `UserStore`.
*   **Session (`auth.Session`)**: A temporary session representation. Each session can be **anonymous** (temp guest name, not registered) or **authenticated** (tied to a `User`). Sessions are identified by a `SessionID` (UUID) stored in the client's secure HTTP-only `session_id` cookie. Valid sessions are cached for 24 hours in the `SessionStore`.
*   **PlayerKey (`server.PlayerKey`)**: A unique string key used to bridge the authentication layer (`auth`) and the gameplay layer (`server`). 
    *   For authenticated users, the format is `"user:<username>"`.
    *   For anonymous users, the format is `"anon:<sessionID>"`.
*   **Player (`server.Player`)**: The in-memory identity of an active player on the server. Players are cached in `WebSocketHandler.players`. They track matchmaking queues they have joined and matches they are currently playing.
*   **Client (`server.Client`)**: Represents a single active WebSocket connection (`*websocket.Conn`). A `Client` runs concurrent read/write loops. Every `Client` belongs to exactly one `Player`, but a `Player` can have multiple `Client` connections.
*   **Matchmaker (`server.Matchmaker`)**: A central component operating on an actor-like model (running in a single goroutine). It maintains queues for different time formats and pairs matching players.
*   **Match (`server.Match`)**: An active game session between two players. It runs in its own goroutine and communicates with clients via thread-safe, sequential event channels (`EventChan`), shielding game state mutations from concurrent access.
*   **Engine (`chess.Engine`)**: A pure, dependency-free domain model that represents the chess game state (board, castling rights, active color, move history) and handles move verification/legal move generation.

## Concept Relationships & Multiplicities

A player can open multiple tabs in their browser, which means one `Player` can map to multiple WebSocket `Client` connections simultaneously. Similarly, a registered `User` can have multiple active sessions on different devices.

```mermaid
classDiagram
    direction LR
    class User {
        Username string
        DisplayName string
    }
    class Session {
        ID SessionID
        Username string
    }
    class Player {
        Key PlayerKey
        Username string
        DisplayName string
    }
    class Client {
        Conn *websocket.Conn
    }
    class Match {
        White *Player
        Black *Player
    }
    User "1" --> "0..*" Session : owns
    Session "1" --> "1" PlayerKey : maps to
    PlayerKey "1" --> "0..1" Player : identifies (cached)
    Player "1" --> "0..*" Client : connected via
    Player "2" --> "0..1" Match : plays in
```

## High-Level System Architecture

The codebase separates concerns into isolated packages. `chess` and `cache` are generic, utility-level packages with no internal dependencies, while `server` orchestrates them using sessions from `auth`.

```mermaid
graph TD
    subgraph Main [main]
        main_go["main.go (Bootstrap)"]
        config_go["config.go (Flags & Config)"]
    end

    subgraph Auth [internal/auth]
        AuthHandler["AuthHandler (HTTP API)"]
        SessionStore["SessionStore"]
        UserStore["UserStore (Bcrypt Auth)"]
    end

    subgraph DB [internal/db]
        DBQueries["Queries (sqlc generated)"]
        DBConnection["Pool (pgxpool)"]
    end

    subgraph Server [internal/server]
        WebSocketHandler["WebSocketHandler (/ws)"]
        Matchmaker["Matchmaker (Queue & Pair)"]
        Match["Match (Game Loop)"]
    end

    subgraph Chess [internal/chess]
        Engine["Engine (Rule validator)"]
        Position["Position (Bitboards)"]
    end

    subgraph Cache [internal/cache]
        CacheWithCleanup["Cache (With TTL Eviction)"]
        CacheWithoutCleanup["Cache (Persistent / No TTL Eviction)"]
    end

    main_go --> AuthHandler
    main_go --> WebSocketHandler
    main_go --> Matchmaker
    main_go --> DBConnection

    AuthHandler --> UserStore
    AuthHandler --> SessionStore
    SessionStore -->|Session| CacheWithCleanup
    UserStore -->|Queries| DBQueries

    WebSocketHandler --> SessionStore
    WebSocketHandler --> Matchmaker
    WebSocketHandler -->|Player| CacheWithCleanup

    Matchmaker --> Match
    Matchmaker -->|Match| CacheWithoutCleanup
    Match --> Engine
```

## User Connection & Registration Flow

This flowchart outlines how anonymous visitors receive temporary sessions, connect via WebSockets, register/login to upgrade their session, and reconnect under their authenticated identity.

```mermaid
sequenceDiagram
    autonumber
    actor Browser as Client Browser
    participant Auth as auth.AuthHandler
    participant Sessions as auth.SessionStore
    participant WS as server.WebSocketHandler
    
    Note over Browser, Sessions: Phase 1: Anonymous Initialization
    Browser->>Auth: GET /api/auth/check (No Cookie)
    Auth->>Sessions: CreateAnonSession()
    Sessions-->>Auth: Session (ID, temp DisplayName)
    Auth-->>Browser: 200 OK + Set-Cookie: session_id=<UUID>
    
    Note over Browser, WS: Phase 2: WebSocket Connection
    Browser->>WS: Connect to /ws (Sends session_id Cookie)
    WS->>Sessions: Get Session by ID
    Sessions-->>WS: Session info
    WS->>WS: Get or Create Player (anon:<UUID>)
    WS->>WS: Register Client to Player
    WS-->>Browser: player_info message + matchmaking_update
    
    Note over Browser, Auth: Phase 3: Registration / Upgrade
    Browser->>Auth: POST /api/auth/register (username, password, displayName)
    Auth->>Sessions: CreateSession(username, displayName)
    Sessions-->>Auth: Authenticated Session
    Auth-->>Browser: 200 OK + Set-Cookie: session_id=<UUID> (new session)
    
    Note over Browser, WS: Phase 4: WS Reconnection
    Browser->>WS: Reconnect to /ws (New Cookie)
    WS->>Sessions: Get Session by ID
    Sessions-->>WS: Session info (authenticated)
    WS->>WS: Get or Create Player (user:username)
    WS->>WS: Register Client to Player
    WS-->>Browser: player_info (authenticated)
```

## Matchmaking & Gameplay Lifecycle

Games are created by the `Matchmaker` and managed in separate goroutines by the `Match` instance. Client interaction occurs via `MessageHandlers` which dispatch to a serialized internal `EventChan` to prevent concurrent modification of the active game state.

```mermaid
sequenceDiagram
    autonumber
    participant C1 as Client 1 (Player 1)
    participant C2 as Client 2 (Player 2)
    participant MM as server.Matchmaker
    participant M as server.Match
    
    C1->>MM: Join matchmaking queue (TimeFormat)
    Note over MM: Player 1 added to queue
    
    C2->>MM: Join matchmaking queue (TimeFormat)
    Note over MM: Player 2 added to queue
    
    Note over MM: Matchmaker pairs Client 1 and Client 2
    MM->>M: Create Match & Start match.Run() goroutine
    MM->>MM: Register Match (adds to activeMatches)
    
    M-->>C1: WS Message: start_match (Color: White)
    M-->>C2: WS Message: start_match (Color: Black)
    
    Note over C1, M: Game Loop
    C1->>M: WS Message: move (e2 -> e4, via EventChan)
    Note over M: Match.Run() processes event sequentially
    M->>M: Validate move & update clock
    M-->>C1: WS Message: board (Updated state, legal moves, clocks)
    M-->>C2: WS Message: board (Updated state, legal moves, clocks)
    
    Note over M: Checkmate / Resign / Timeout / Draw
    M-->>C1: WS Message: end_match (Result)
    M-->>C2: WS Message: end_match (Result)
    M->>MM: Notify Match Ended
    MM->>MM: Unregister Match (removes from activeMatches)
    Note over M: Goroutine exits
```
