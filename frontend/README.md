# Frontend

The frontend uses React + TypeScript + Vite.

## Setup

Create a `.env` file with the following:

```
VITE_WS_URL="ws://localhost:8085/ws"
```

In production the `VITE_WS_URL` should be set to the production websocket endpoint.

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
