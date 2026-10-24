---
title: Operations and Troubleshooting
nav_order: 4
---

# Operations and Troubleshooting

[Documentation index](README.md) | [User Guide](user-guide.md)

## Before a Session

1. Back up the last working `config.json`.
2. Start DMXBOX from the directory containing the executable and `static/`.
3. Confirm that the control page shows the expected groups.
4. Check Browser only. Turn it off when external control is required.
5. Match fixture addresses and modes against the configured channel layout.
6. Test Fade In, Fade Out, and Cut for each group and observe the fixtures.
7. If using OSC, verify Mute and Unmute on the intended mixer channels.
8. If using external control, test a command from that software as well.

A successful health check confirms that the HTTP server responds.
Check physical DMX and mixer operation separately.

## Understanding State

| Display | Meaning | What to check |
| --- | --- | --- |
| Output percentage | Software output relative to configured maximum values | Whether fixtures actually follow it |
| `waiting` | Waiting for a delayed start | The command's `interval` |
| `fading` | A fade is active | `duration` and `dmx.fadeInterval` |
| `idle` | No waiting or active fade | Output can still be nonzero |
| `Muted` / `Unmuted` | Last completely sent OSC request | The actual mixer state |
| `State unavailable` | State cannot be obtained, or OSC state is not yet known | Connection, logs, and OSC destination |

Output percentage is not fade progress. A dimmer with a configured maximum of 128 reports 100% when it outputs 128.
Direct changes on the mixer are not automatically reported back to DMXBOX.

## Switching to Manual Control

1. Expand **Control mode** at the bottom of the settings page.
2. Turn **Browser only** on. No Update action is needed.
3. Confirm the mode on the control page and perform the required browser operations.
4. Stop commands from the external controller and investigate the cause.
5. When ready to return to external control, turn the mode off and verify one command.

This mode restricts Fade and Mute requests from sources such as TCP.
It does not stop an already active fade or black out the output.
The mode resets after restarting DMXBOX and is not an authentication mechanism.

## Reading the Operation Log

Open **Operation Log** from the application menu.
The page displays the latest 100 entries, newest first, and refreshes every two seconds.

| Column | Content |
| --- | --- |
| Time | Recorded time, displayed in the browser's local time zone |
| Source | Input module, such as `http` or `tcp` |
| Target | Destination module, such as `dmx` or `osc` |
| Action | `fade` or `mute` |
| Arguments | Requested group ID, direction, timing, and other arguments |

An entry records that a request was handed to an internal module.
It does not confirm fade completion or a physical device response.
Compare it with runtime diagnostic logs around the same time.

Entries are appended to `operations.jsonl` in the working directory, one JSON object per line.
At startup, DMXBOX reads recent history from that file and retains up to 1000 entries in memory.
There is no automatic rotation. Monitor file size and archive it while DMXBOX is stopped.
Check runtime logs for errors opening the file.

## Diagnosing Symptoms

| Symptom | Check in this order |
| --- | --- |
| A group is missing | Working directory → loaded configuration → `dmx.groups` → restart notice |
| Gauge moves but fixtures do not | Enabled output → port/destination → fixture addresses and modes → cabling |
| Fade finishes immediately | Cut action → `duration:0` → `fadeInterval` |
| Only external control fails | Browser only → TCP module enabled → restart → listening port → command delimiter |
| Mute state remains unknown | OSC enabled → restart → channel list → send errors |
| Connection drops after saving | HTTP listener changes → new URL → reload → runtime logs |
| Settings reset at every startup | Working directory → intended `config.json` → write permissions |

TCP `ack` and HTTP `result: OK` do not confirm successful physical output.
See [TCP Protocol](tcp-protocol.md) and [HTTP API and Integration](http-api.md).

## Changing Settings and Recovering

Changing output destinations or group layouts can recreate connections and reset output state.
Check device state first and make these changes during an appropriate pause in operation.

When saving reports that a restart is required, do not assume the saved settings are already active.
After restarting, check groups, output destinations, and external control again.
If problems persist, stop DMXBOX, restore the backed-up `config.json`, and start it again.

Shutdown does not explicitly send a blackout frame.
Perform the required Fade Out or device-side operation first, and verify loss-of-signal behavior.

## Information for a Bug Report

- Release version, operating system, and fixture/adapter models
- Time, reproduction steps, expected result, and actual result
- Relevant runtime and operation logs
- Configuration with private information removed
- Whether the request came from the browser, HTTP, or TCP
- Whether it reproduces with console output alone or requires hardware

Review IP addresses, equipment names, and operational details before publishing logs or configuration.
