# DMXBOX Lighting and Output Specifications

[日本語](../ja/lighting-specifications.md) | [Documentation index](README.md) | [Configuration Reference](configuration.md)

## Scope

This document describes the lighting-related behavior implemented by DMXBOX: the internal DMX universe, supported logical device models, fades, serial DMX output, Art-Net output, and operational limitations. It is not a compatibility guarantee for every USB adapter or fixture; verify each production setup with the actual hardware.

## DMX Universe

| Property | Value |
| --- | --- |
| Slots | 512 |
| Slot numbering in configuration | 1–512 |
| Internal value | 8-bit, 0–255 |
| Initial universe state | All zero |
| Configurable update rate | `dmx.fps` |

Device start channels are one-based. Every channel used by a device must fit within the universe. Overlapping devices write to the same slots and are not rejected by the current implementation; avoid overlaps unless the resulting behavior is explicitly intended and tested.

## Logical Device Models

### `dimmer`

| Relative slot | Purpose | Fade In target |
| --- | --- | --- |
| 1 | Intensity | `max[0]` |

The model uses one DMX slot. Fade Out targets zero.

Example:

```json
{
    "model": "dimmer",
    "channel": 1,
    "max": [255]
}
```

### `wclight`

| Relative slot | Purpose | Fade In target |
| --- | --- | --- |
| 1 | Cool white | `max[0]` |
| 2 | Warm white | `max[1]` |
| 3 | Flash/strobe | `max[2]`; the frontend normally keeps this at zero |

The model uses three consecutive slots, so its latest valid start channel is 510. Fade Out targets zero for all three slots.

Example:

```json
{
    "model": "wclight",
    "channel": 10,
    "max": [128, 128, 0]
}
```

DMXBOX does not identify physical fixture models through RDM. The logical model must match the channel mode configured on the fixture.

## Groups

A group contains one or more logical devices. Fade and cut operations apply to every device in the group. Group IDs are used by the HTTP and TCP APIs, while `name` is displayed in the browser.

The current implementation does not reject:

- devices that overlap other devices;
- the same device slots appearing in multiple groups;
- incompatible fixture modes configured on the physical equipment.

Treat the configuration as the patch sheet and review it before enabling physical output.

## Fade and Cut Behavior

At the start of an operation, each device captures its current universe values. It then linearly interpolates each slot toward:

- its configured `max` values for Fade In;
- zero for Fade Out.

`dmx.fadeInterval` is the default duration in seconds. A command-specific `duration` overrides it. `duration:0` performs a one-update cut. A positive `interval` delays the start of the operation.

Starting another fade on a device replaces that device's timing and target using the values present when the new operation starts.

## Output Scheduling

The DMX loop runs at `dmx.fps`. Renderers receive an update when a device changes and periodically while unchanged. Use a positive, practical frame rate; 30 FPS is the default.

## Console Output

The `console` renderer is intended for development and configuration checks without physical output hardware. It satisfies the DMX renderer requirement and allows group behavior to be inspected through logs.

## Serial DMX (`ftdi`)

The current serial renderer uses a serial port exposed by the operating system. Despite the configuration name, compatibility depends on the adapter and driver presenting a port that supports the required break operation.

| Serial property | Value |
| --- | --- |
| Baud rate | 250000 |
| Data bits | 8 |
| Parity | None |
| Stop bits | 2 |
| Break before frame | 1 millisecond |
| Start code | `0x00` |
| Data slots | 512 |

Each output writes a break, one zero start-code byte, and the complete 512-byte universe.

If output fails, DMXBOX closes the port and attempts to reopen it on a later output. Reopen errors are logged indirectly through subsequent failures; verify recovery with the actual adapter.

Example port names:

- Windows: `COM3`
- Linux: `/dev/ttyUSB0`

## Art-Net

DMXBOX sends ArtDMX UDP packets.

| Property | Value |
| --- | --- |
| Default UDP port | 6454 |
| Protocol version | 14 |
| Payload length | 512 slots |
| Sequence | Incrementing 8-bit value; zero is skipped after wrap |
| Physical byte | 0 |
| Port-Address | Combined from `net`, `subuni`, and `universe` |

The source interface is selected by finding a local interface whose network contains the configured destination address. The configured address must therefore be reachable through one of the machine's interfaces.

Art-Net input is not part of the public DMXBOX feature set. The implementation opens the UDP socket for transport and diagnostic receive logging, but received Art-Net data does not modify the DMX universe.

### Receive Timeout

The receive loop uses a 500 ms read deadline so that it can observe stop requests. Waiting and transmission continue when no packet arrives before the deadline. A socket read error other than a timeout stops Art-Net and is logged.

## OSC Output

OSC is used for mixer mute/unmute control rather than DMX fixture data. For every configured channel, DMXBOX replaces `{}` in the OSC address template and sends either an `int32` or `float32` value. `inverse` swaps the mute and unmute values.

## Startup, Reload, and Shutdown

- The internal DMX universe starts at zero.
- Initializing a logical device clears the universe before normal rendering begins.
- Saving configuration reloads modules and recreates configured outputs.
- Serial and Art-Net transports are closed during normal shutdown.
- DMXBOX does not currently send an explicit all-zero blackout frame immediately before closing outputs.

The physical result after shutdown depends on the receiver or fixture's signal-loss behavior. Define and test the required behavior at the system level; do not assume shutdown itself guarantees blackout.

## Production Verification Checklist

- Confirm the physical fixture mode matches each logical device model.
- Confirm every device fits within slots 1–512 and does not overlap unintentionally.
- Test slot boundaries, especially 1, 255/256, and the final slots.
- Test Fade In, Fade Out, Cut In, and Cut Out with actual fixtures.
- Test USB disconnection and reconnection if using serial DMX.
- Test network-interface selection and receiver loss if using Art-Net.
- Confirm the expected state during configuration reload and process shutdown.
- Keep a hardware-accessible blackout or power-control method for emergencies.
