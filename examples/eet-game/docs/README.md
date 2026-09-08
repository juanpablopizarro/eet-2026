# 👾 EET Space Invaders

A browser-based Space Invaders clone with a Go WebSocket backend and a neon pixel-art canvas frontend.

## Project Layout

```
eet-game/
├── cmd/server/         # Entry point
├── internal/
│   ├── game/           # Core game loop, entities, collision, levels
│   ├── server/         # HTTP + WebSocket server
│   └── protocol/       # JSON message types
├── pkg/vec2/           # 2D vector math
├── web/static/         # HTML5 Canvas UI (index.html, game.js, style.css)
├── configs/            # Configuration files
└── docs/               # This file
```

## Running

```bash
# Install dependencies
go mod tidy

# Run the server (development)
make run

# Open in browser
open http://localhost:8080
```

## Controls

| Key | Action |
|-----|--------|
| `←` / `A` | Move left |
| `→` / `D` | Move right |
| `Space` | Fire |
| `P` | Pause / Resume |
| `Enter` | Start / Restart |

## Architecture

- The **Go server** runs the authoritative game loop at 60fps
- The **browser client** sends input events and renders the received game state
- Communication uses **WebSockets** via `github.com/coder/websocket`
- Each browser connection gets its own isolated game instance

## Scoring

| Target | Points |
|--------|--------|
| Bottom aliens (rows 1-2) | 10 |
| Middle aliens (rows 3-4) | 20 |
| Top alien (row 5) | 30 |
| UFO 🛸 | 50–300 (random) |
