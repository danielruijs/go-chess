# Copilot instructions for go-chess

## Repository summary

- Go + React monorepo for a real-time chess app. The Go backend serves WebSocket gameplay traffic and Prometheus metrics. The React frontend renders the game UI and connects to the backend WebSocket.
- Repo with two main projects: backend/ (Go) and frontend/ (React + TypeScript + Vite).

## Tech stack and versions (verified locally on 2026-05-21)

- Go 1.26 (go.mod and backend CI use 1.26). Local tool check: go1.26.0.
- Node.js 24 (frontend CI uses Node 24). Local tool check: v24.12.0 and npm 11.7.0.
- Task runner: go-task. Local tool check: task 3.44.1.
- Frontend: React 19, TypeScript 6, Vite 8, Tailwind v4, Jest, ESLint, Prettier.
- Backend linting: golangci-lint (required for backend Taskfile lint step).

## Project layout (where to look)

- backend/: Go server.
  - main entry: backend/main.go (starts WebSocket server and metrics).
  - flags: backend/config.go.
  - chess engine/core: backend/internal/chess/.
  - server runtime: backend/internal/server/ (matchmaking, websocket, metrics).
  - monitoring dashboard: backend/monitoring/grafana-dashboard.json.
  - Dockerfile: backend/Dockerfile (multi-stage build).
- frontend/: React app.
  - main entry: frontend/src/main.tsx and frontend/src/App.tsx.
  - websocket context: frontend/src/contexts/WebSocketContext.tsx.
  - chess UI components: frontend/src/components/.
  - pages: frontend/src/pages/ (Home, Game).
  - tests: frontend/src/tests/.
  - config: frontend/vite.config.ts, jest.config.ts, eslint.config.js, tsconfig\*.json.
  - env template: frontend/.env.example (VITE_WS_URL).
- Root:
  - Taskfile.yml (wraps backend/frontend CI).
  - docker-compose.yml (backend only).
  - .github/workflows/backend.yml and frontend.yml (CI and deploy steps).

## Build, test, lint, run (follow these sequences)

### Backend (Go)

- Always run: `cd backend && task ci` (runs all backend CI steps: deps, build, test, vet, lint).

### Frontend (React + Vite)

- Always run: `cd frontend && task ci` (runs all frontend CI steps: deps, build, test, lint, format).

## Root file list (for quick navigation)

- README.md
- Taskfile.yml
- docker-compose.yml
- renovate.json
- backend/
- frontend/
- .github/
- .git/

## Top-level directory contents

- backend/: README.md, Taskfile.yml, Dockerfile, go.mod, go.sum, main.go, config.go, internal/, monitoring/.
- frontend/: README.md, Taskfile.yml, package.json, package-lock.json, vite.config.ts, jest.config.ts, eslint.config.js, tsconfig\*.json, .env.example, public/, src/.
- .github/workflows/: backend.yml, frontend.yml.

## Key entrypoints (mental map)

- Backend server entrypoint: backend/main.go (websocket and metrics servers).
- Frontend entrypoint: frontend/src/main.tsx (mounts App).

## Code review instructions

These instructions guide Copilot code review across all files in this repository.
General principles that apply to all reviews:

- Flag any added code that is not strictly necessary to achieve the desired functionality.
- Keep complexity as low as possible; favor the simplest correct solution.
- Avoid new abstractions or helpers unless they clearly reduce duplication in this change.
- Call out over-engineering or speculative changes.

### Security Critical Issues

- Check for hardcoded secrets, API keys, or credentials
- Look for SQL injection and XSS vulnerabilities
- Verify proper input validation and sanitization
- Review authentication and authorization logic

### Performance Red Flags

- Identify N+1 database query problems
- Spot inefficient loops and algorithmic issues
- Check for memory leaks and resource cleanup
- Review caching opportunities for expensive operations

### Code Quality Essentials

- Functions should be focused and appropriately sized (under 50 lines)
- Use clear, descriptive naming conventions
- Ensure proper error handling throughout
- Remove dead code and unused imports

## Trust these instructions

- Follow this document first; only search the repo if you need details not covered here or if a command fails unexpectedly.
