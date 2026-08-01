# DMXBOX TCP Protocol

[日本語](../ja/tcp-protocol.md) | [Documentation index](README.md) | [Configuration Reference](configuration.md)

## Purpose

The optional TCP module accepts compact commands from local control software or dedicated control devices and routes them to DMX fade or OSC mute actions.

Enable it with:

```json
{
    "input": {
        "modules": ["http", "tcp"],
        "tcp": {
            "ip": "127.0.0.1",
            "port": 50000
        }
    }
}
```

## Transport

| Property | Current behavior |
| --- | --- |
| Transport | TCP |
| Default address | `127.0.0.1:50000` |
| Direction | Client sends commands; server sends acknowledgements |
| Command separators | CRLF (`\r\n`), LF (`\n`), or CR (`\r`) |
| Read buffer | 512 bytes |
| Acknowledgement | ASCII `ack\r\n` |
| Connection lifetime | The server accepts more than one read on a connection; one-command-per-connection is supported and recommended |

Commands and identifiers are expected to use ASCII-compatible text. The implementation does not perform an explicit character-encoding conversion.

## Command Grammar

```text
fade-command = ("fadeIn" | "fadeOut") SP group [SP options]
options      = option *("," option)
option       = name ":" value
mute-command = "mute" [SP boolean]
boolean      = "true" | "false"
```

Only `duration` and `interval` have defined fade semantics. Unknown option names are currently forwarded internally but ignored by the DMX handler.

## Fade Commands

```text
fadeIn <group>
fadeOut <group>
```

- `<group>` is a key from `dmx.groups`.
- `fadeIn` transitions each device toward its configured `max` values.
- `fadeOut` transitions each device toward zero.
- Without options, `dmx.fadeInterval` is used as the duration.

Examples:

```text
fadeIn stage
fadeOut stage
```

### Per-command Options

```text
fadeIn <group> duration:<seconds>,interval:<seconds>
fadeOut <group> duration:<seconds>,interval:<seconds>
```

| Option | Meaning |
| --- | --- |
| `duration` | Transition duration in seconds. `0` performs an immediate cut. |
| `interval` | Delay before the transition starts, in seconds. Values greater than zero delay the start. |

Examples:

```text
fadeIn stage duration:1.5,interval:0
fadeOut stage duration:0,interval:0
```

Values are parsed as 32-bit floating-point numbers. Invalid values are ignored and the configured defaults are used.

## Mute Commands

```text
mute
mute true
mute false
```

- `mute` and `mute true` request mute.
- `mute false` requests unmute.
- The OSC module must be enabled through `output.target` for the command to reach a receiver.

## Compatibility Aliases

Aliases are expanded only when a complete single-token command is received.

| Alias | Expanded command | Availability |
| --- | --- | --- |
| `fi` | `fadeIn stg` | When group `stg` exists |
| `fo` | `fadeOut stg` | When group `stg` exists |
| `ci` | `fadeIn stg interval:0` | When group `stg` exists |
| `co` | `fadeOut stg interval:0` | When group `stg` exists |
| `fai` | `fadeIn aud` | When group `aud` exists |
| `fao` | `fadeOut aud` | When group `aud` exists |
| `mute` | `mute true` | Always |
| `unmute` | `mute false` | Always |

`ci` and `co` currently set `interval:0` but do not set `duration:0`; therefore, they do not necessarily perform an immediate cut. Use an explicit `duration:0` command when an immediate transition is required.

## Responses

After processing each line fragment, the server writes:

```text
ack\r\n
```

The current acknowledgement means that the input was read and passed through the command parser. It does **not** guarantee that:

- the command name was recognized;
- the group exists;
- the DMX or OSC module accepted the internal message;
- the physical output device completed the operation.

There is currently no negative acknowledgement or structured error response.

## Client Example

Python:

```python
import socket

with socket.create_connection(("127.0.0.1", 50000), timeout=2) as client:
    client.sendall(b"fadeIn stage\r\n")
    response = client.recv(512)
    assert response == b"ack\r\n"
```

## Known Constraints

- TCP does not preserve application message boundaries. The current implementation parses each socket read independently rather than buffering until a complete line. Keep commands short, send one command at a time, and wait for `ack\r\n` before sending the next command.
- A single command must fit within the 512-byte read buffer.
- Empty and unknown commands may still receive `ack\r\n`.
- Whitespace is significant: use one ASCII space between the command, group, and options.
- Configuration reload recreates the TCP listener and can close or interrupt active connections.
- The protocol currently has no authentication or encryption. Bind it to `127.0.0.1` unless it is protected by a trusted network and host firewall.
