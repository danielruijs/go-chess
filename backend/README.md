# Backend

The backend uses Go.

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
