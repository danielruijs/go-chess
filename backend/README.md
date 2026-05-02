# Backend

The backend uses Go.

## Running the backend

To run the backend, use the following command:

```bash
task run
```

The backend serves gameplay traffic on `localhost:8085`.

### Metrics

Prometheus metrics are served on `localhost:2115/metrics`. Backend metrics include websocket connection counts, matchmaking queue depth, active matches, match durations, move and draw event totals, clock timeouts, and match results by outcome and reason.

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
