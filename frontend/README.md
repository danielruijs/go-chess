# Frontend

The frontend uses React + TypeScript + Vite.

## Setup

Create a `.env` file with the following:

```env
VITE_WS_URL="ws://localhost:8085/ws"
VITE_API_URL="http://localhost:8085"
```

In production, `VITE_WS_URL` and `VITE_API_URL` should be set to the production WebSocket and API endpoints respectively.

## Running the frontend

To run the frontend, use the following command:

```bash
task run
```

## Linting and running tests

To lint and run tests for the frontend, use the following command:

```bash
task ci
```
