## Code Review Instructions

These instructions guide automated code review across this repository. Apply judgment based on file type and location — not every rule fires on every file.

### General Principles
- Flag any added code not strictly necessary for the stated functionality.
- Favor the simplest correct solution; complexity must be justified by the problem.
- Avoid new abstractions or helpers unless they clearly reduce duplication in this change.
- Call out speculative generality: config options, flags, or extension points with no current caller.
- Prefer editing/deleting over adding when a simpler fix exists.

### Security Critical
- Check for hardcoded secrets, API keys, credentials, or tokens in code, config, or test fixtures.
- SQL injection: raw string-built queries anywhere user input flows in; verify parameterized queries / prepared statements are used consistently.
- Input validation at every trust boundary: HTTP handlers, WebSocket message handlers, and any place external data enters the system — not just at the outermost layer.
- AuthN/authZ: confirm a user can only act on their own sessions/games; check for missing ownership checks on game IDs, move endpoints, or admin routes (e.g. a player submitting moves for a game/color they're not part of).
- Session and cookie handling: expiry, invalidation on logout, secure/httpOnly flags.
- Rate limiting / abuse potential on endpoints that are cheap to spam.

### Concurrency & Data Integrity
- Race conditions around shared mutable state: caches, in-memory game session stores, tickers/timers, goroutines/async tasks touching the same struct or map without synchronization.
- Correct use of transactions for multi-step DB writes.
- Resource cleanup: goroutines, timers, WebSocket connections, DB connections/rows — confirm they're closed/cancelled on all exit paths, including error paths.
- Idempotency on operations that could be retried or double-submitted (reconnect after dropped WebSocket, duplicate move submission).

### Performance
- Missing indexes for columns used in frequent WHERE/JOIN/ORDER BY.
- Spot inefficient loops and algorithmic issues
- Check for memory leaks and resource cleanup
- Review caching opportunities for expensive operations
- Unnecessary re-renders or re-computation in the frontend (missing memoization only where profiling/logic actually shows it matters — don't flag its absence as a default).

### API & Contract Consistency
- Frontend/backend type mismatches: request/response shapes, enums that must stay in sync across the boundary.
- Breaking changes to existing API responses without versioning or frontend updates in the same change.
- WebSocket message schemas: check both send and receive sides are updated together.

### Code Quality
- Functions focused and reasonably sized (~50 lines as a guideline, not a hard rule).
- Clear, descriptive naming; avoid abbreviations that aren't already idiomatic in the language.
- Errors are handled or explicitly and deliberately ignored (with reasoning) — never silently swallowed.
- No dead code, unused imports, or commented-out blocks left in.
- No unnecessary styling/config bloat: Flag classes, style props, or config options that don't change behavior or duplicate something already available (e.g. redundant Tailwind classes).
- No duplicated logic: Flag logic, blocks, or queries copy-pasted instead of extracted or reused (e.g. near-identical functions or components) — not just exact duplicates, but close enough that it should be shared.
