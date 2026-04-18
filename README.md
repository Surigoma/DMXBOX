# DMXBOX

[日本語](README.ja.md)

DMXBOX is a DMX lighting control system. It features a web-based frontend for configuration and control, supporting HTTP/TCP input, DMX hardware/Art-Net output, device models (dimmer, WC Light), and group management. It also supports mute processing for mixers using OSC.

## Main Features

- **Backend (Go)**:
  - DMX server (FPS control, fade effects).
  - DMX hardware (FTDI, etc.), Art-Net output support.
  - HTTP API (configuration, console, DMX control, OSC mapping).
  - TCP server (custom input).
  - OSC server (mixer control).
  - Modular architecture (enable/disable via config).

- **Frontend (React + TypeScript + Vite)**:
  - Device management: groups, models (dimmer, WC Light).
  - I/O settings: HTTP, TCP, OSC, Art-Net, DMX.
  - Controls: Mute, Fade, Error display.
  - Responsive UI, test utilities.

- **Configuration**: JSON-based (ports, devices, modules).
- **Testing**: Backend/frontend unit tests.
- **Build**: Cross-platform with Taskfile.

## Prerequisites

### Task Runner
- [Task](https://taskfile.dev/)

### Backend
- [Go 1.21+](https://go.dev/)
- [swag](https://github.com/swaggo/swag) (for API documentation)
- [Air](https://github.com/air-verse/air) (hot reload, optional)

### Frontend
- [Node.js 18+](https://nodejs.org/)
- [yarn](https://yarnpkg.com/)

## Quick Start

1. Clone the repository.
2. Run `task dev` to start development (backend + frontend hot reload).
3. Open `http://localhost:5173` (frontend).
4. Configure via UI or `backend/config.json`.

## Configuration

Main configuration: `backend/config.json` (auto-generated from defaults).

Default example:
```json
{
  "modules": {
    "http": true,
    "tcp": true
  },
  "output": {
    "target": ["console"],
    "dmx": {"port": "COM1"},
    "artnet": {"addr": "2.255.255.255/8", "universe": 0}
  },
  "http": {"ip": "127.0.0.1", "port": 8000},
  "tcp": {"ip": "127.0.0.1", "port": 50000},
  "dmx": {
    "groups": {
      "group1": {
        "name": "Group 1",
        "devices": [{"model": "dimmer", "channel": 1}]
      }
    },
    "fadeInterval": 0.7
  },
  "osc": {"ip": "127.0.0.1", "port": 8765, "format": "/yosc:req/set/MIXER:Current/InCh/Fader/On/{}/1"}
}
```

- Enable modules with `modules`.
- Define DMX groups/devices.
- Save changes via HTTP API.

## Development

- `task dev`: Air (backend) + Vite (frontend).
- `task test`: Run all tests.
- `task test_watch`: Watch mode.

## Build

Cross-platform support (Windows/Linux/macOS). Tested on Windows/Linux only.

- `task build`: Build for current platform.
- `task build_all`: Build for all targets (Win/Linux).
Outputs to `dist/`.

`task default` creates directories + builds + copies config.

## API

`http://localhost:8000`

Endpoints:
- `/api/v1/config`: Configuration GET/POST.
- `/api/v1/console`: Console output.
- `/api/v1/dmx`: DMX control.
- `/api/v1/health`: Health check.
- `/api/v1/osc`: OSC mapping.

For API details, see `http://localhost:8000/docs/index.html`.

## Screenshots

### Control Screen
![Control](images/control.png)

### Configuration Screen
![Configuration](images/settings.png)

## Testing

- Backend: `cd backend; task test`
- Frontend: `cd frontend; task test`
- All: `task test`

## Architecture

- **main.go**: Module management (DMX, HTTP, TCP, OSC).
- Message passing via channels.
- Graceful shutdown on SIGINT.

## Supported Devices

- `dimmer`: Simple dimming control (single channel).
- `wclight`: WC Light (white color adjustment lighting).
  - Channel configuration: [cool (cool white), warm (warm white), flash (unused: always 0)]

## Troubleshooting

- Configuration load failure: Check JSON syntax (`backend/config.go` Load function).
- DMX hardware: Check ports (Windows: COM1, etc.; Linux: /dev/ttyUSB0, etc.).
- Port conflicts: Change IP/Port in config.

Check logs (JSON format, using slog).

## License
See [LICENSE](LICENSE).