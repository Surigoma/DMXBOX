# DMXBOX Configuration Reference

[日本語](../ja/configuration.md) | [Documentation index](README.md)

DMXBOX reads `config.json` from the backend working directory. In a source checkout, the development file is `backend/config.json`. Saving from the browser replaces the complete configuration and reloads backend modules.

Back up a working configuration before editing it manually.

## Complete Example

```json
{
    "output": {
        "target": ["console"],
        "ftdi": { "port": "COM3" },
        "artnet": {
            "addr": "127.0.0.1",
            "universe": 0,
            "subuni": 0,
            "net": 0
        },
        "osc": {
            "ip": "127.0.0.1",
            "port": 8765,
            "format": "/yosc:req/set/MIXER:Current/InCh/Fader/On/{}/1",
            "type": "int",
            "inverse": true,
            "channels": [1, 2, 3, 4]
        }
    },
    "input": {
        "modules": ["http"],
        "http": {
            "ip": "127.0.0.1",
            "port": 8080,
            "accepts": [
                "http://localhost:8080",
                "http://127.0.0.1:8080"
            ]
        },
        "tcp": {
            "ip": "127.0.0.1",
            "port": 50000
        }
    },
    "dmx": {
        "groups": {
            "stage": {
                "name": "Stage",
                "devices": [
                    { "model": "dimmer", "channel": 1, "max": [255] },
                    { "model": "wclight", "channel": 10, "max": [255, 255, 0] }
                ]
            }
        },
        "fadeInterval": 0.7,
        "delay": 0,
        "fps": 30
    }
}
```

## `input`

### `input.modules`

List of enabled input modules:

- `http`: HTTP API and production frontend
- `tcp`: plain TCP command server

The normal local configuration should include `http`.

### `input.http`

| Field | Type | Description |
| --- | --- | --- |
| `ip` | string | Listen address. Use `127.0.0.1` for local-only access. |
| `port` | integer | HTTP port, normally `8080`. |
| `accepts` | string array | Allowed CORS origins. |

Binding to a non-loopback address exposes the API and configuration interface to the network. DMXBOX does not currently provide application-level authentication, so only do this on a trusted network with appropriate host firewall rules.

### `input.tcp`

| Field | Type | Description |
| --- | --- | --- |
| `ip` | string | TCP listen address. |
| `port` | integer | TCP listen port, normally `50000`. |

See the [TCP Protocol](tcp-protocol.md) for command framing, responses, commands, and compatibility aliases.

## `output`

### `output.target`

List of enabled outputs:

- `console`: log DMX output without hardware
- `ftdi`: serial DMX output
- `artnet`: Art-Net output
- `osc`: OSC mute/unmute output

Multiple outputs can be enabled together. At least one DMX renderer (`console`, `ftdi`, or `artnet`) is required by the DMX module.

### `output.ftdi`

| Field | Type | Description |
| --- | --- | --- |
| `port` | string | Serial device, for example `COM3` or `/dev/ttyUSB0`. |

### `output.artnet`

| Field | Type | Range | Description |
| --- | --- | --- | --- |
| `addr` | string | — | Destination IPv4 address. It must be reachable through a local interface. |
| `port` | integer | optional | UDP port. Defaults to 6454. |
| `net` | integer | 0–127 | Art-Net Net field. |
| `subuni` | integer | 0–15 | Art-Net Sub-Net field. |
| `universe` | integer | 0–15 | Art-Net Universe field. |

### `output.osc`

| Field | Type | Description |
| --- | --- | --- |
| `ip` | string | OSC receiver address. |
| `port` | integer | OSC receiver UDP port. |
| `format` | string | OSC address template. Every `{}` is replaced with a channel. |
| `type` | string | `int` for `int32`, or `float` for `float32`. |
| `inverse` | boolean | Invert the mute value sent to the receiver. |
| `channels` | integer array | Values substituted into the address template. |

## `dmx`

| Field | Type | Description |
| --- | --- | --- |
| `groups` | object | Map from stable group IDs to group definitions. |
| `fadeInterval` | number | Default fade duration in seconds. |
| `delay` | number | Configured effect delay in seconds. |
| `fps` | number | DMX update rate. The frontend accepts 0–99; use a positive value. |

### Groups

Each key in `dmx.groups` is used by the HTTP and TCP control APIs. A group contains:

| Field | Type | Description |
| --- | --- | --- |
| `name` | string | Display name shown in the browser. |
| `devices` | array | Devices controlled together by the group. |

### Devices

| Field | Type | Description |
| --- | --- | --- |
| `model` | string | `dimmer` or `wclight`. |
| `channel` | integer | One-based DMX start channel, 1–512. |
| `max` | integer array | Per-channel Fade In maximum values, normally 0–255. |

Device layouts:

- `dimmer`: one channel; `max` requires one value.
- `wclight`: three channels: cool white, warm white, and flash. The flash channel is currently always zero.

The complete device must fit within the 512-channel universe. For example, `wclight` can start no later than channel 510.

See [Lighting and Output Specifications](lighting-specifications.md) for serial DMX, Art-Net, device behavior, fades, and shutdown behavior.

## Development Frontend Configuration

`frontend/public/config.json` is separate from the backend configuration:

```json
{
    "backendPort": 8080
}
```

It tells the Vite development frontend which local backend port to use. In a production release, both frontend and backend are served by the Go HTTP server.
