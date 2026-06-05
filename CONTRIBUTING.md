# Contributing

Thank you for your interest in contributing to tasktimer!

## Prerequisites

- [Go](https://go.dev/dl/) 1.24 or later
- Git

Verify your Go installation:

```sh
go version
```

## Getting started

1. Fork the repository on GitHub and clone your fork:

```sh
git clone https://github.com/YOUR_USERNAME/tasktimer.git
cd tasktimer
```

2. Install dependencies:

```sh
go mod download
```

## Running the app

Run directly from source without building a binary:

```sh
go run . 
```

With a named project:

```sh
go run . -p myproject
```

Run a subcommand:

```sh
go run . report
go run . list
go run . paths
```

## Building a binary

```sh
go build -o tt .
```

On Windows:

```sh
go build -o tt.exe .
```

Then run the resulting binary:

```sh
./tt          # macOS/Linux
.\tt.exe      # Windows
```

## Running tests

```sh
go test ./...
```

With verbose output:

```sh
go test -v ./...
```

## Project structure

```
.
├── main.go                  # Entry point
├── internal/
│   ├── cmd/                 # CLI commands (Cobra)
│   │   ├── root.go          # Root command and wiring
│   │   ├── report.go        # tt report  (--since / --until flags)
│   │   ├── start.go         # tt start <title>
│   │   ├── stop.go          # tt stop
│   │   ├── pause.go         # tt pause
│   │   ├── resume.go        # tt resume
│   │   ├── edit.go          # tt edit
│   │   ├── list.go          # tt list  (-v flag)
│   │   ├── paths.go         # tt paths
│   │   ├── tojson.go        # tt to-json
│   │   └── fromjson.go      # tt from-json
│   ├── ui/                  # Bubble Tea TUI models
│   │   ├── main.go          # Root model (task list + input)
│   │   ├── project_timer.go # Running total timer + taskDuration helper
│   │   ├── common.go        # Shared styles
│   │   ├── markdown.go      # Markdown report renderer
│   │   └── json.go          # JSON serialization helpers
│   ├── model/
│   │   └── model.go         # Task data model (Task, ExportedTask)
│   └── store/
│       └── store.go         # Badger DB persistence layer
└── .github/
    └── workflows/
        └── build.yml        # CI: test + build + release
```

Key dependencies:

| Package | Purpose |
|---|---|
| [Bubble Tea](https://github.com/charmbracelet/bubbletea) | TUI framework |
| [Bubbles](https://github.com/charmbracelet/bubbles) | TUI components (list, text input, spinner) |
| [Badger v3](https://github.com/dgraph-io/badger) | Embedded key-value database |
| [Cobra](https://github.com/spf13/cobra) | CLI framework |
| [Glamour](https://github.com/charmbracelet/glamour) | Terminal Markdown rendering |

## Making changes

1. Create a branch for your change:

```sh
git checkout -b your-feature-name
```

2. Make your changes and ensure tests pass:

```sh
go test ./...
go build -o tt .
```

3. Commit and push:

```sh
git add .
git commit -m "short description of what and why"
git push origin your-feature-name
```

4. Open a pull request against the `main` branch.

## Guidelines

- Keep pull requests focused on a single change.
- If you're fixing a bug, include a clear description of the original behavior and the fix.
- If you're adding a feature, open an issue first to discuss it.
- Match the existing code style — run `go vet ./...` and `go fmt ./...` before committing.
