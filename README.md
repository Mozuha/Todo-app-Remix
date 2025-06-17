# Todo app using Remix, Gin, PostgreSQL

## Tech stack
- Remix (TypeScript)
- Gin (Go)
- sqlc, pgx, golang-migrate
- PostgreSQL
- Docker

## Auth logic
**[Problem]**  
Normally, one cannot actively invalidate the JWT token and thus need to wait for its expiration. Until its expiration, anyone who have access to the valid token can be authenticated even if the original user is logged out.

**[Workaround]**  
Include the user's session ID into JWT token so the token can be invalidated when the user logged out by invalidating the user session.

## Note
Frontend (WIP)
