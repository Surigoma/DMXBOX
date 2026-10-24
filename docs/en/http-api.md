---
title: HTTP API
nav_order: 1
parent: External Integration
---

# HTTP API and Integration

[Documentation index](README.md) | [TCP Protocol](tcp-protocol.md)

## Prerequisites

These examples assume DMXBOX runs at `http://127.0.0.1:8080` and has a group named `stage`.
Use the key from `dmx.groups`, not its display name, in place of `stage`.
Control requests change actual output, so try them with `console` output first.
In Windows PowerShell, use `curl.exe` in place of `curl` if needed.

The running application's Swagger UI is available at `http://localhost:8080/docs/index.html`.

## Main Endpoints

| Method | Path | Purpose |
| --- | --- | --- |
| GET | `/api/v1/health` | HTTP server liveness |
| GET | `/api/version` | Running version |
| GET | `/api/features` | Enabled features |
| GET | `/api/endpoints` | Registered routes |
| GET | `/api/v1/config/fade` | Group configuration |
| GET | `/api/v1/fade-state` | Current software output states by group |
| POST | `/api/v1/fade/{group}` | Fade or Cut request |
| GET | `/api/v1/mute-state` | Last completely sent OSC mute state |
| POST | `/api/v1/mute` | Mute or Unmute request |
| GET | `/api/v1/control-mode` | Browser only state |
| GET | `/api/v1/control-state/stream` | Combined Fade and Mute state stream |
| GET | `/api/v1/operations?limit=100` | Recent operation log |
| GET | `/api/v1/config/all` | Complete configuration |
| POST | `/api/v1/config/save` | Save complete configuration |

## Start with Read Requests

```sh
curl "http://127.0.0.1:8080/api/v1/health"
curl "http://127.0.0.1:8080/api/v1/config/fade"
curl "http://127.0.0.1:8080/api/v1/control-mode"
```

Health returns `status: "ok"` and a timestamp.
Group configuration is an object keyed by group ID.
When `browserOnly` is `true`, control from external scripts is restricted.

## Fade In, Fade Out, and Cut

```sh
curl -X POST "http://127.0.0.1:8080/api/v1/fade/stage?isIn=true&duration=1.5&interval=0"
curl -X POST "http://127.0.0.1:8080/api/v1/fade/stage?isIn=false&duration=1.5&interval=0"
curl -X POST "http://127.0.0.1:8080/api/v1/fade/stage?isIn=false&duration=0&interval=0"
```

| Parameter | Meaning |
| --- | --- |
| `isIn` | `true` targets configured maximum values; `false` targets zero. Defaults to `true` |
| `duration` | Transition time in seconds. Defaults to `dmx.fadeInterval`; `0` performs a Cut |
| `interval` | Delay before starting, in seconds. Use `0` for immediate start |

Pass parameters in the query string, not a JSON body.
Quote the complete URL so the shell does not interpret `&`.
HTTP 200 with `{"result":"OK"}` means the request was handed to an internal module, not that a device completed it.
That response alone does not detect a missing group, so check group configuration first.

## OSC Mute

```sh
curl -X POST "http://127.0.0.1:8080/api/v1/mute?isMute=true"
curl -X POST "http://127.0.0.1:8080/api/v1/mute?isMute=false"
curl "http://127.0.0.1:8080/api/v1/mute-state"
```

`isMute` defaults to `true`. Enable the `osc` output and configure the intended receiver first.
State contains `isMute: true`, `false`, or `null` when unknown.
It records send completion, not an acknowledgment from the mixer.

## State Stream

```sh
curl -N "http://127.0.0.1:8080/api/v1/control-state/stream"
```

This is `application/x-ndjson`, not SSE. It sends an initial snapshot, then checks for changes every 250 milliseconds.
A JSON line is sent only when state changes. Empty lines every 15 seconds keep the connection alive.

```json
{"fade":{"stage":{"level":0.5,"state":"fading","isIn":true}},"mute":{"isMute":null}}
```

This is an illustrative snapshot. `level` is an output ratio from 0 to 1; `state` is `idle`, `waiting`, or `fading`.
Split incoming data on newline boundaries rather than network chunks and ignore empty lines before parsing JSON.
Reconnect after a disconnection and use the new initial snapshot as the current state.
If using a proxy, configure it to avoid buffering the stream.

## Saving Configuration

Read the complete object from `GET /api/v1/config/all`, edit it, and send it to
`POST /api/v1/config/save` with `Content-Type: application/json`.
This is not a partial-update PATCH API. A request containing only one field can fail validation or replace other fields with unintended values.

A successful save returns `result: true`. If `restartRequired` is `true`, restart DMXBOX to apply the saved update.
Even a normal reload can interrupt connections; integration clients must handle reconnection.
Malformed JSON or invalid settings return 400; a persistence failure returns 500.

## Errors and Control Restrictions

| Symptom | Check |
| --- | --- |
| Connection refused | Running process, listener address/port, firewall |
| Control returns 403 | Whether Browser only is enabled |
| Control returns 500 | Enabled output modules and internal message errors |
| 200 but no physical change | Group ID, output destination, hardware, runtime logs |
| Another browser origin fails | Origin settings in `input.http.accepts` |

Change Browser only through the Web UI. Its identifying headers are not authentication credentials.
CORS is not authentication either; restrict access to trusted networks.
Blind retries can duplicate operations. Check state and operation logs before retrying a control request.
