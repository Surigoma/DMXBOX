# DMXBOX User Guide

[日本語](../ja/user-guide.md) | [Documentation index](README.md)

## Intended Environment

DMXBOX is primarily designed to run its backend and browser interface on the same computer. The default HTTP listener is `127.0.0.1:8080`, so it is not exposed to other computers unless the configuration is changed.

Release builds are provided for Windows x86-64, Linux x86-64, and Linux ARM64.

## Release Contents

Keep these items in the same directory:

```text
dmxbox_<platform>_<version>[.exe]
static/
[config.json]
```

DMXBOX resolves `config.json` and `static/` relative to the current working directory. `config.json` is optional on first start and is generated when missing. Start DMXBOX from the directory containing the executable and `static/`.

## First Start

1. Extract the release archive.
2. If the archive contains `config.json`, review its HTTP listener and output targets.
3. Start the executable.
4. Open `http://localhost:8080/gui/`.
5. Confirm the configuration before connecting or enabling production lighting hardware.

If `config.json` is missing, DMXBOX creates a default one in the working directory. The default output is `console`, which is useful for checking operation without DMX hardware.

## Screens

### Control

The control screen displays configured DMX groups. Each group provides Fade In and Fade Out controls. Cut controls perform an immediate transition. If OSC output is enabled, Mute and Unmute controls are also displayed.

One background gauge spans both Fade buttons and shows the group's current DMX output level relative to its configured maximum. It grows during Fade In and shrinks during Fade Out. The group name, waiting/fading/idle status and output percentage appear above the buttons, including changes made through TCP or another browser. The percentage indicates output level, not fade completion. Fade and mute states share one HTTP stream; the backend checks every 250 milliseconds and sends updates only when state changes. Values reflect the backend's software DMX output, rather than physical fixture feedback. If the stream disconnects, the screen shows `State unavailable` while keeping the control buttons usable, and reconnects automatically after 3 seconds.

Above Mute/Unmute, `Muted` or `Unmuted` indicates the last completely sent OSC command, and the corresponding button is highlighted. Changes made through TCP or another browser also appear here. This is not feedback from the receiving device. Before the first command, while sending, after a send error or when status retrieval fails, the screen shows `State unavailable`. Both controls remain usable when the state is unknown.

### Settings

The settings screen edits:

- HTTP and TCP input
- FTDI, Art-Net, and OSC output
- DMX frame rate and fade timing
- DMX groups, device models, channels, and maximum values

Editing shows an unsaved-changes notice. Press **Update** to save the complete configuration; the notice clears after a successful save. Failed saves retain your edits. Navigating away with unsaved edits asks whether to stay or discard them. Saving reloads the backend modules, so active connections may be recreated.

Expand **Control mode** at the bottom of the settings screen to toggle **Browser only**. This mode blocks commands from TCP and other external controllers so you can operate manually from the browser if the control system has a problem. Changes take effect immediately without pressing **Update**, and the mode clears when the backend restarts. The current ON/OFF state also appears in the collapsed heading and on the control screen; failed status retrieval shows `Unavailable`.

## Stopping DMXBOX

Use Ctrl+C in the terminal. Wait for shutdown logging before closing the terminal or removing connected hardware.

## Updating

1. Stop the current process.
2. Back up `config.json`.
3. Extract the new release into a separate directory.
4. Copy the existing configuration after reviewing release notes for configuration changes.
5. Start the new executable and check `/api/v1/health` and the control screen.
6. Keep the previous release until the new version has been verified with the actual hardware.

## Troubleshooting

### The browser interface does not open

- Confirm that the process is still running.
- Open `http://localhost:8080/gui/` rather than the development URL.
- Check whether another application is using port 8080.
- Confirm `input.http.ip` and `input.http.port` in `config.json`.

### Configuration cannot be loaded

- Validate the JSON syntax.
- Start DMXBOX from the directory containing `config.json`.
- Restore a known-good configuration or temporarily move the invalid file so a default can be generated.

### Serial DMX cannot be opened

- Confirm the configured port, such as `COM3` or `/dev/ttyUSB0`.
- Close other programs using the device.
- On Linux, confirm that the user can access the serial device.
- Reconnect the USB device and verify whether its device name changed.

### Art-Net output is not received

- Check the destination address and network interface.
- Confirm Net, Sub-Universe, and Universe values.
- Check the operating-system firewall and receiver configuration.

### OSC mute does not work

- Confirm that `osc` is present in `output.target`.
- Check the target IP, port, address format, data type, and channel list.
- Confirm whether the receiving mixer expects normal or inverse values.

Runtime diagnostics are written to standard output. Include the relevant log lines, platform, release version, and redacted configuration when reporting a problem.

## Further Reading

- [Configuration Reference](configuration.md)
- [TCP Protocol](tcp-protocol.md)
- [Lighting and Output Specifications](lighting-specifications.md)
- [Developer Guide](developer-guide.md)
- API documentation: `http://localhost:8080/docs/index.html`
