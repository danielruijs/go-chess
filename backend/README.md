# Backend

The backend uses Go.

## Architecture & Conceptual Flow

For a comprehensive overview of the backend's architecture, including concept explanations, and data flow diagrams, see the [Architecture Documentation](ARCHITECTURE.md).

## Running the backend

To run the backend locally:

1. **Start the database and run migrations:**
    ```bash
    task db
    ```
    This starts the PostgreSQL container, exposes port `5433` to the host, and runs database migrations. Use `task db-down` to stop the database.

2. **Run and iterate on the backend Go application:**
   ```bash
   task run
   ```
   This compiles and runs the Go backend locally on your host machine, connecting to the running PostgreSQL container.

The backend serves gameplay traffic on `localhost:8085`.

### Monitoring

For monitoring the backend, see the [Monitoring](monitoring/README.md) documentation.

## Deploying in a container

To run the full stack in a containerized environment, execute the following from the root directory:

```bash
docker compose up -d
```

This will start PostgreSQL, run database schema migrations, and deploy the backend service.

## Linting and running tests

To lint and run tests for the backend, use the following command:

```bash
task ci
```

## API Endpoints

The backend exposes the following REST endpoints and WebSocket channels on port `8085`.

### Session Cookie Details

All authentication endpoints utilize a secure session cookie named `session_id`.

- **Configuration**: Domain configured with the required `-cookie-domain` flag, `Secure=true` (supported on `localhost`/`127.0.0.1` by modern browsers over HTTP), `HttpOnly=true`, `SameSite=Lax`.
- **CORS Policy**: Requests with credentials enabled (`{ credentials: "include" }`) are accepted from origins specified in the `-allowed-origins` command-line flag.

---

### 1. Check Session / Initialize

Check if the user has an active session. If the session cookie is missing or invalid, the server automatically initializes a new anonymous session and returns it in the response cookie.

- **URL**: `/api/auth/check`
- **Method**: `GET`
- **Response Headers**: `Set-Cookie: session_id=<sessionID>` (if cookie is set/renewed)
- **Response Content-Type**: `application/json`
- **Response Body**:
  ```json
  {
    "username": "",
    "displayName": "Anonymous",
    "isAuthenticated": false
  }
  ```
  _Or if logged in:_
  ```json
  {
    "username": "alice",
    "displayName": "Alice Smith 👑",
    "isAuthenticated": true
  }
  ```

---

### 2. User Registration

Registers a new user account and upgrades the active session.

- **URL**: `/api/auth/register`
- **Method**: `POST`
- **Request Content-Type**: `application/json`
- **Request Body**:
  ```json
  {
    "username": "alice",
    "password": "SuperSecretPassword1",
    "displayName": "Alice Smith 👑"
  }
  ```
- **Response Headers**: `Set-Cookie: session_id=<sessionID>`
- **Response Content-Type**: `application/json`
- **Response Body (Success 200 OK)**:
  ```json
  {
    "username": "alice",
    "displayName": "Alice Smith 👑",
    "isAuthenticated": true
  }
  ```
- **Error Responses**:
  - `400 Bad Request`: Username already taken, invalid username, display name or password.

---

### 3. User Login

Authenticates user credentials and upgrades/establishes the session.

- **URL**: `/api/auth/login`
- **Method**: `POST`
- **Request Content-Type**: `application/json`
- **Request Body**:
  ```json
  {
    "username": "alice",
    "password": "SuperSecretPassword1"
  }
  ```
- **Response Headers**: `Set-Cookie: session_id=<sessionID>`
- **Response Content-Type**: `application/json`
- **Response Body (Success 200 OK)**:
  ```json
  {
    "username": "alice",
    "displayName": "Alice Smith 👑",
    "isAuthenticated": true
  }
  ```
- **Error Responses**:
  - `401 Unauthorized`: Invalid username or password.

---

### 4. User Logout

Logs out the user, invalidates the session record on the backend, and clears the client's session cookie.

- **URL**: `/api/auth/logout`
- **Method**: `POST`
- **Response Headers**: `Set-Cookie: session_id=; Max-Age=0` (expires in the past)
- **Response Body**: _Empty (200 OK)_

---

### 5. Get Match

Retrieves the recorded move history, board positions, move timestamps, and final result of an ended match by its public ID for post-game analysis and replay.

- **URL**: `/api/match/{publicId}`
- **Method**: `GET`
- **Response Content-Type**: `application/json`
- **Response Body (Success 200 OK)**:
  ```json
  {
    "whitePlayerName": "Alice Smith 👑",
    "blackPlayerName": "Bob",
    "result": {
      "outcome": "white_win",
      "reason": "checkmate"
    },
    "positions": [
      {
        "index": 0,
        "board": [
          ["r", "n", "b", "q", "k", "b", "n", "r"],
          ["p", "p", "p", "p", "p", "p", "p", "p"],
          [" ", " ", " ", " ", " ", " ", " ", " "],
          [" ", " ", " ", " ", " ", " ", " ", " "],
          [" ", " ", " ", " ", " ", " ", " ", " "],
          [" ", " ", " ", " ", " ", " ", " ", " "],
          ["P", "P", "P", "P", "P", "P", "P", "P"],
          ["R", "N", "B", "Q", "K", "B", "N", "R"]
        ],
        "whiteTimeMs": 300000,
        "blackTimeMs": 300000
      },
      {
        "index": 1,
        "board": [ ... ],
        "san": "e4",
        "whiteTimeMs": 298500,
        "blackTimeMs": 300000
      }
    ]
  }
  ```
- **Error Responses**:
  - `400 Bad Request`: Invalid match public ID format.
  - `404 Not Found`: Match not found.
  - `409 Conflict`: Match is currently in progress.

---

### 6. Get User Profile & Match History

Retrieves public user profile information, performance statistics, and match history for the specified user.

- **URL**: `/api/users/{username}`
- **Method**: `GET`
- **Response Content-Type**: `application/json`
- **Response Body (Success 200 OK)**:
  ```json
  {
    "user": {
      "username": "alice",
      "displayName": "Alice Smith 👑",
      "createdAt": "2024-01-15T12:00:00Z"
    },
    "stats": {
      "white": {
        "wins": 10,
        "losses": 2,
        "draws": 1
      },
      "black": {
        "wins": 8,
        "losses": 4,
        "draws": 0
      }
    },
    "matches": [
      {
        "publicId": "7b6e92f1-9d2a-4c28-98e3-4f91d08e1a12",
        "playedColor": "white",
        "opponentDisplayName": "Bob",
        "opponentUsername": "bob",
        "result": {
          "outcome": "white_win",
          "reason": "checkmate"
        },
        "timeFormat": {
          "initialMs": 300000,
          "incrementMs": 0
        },
        "moveCount": 35,
        "createdAt": "2024-01-16T15:30:00Z"
      }
    ]
  }
  ```
- **Error Responses**:
  - `400 Bad Request`: Missing username in URL path.
  - `404 Not Found`: User not found.

---

### 7. WebSocket Gameplay Connection

Establishes a persistent bi-directional connection for gameplay and matchmaking updates.

- **URL**: `/ws`
- **Protocol**: `ws://` / `wss://`
- **Authentication**: Authenticates using the `session_id` cookie passed automatically by the browser in the HTTP upgrade handshake request header.
- **Initial Connection Handshake**:
  Immediately upon connection, the server sends a `player_info` message indicating the player's server-assigned username, display name, and authentication status:
  ```json
  {
    "type": "player_info",
    "data": {
      "username": "",
      "displayName": "Anonymous",
      "isAuthenticated": false
    }
  }
  ```
