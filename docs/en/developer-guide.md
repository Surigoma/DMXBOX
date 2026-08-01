# DMXBOX Developer Guide

[日本語](../ja/developer-guide.md) | [Documentation index](README.md)

## Source Tree

```text
backend/                 Go application and tests
  artnet/                Art-Net transport
  config/                Configuration model and persistence
  dmxServer/             DMX groups, devices, timing, and outputs
  httpServer/            HTTP server and controllers
  oscServer/             OSC message output
  packageModule/         Module lifecycle and message routing
  tcpServer/             TCP command input
frontend/                React and TypeScript application
  src/component/         Control and configuration components
  src/contexts/          Frontend runtime configuration
  src/routes/            TanStack Router routes
  src/test/              Shared browser-test helpers
.github/workflows/       Tag-triggered release workflow
Taskfile.yml             Repository-level task entry points
```

## Architecture

`backend/main.go` loads configuration, registers enabled modules, initializes the module manager, and starts each module. Modules communicate through `message.Message` values routed by `packageModule.ModuleManagerType`.

The main modules are:

- `http`: HTTP API and static frontend
- `tcp`: optional TCP command listener
- `dmx`: DMX state, fades, frame timing, and renderer output
- `osc`: optional OSC mute output

The DMX server owns a 512-byte universe. Device models update slices of that universe, while configured renderers send the result to console, serial DMX, or Art-Net.

The frontend uses TanStack Router, Material UI, React Hook Form, SWR, and Zod. In development it runs separately through Vite. In production the generated `static/` directory is served by the Go HTTP server.

## Requirements

- [Task](https://taskfile.dev/)
- Go compatible with `backend/go.mod`
- Node.js compatible with Vite 8
- Corepack and Yarn 4
- Chromium for frontend browser tests

Optional tools used by development tasks or reports include Air, swag, gotestsum, gocover-cobertura, and gopogh.

## Setup and Development

Enable Corepack and install dependencies as required by the local environment, then run:

```sh
task dev
```

Development endpoints:

- Frontend: `http://localhost:5173`
- Backend: `http://localhost:8080`
- Swagger UI: `http://localhost:8080/docs/index.html`

The frontend backend-port setting is `frontend/public/config.json`. The backend development configuration is `backend/config.json`.

## Testing

Use the repository task instead of invoking the browser tests directly:

```sh
task test
```

This runs frontend browser tests headlessly with coverage, followed by backend tests and report generation. Reports are written to:

```text
frontend/test/
backend/test/
```

Watch mode:

```sh
task test_watch
```

Frontend lint:

```sh
cd frontend
yarn lint
```

Backend static analysis:

```sh
cd backend
go vet ./...
```

The backend report task intentionally continues after the test command so reports can be inspected when a local test fails. Review the test output and generated report rather than relying only on the final task exit status.

## API Development

Routes are registered in `backend/httpServer/httpServer.go`. Controllers live below `backend/httpServer/controller/`.

After changing Swagger annotations, regenerate the API documentation:

```sh
cd backend
task update_swag
```

Current API groups include health, configuration, fade control, OSC mute control, feature discovery, endpoint discovery, and version information. See the running Swagger UI for request and response details.

## Adding a Device Model

1. Implement a device constructor below `backend/dmxServer/devices/spec/`.
2. Set its model name and channel count.
3. Register the constructor in `dmxServer.DeviceTypes`.
4. Add the model and validation to `frontend/src/types.ts`.
5. Add a frontend editor below `frontend/src/component/settings/device/model/`.
6. Add backend and frontend tests, including channel-boundary cases.
7. Update the configuration reference.

## Adding an Output Renderer

1. Implement the renderer through `dmxServer/controller.Controller`.
2. Register it in `dmxServer.RenderTypes`.
3. Extend the backend configuration model.
4. Add the frontend schema and settings editor.
5. Test initialization, output, failure, reload, and finalization.
6. Document the configuration fields and platform requirements.

## Building

Build the current supported platform:

```sh
task build
```

Build Windows x86-64, Linux x86-64, and Linux ARM64:

```sh
task build_all
```

Run the default task to create distribution directories, build all targets, and copy the current backend configuration:

```sh
task
```

Artifacts are written below `dist/`. The frontend is built into `dist/static/`.

Taskfile version variables and the embedded `backend/.version` file currently control artifact and runtime versions. Keep them synchronized when preparing a release.

## Release Process

The project uses local verification as the primary release gate. GitHub Actions is intentionally limited to manually created `v*` tags.

Suggested checklist:

1. Confirm the working tree and intended release changes.
2. Run `task test` and inspect both frontend and backend results.
3. Run `yarn lint` in `frontend/` and `go vet ./...` in `backend/`.
4. Run `task build_all` from a clean distribution directory.
5. Test the packaged frontend and `/api/v1/health`.
6. Verify FTDI, Art-Net, fade, cut, TCP, and OSC behavior as applicable.
7. Update versions and release notes.
8. Commit the release state.
9. Create and push the chosen `v*` tag.
10. Verify the GitHub Release and downloaded archive.

## Contribution Notes

- Keep unrelated local changes intact.
- Format Go code with `gofmt`.
- Keep TypeScript lint clean.
- Add regression tests for behavior changes.
- Keep English and Japanese documentation aligned when changing public behavior.
- Generated `frontend/src/routeTree.gen.ts` should not be edited manually.

## Related Documentation

- [User Guide](user-guide.md)
- [Configuration Reference](configuration.md)
- [TCP Protocol](tcp-protocol.md)
- [Lighting and Output Specifications](lighting-specifications.md)
- [Root README](../../README.md)
