# Backend

The backend uses Go.

## Architecture & Conceptual Flow

For a comprehensive overview of the backend's architecture, including concept explanations, and data flow diagrams, see the [Architecture Documentation](ARCHITECTURE.md).

## Running the backend

To run the backend, use the following command:

```bash
task run
```

The backend serves gameplay traffic on `localhost:8085`.

### Monitoring

For monitoring the backend, see the [Monitoring](monitoring/README.md) documentation.

## Deploying in a container

Use the docker compose file in the project root to run the backend in a container by running the following command in the project root:

```bash
docker compose up -d backend
```

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

- **URL**: `/api/check`
- **Method**: `GET`
- **Response Headers**: `Set-Cookie: session_id=<sessionID>` (if cookie is set/renewed)
- **Response Content-Type**: `application/json`
- **Response Body**:
  ```json
  {
    "username": "",
    "displayName": "Anonymous 1",
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

- **URL**: `/api/register`
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

- **URL**: `/api/login`
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

- **URL**: `/api/logout`
- **Method**: `POST`
- **Response Headers**: `Set-Cookie: session_id=; Max-Age=0` (expires in the past)
- **Response Body**: _Empty (200 OK)_

---

### 5. WebSocket Gameplay Connection

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
      "displayName": "Anonymous 1",
      "isAuthenticated": false
    }
  }
  ```
