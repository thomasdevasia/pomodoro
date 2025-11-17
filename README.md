# Pomodoro — Terminal Pomodoro Timer

A small and simple terminal-based Pomodoro timer built with Bubble Tea.

## Prerequisites

- Go 1.25+ installed and available in your `PATH`.
- (Optional, Linux) `notify-send` (libnotify) and `spd-say` (speech-dispatcher) for desktop & spoken notifications.

## Build

Build a local binary:

```bash
go build -o bin/pomodoro ./cmd/pomodoro
```

Run directly (development):

```bash
go run ./cmd/pomodoro
```

Install to `GOBIN` / `GOPATH/bin`:

```bash
go install ./cmd/pomodoro
# then run as
pomodoro
```

## Reproducible-ish build

Use these flags to reduce embedded environment differences and produce a smaller, more deterministic binary:

```bash
export SOURCE_DATE_EPOCH=1700000000
CGO_ENABLED=0 go build -trimpath -buildvcs=false -ldflags="-s -w" -o bin/pomodoro ./cmd/pomodoro
```

For full reproducibility, build inside CI or a pinned container image with a fixed Go version and module cache.

## Usage

1. Start the app:

```bash
./bin/pomodoro
```

2. Fill the form:

- Enter Title (defaults to `Focus Time` when empty)
- Give Description (optional)
- Select Duration (in minutes)

3. While running, the UI shows a timer and a progress bar. When the timer finishes you will get a desktop notification and a spoken alert (if `notify-send` and `spd-say` are available).

## Keybindings

- `ctrl+c` : Quit

## Development

- Main entry: `./cmd/pomodoro`.
- UI / logic: `internal/controller/controller.go`.
- Common commands (from repo root):

```bash
make build        # build local binary to ./bin/pomodoro
make repro-build  # reproducible-ish build
make install      # install to GOBIN / GOPATH/bin
make run          # build then run
make test         # run tests (if any)
```

## Recommended additions

- `LICENSE` (MIT/Apache) — included here as MIT by default.
- `.gitignore` to exclude `bin/` and compiled artifacts.
- CI workflow: run `go test`, `go vet`, and `go build` on PRs.
- Use `goreleaser` if you want multi-arch release artifacts.

## License

This project is licensed under the MIT License — see the `LICENSE` file for details.
