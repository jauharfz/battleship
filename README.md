# Battleship

A real-time multiplayer Battleship web game built in Go and TypeScript. Players arrange a fleet of 5 ships, enter a matchmaking queue (random matchmaking or invite a friend via a room link), and battle over a server-authoritative WebSocket connection.

I built this project to understand how stateful, concurrent game servers work in practice. The backend is written in Go without a web framework, combining REST for authentication and lobby matchmaking with an in-memory actor-based room loop for the real-time WebSocket match. The frontend is vanilla TypeScript and CSS Grid, compiled with Vite, using browser Pointer Events for client-side fleet placement.

## Tech stack

- **Backend**: Go (standard library `net/http`, `sync`, `context`, `log/slog`)
- **Networking**: `github.com/gorilla/websocket` for WebSocket RFC 6455 connections
- **Auth & Tokens**: `github.com/golang-jwt/jwt/v5` (HS256 JWT for REST auth, one-time UUID tickets for WebSocket handshake)
- **Frontend**: Vanilla TypeScript, CSS Grid, HTML5 Pointer Events
- **Build & Dev**: Vite (proxying `/api` and `/ws` to the Go backend)

## How it works

The system coordinates game flow across several communication phases:

1. **Authentication & Pre-match**: The client registers as a guest (`POST /api/v1/auth/guest`) to get a signed JWT. Before queuing, the client validates fleet coordinates locally and can optionally verify them via `QUERY /api/v1/validate`.
2. **Matchmaking**: The client requests matchmaking (`POST /api/v1/matchmaking`). For random matching, a FIFO queue pairs two waiting players; for friend matching, the first player registers a room ID and the second joins with the same ID. Once paired, the server generates a unique `room_id` and issues a one-time ticket for each player.
3. **Ticket-based WebSocket Upgrade**: Clients poll `GET /api/v1/matchmaking/status`. When matched, they receive a short-lived ticket and connect to `/ws?ticket=<ticket_id>`. The server atomically consumes the ticket during upgrade, preventing long-lived JWTs in WebSocket query parameters and rejecting unauthorized connections with HTTP 401.
4. **Room Actor Loop**: Each active match runs in its own isolated goroutine (`Room.Run()`), receiving incoming client messages via an ingress channel and managing state transitions without global locks.

## Architecture and design notes

- **Server-authoritative with hidden state**: Ship coordinates exist exclusively on the server (`game.Match`). Clients never receive their opponent's ship layout—only shot coordinates and whether it resulted in a hit, miss, or sunken ship.
- **Single-Writer rule per connection**: Each WebSocket client runs dedicated `ReadPump` and `WritePump` goroutines. Messages to a client are pushed into a buffered channel (`c.Send`), ensuring concurrent writes to a single connection never occur.
- **One-time ticket handshake**: WebSocket handshakes in browsers cannot easily pass custom authorization headers. Rather than passing long-lived JWTs in query strings (which leak into server access logs and proxy histories), the server exchanges the JWT for a single-use UUID ticket with a 6-second TTL that is deleted immediately upon consumption (`TicketStore.Consume`).
- **Bounded timeouts on match start**: A 6-second `time.Timer` guards the room joining phase. If a matched player fails to establish their WebSocket connection, the room aborts cleanly (`EVENT_GAME_OVER` with `OPPONENT_DISCONNECTED_TIMEOUT`) and drains remaining clients to prevent orphaned goroutines.
- **Pointer Events for fleet placement**: Ship dragging and rotation are implemented using native W3C Pointer Events (`pointerdown`, `pointermove`, `pointerup`) with `setPointerCapture`. The client maps screen pixels to grid cells using pure geometry, avoiding the mobile limitations and drag-image quirks of the HTML5 Drag and Drop API.
- **Request hardening**: REST handlers enforce payload limits with `http.MaxBytesReader` and use `json.NewDecoder` with `DisallowUnknownFields` to reject unexpected parameters.
- **Graceful shutdown**: `signal.NotifyContext` intercepts SIGINT/SIGTERM. Active background loops (`Matching`, `TicketSweeper`) terminate cleanly via context cancellation, while the HTTP server drains connections under a 5-second shutdown context.

## Game rules and protocol

The game follows classic Battleship rules on a 10×10 grid (0-indexed, `0` to `9`):
- Fleet composed of exactly 5 ships: Aircraft Carrier (5), Battleship (4), Heavy Cruiser (3), Light Cruiser (3), and Destroyer (2).
- Ships cannot overlap or extend beyond the 10×10 boundary.
- Alternating turns with one shot per turn. Shooting an already targeted cell returns an `ALREADY_SHOT` error.
- Victory occurs when all 5 ships of the opposing player are destroyed. Disconnecting mid-battle triggers an automatic forfeit.

The full REST API and WebSocket wire protocol contract is documented in [docs/protocol.md](docs/protocol.md).

## Project structure

```text
├── cmd/
│   └── server/          # Application entrypoint and server bootstrap
├── internal/
│   ├── auth/            # JWT signing/verification and one-time ticket store
│   ├── game/            # Pure game engine (board, ships, turns, win conditions)
│   ├── lobby/           # Matchmaking engine (random FIFO queue and custom rooms)
│   ├── protocol/        # DTO definitions, WebSocket commands, events, and mappers
│   ├── room/            # Actor loop, RoomManager, and Client read/write pumps
│   └── server/          # HTTP routes, middlewares, REST handlers, and WS upgrade
├── web/
│   ├── public/          # Static assets
│   ├── src/             # Frontend TypeScript modules (app, grid, ship, websocket)
│   ├── index.html       # Game layout and board markup
│   └── vite.config.ts   # Vite dev server and backend proxy configuration
└── docs/
    └── protocol.md      # REST & WebSocket protocol specification
```

## Running locally

### Prerequisites

- Go 1.22+
- Node.js 18+ (for building the frontend)

### 1. Build the Frontend

```bash
cd web
npm install
npm run build
cd ..
```

This compiles TypeScript and outputs the production bundle to `web/dist/`.

### 2. Run the Server

```bash
go run cmd/server/main.go
```

The server starts on `http://localhost:8080`.

Open two separate browser windows (or tabs) at `http://localhost:8080`:
1. Arrange ships on the left board (drag to reposition, click to rotate).
2. Click **Match** on both windows.
3. The server pairs both clients, transitions into battle, and announces whose turn it is.
4. Take turns clicking cells on the opponent grid on the right.

### Running with Frontend Hot Reload (Development)

To work on frontend code with Vite's dev server:

```bash
# Terminal 1: Start Go backend
go run cmd/server/main.go

# Terminal 2: Start Vite dev server
cd web
npm run dev
```

Vite proxies `/api` and `/ws` to `localhost:8080`. Open `http://localhost:5173` in your browser.

## Running tests

The pure game logic (`internal/game`) contains table-driven unit tests verifying ship validation, boundary checks, hit detection, and match state transitions:

```bash
go test -v ./...
```

Run tests with the Go race detector enabled to verify concurrency safety:

```bash
go test -race ./...
```

## Limitations and next steps

- **In-memory state**: Rooms and active matchmaking queues are stored in memory. Restarting the server terminates active matches and resets queues. Adding Redis would allow distributed state and horizontal scaling.
- **No reconnect / resume**: If a player temporarily loses network connectivity, the match immediately registers a forfeit. A production implementation would allow a grace period with state synchronization (`EVENT_RECONNECT_SYNC`).
- **Fixed guest authentication**: The current system uses transient guest identities generated on the fly. Persistent user accounts with database storage (e.g. PostgreSQL) and match history are natural next steps.
- **Deployment automation**: Adding a multi-stage Dockerfile and a Caddy reverse proxy configuration for TLS termination on a public host.
