<p align="center">
  <img alt="header image" src="https://raw.githubusercontent.com/caarlos0/tasktimer/main/static/undraw_Dev_focus_re_6iwt.svg" height="250" />
  <h1 align="center">tasktimer</h1>
  <p align="center">Task Timer (<code>tt</code>) is a dead simple TUI task timer.</p>
</p>

<p align="center">
  <a href="https://github.com/caarlos0/tasktimer/releases/latest"><img src="https://img.shields.io/github/release/caarlos0/tasktimer.svg?style=for-the-badge" alt="Release"></a>
  <a href="LICENSE.md"><img src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge" alt="Software License"></a>
  <a href="https://github.com/caarlos0/tasktimer/actions?query=workflow%3Abuild"><img src="https://img.shields.io/github/actions/workflow/status/caarlos0/tasktimer/build.yml?style=for-the-badge" alt="Build"></a>
  <a href="https://goreportcard.com/report/github.com/caarlos0/tasktimer"><img src="https://goreportcard.com/badge/github.com/caarlos0/tasktimer?style=for-the-badge" alt="Go Report Card"></a>
</p>

---

- [Install](#install)
- [Usage](#usage)
  - [Projects](#projects)
  - [Keyboard shortcuts](#keyboard-shortcuts)
  - [Start, pause, and stop from the CLI](#start-pause-and-stop-from-the-cli)
  - [Report](#report)
  - [Edit](#edit)
  - [List projects](#list-projects)
  - [Backup and restore](#backup-and-restore)
- [FAQ](#faq)

## Install

**homebrew**:

```sh
brew install caarlos0/tap/tt
```

**apt**:

```sh
echo 'deb [trusted=yes] https://repo.caarlos0.dev/apt/ /' | sudo tee /etc/apt/sources.list.d/caarlos0.list
sudo apt update
sudo apt install tt
```

**yum**:

```sh
echo '[caarlos0]
name=caarlos0
baseurl=https://repo.caarlos0.dev/yum/
enabled=1
gpgcheck=0' | sudo tee /etc/yum.repos.d/caarlos0.repo
sudo yum install tt
```

**arch linux**:

```sh
yay -S tasktimer-bin
```

**deb/rpm/apk**:

Download the `.apk`, `.deb` or `.rpm` from the [releases page][releases] and install with the appropriate commands.

**manually**:

Download the pre-compiled binary for your platform from the [releases page][releases] and place it somewhere on your `$PATH`.

To build from source, see [CONTRIBUTING.md](CONTRIBUTING.md).

## Usage

To get started, just run `tt`:

```sh
tt
```

You'll be presented with something like this:

<img width="1312" alt="image" src="https://user-images.githubusercontent.com/245435/132955599-312d5240-995a-41bc-b53d-d9cf714fd2b1.png">

Type a description of what you're working on and press <kbd>Enter</kbd> to start timing.

Press <kbd>Enter</kbd> again (with an empty input) to stop the current timer, or type a new task description and press <kbd>Enter</kbd> to stop the previous task and immediately start the new one.

Each task has its own timer. The total time across all tasks is shown in the header:

<img width="1312" alt="image" src="https://user-images.githubusercontent.com/245435/132955639-dea42092-c48a-478c-bbe1-e29fbf343c3c.png">

### Projects

By default, all tasks are saved under a project named `default`. Use the `-p` flag to work with a named project:

```sh
tt -p myproject
```

This flag is supported by all subcommands.

### Keyboard shortcuts

| Key | Action |
|---|---|
| <kbd>Enter</kbd> | Start a new task / stop the current timer |
| <kbd>r</kbd> | Restart the selected task (copies its name into the input) |
| <kbd>p</kbd> | Pause or resume the current task |
| <kbd>d</kbd> | Delete the selected task |
| <kbd>ESC</kbd> | Stop the current task and blur the input field |
| <kbd>↑</kbd> / <kbd>↓</kbd> | Navigate the task list (when input is not focused) |
| <kbd>Page Up</kbd> / <kbd>Page Down</kbd> | Scroll the task list |
| <kbd>/</kbd> | Filter the task list |
| <kbd>Ctrl</kbd>+<kbd>C</kbd> | Stop the current timer and exit |

### Start, pause, and stop from the CLI

Control timers without launching the TUI — useful in scripts, shell aliases, or CI hooks:

```sh
tt start writing tests     # stop any running task, start a new one
tt pause                   # freeze the current timer (time stops accumulating)
tt resume                  # continue the paused task from where it left off
tt stop                    # stop all running or paused tasks
```

These commands respect the `-p` flag for named projects:

```sh
tt -p myproject start "fixing the build"
tt -p myproject pause
```

A paused task shows `[paused]` in the TUI task list and its elapsed time is frozen until resumed. Press <kbd>p</kbd> in the TUI to toggle pause/resume on the active task.

### Report

Generate a Markdown report for the current project:

```sh
tt report
# or
tt r
```

Output goes to `STDOUT`, so you can save or pipe it:

```sh
tt report > report.md
tt -p myproject report | pbcopy
```

Filter by date range with `--since` and `--until` (format: `YYYY-MM-DD`):

```sh
tt report --since 2026-05-19 --until 2026-05-23   # this week only
tt report --since 2026-05-22                       # today onwards
```

The total time shown in the report header reflects only the filtered tasks.

<img width="1312" alt="image" src="https://user-images.githubusercontent.com/245435/132955650-a2b0cfd1-eb38-4ecb-9116-20ca815fe01a.png">

### Edit

If you forgot to stop a timer, use the edit command to fix it:

```sh
tt edit
# or
tt e
```

<img width="1312" alt="image" src="https://user-images.githubusercontent.com/245435/132955661-1349a06d-9c24-45ee-97a3-583ad8e066c9.png">

The current project is exported to a temporary JSON file and opened in your `$EDITOR`. Save and close the file to apply the changes.

> **Note:** `$EDITOR` must be set in your environment (e.g. `export EDITOR=vim`).

### List projects

List all projects that have recorded data:

```sh
tt list
```

Add `-v` to see task counts and total time per project:

```sh
tt list -v
# default (42 tasks, 3h 15m)
# work    (18 tasks, 9h 40m)
```

### Backup and restore

Export and import task data as JSON:

```sh
tt to-json backup.json
tt from-json backup.json
```

## FAQ

### Where are data and logs stored?

Depends on the OS. Run this to see the paths used on your machine:

```sh
tt paths
```

### How do I get help for a specific command?

Pass `--help` to any command:

```sh
tt --help
tt report --help
```

## Stargazers over time

[![Stargazers over time](https://starchart.cc/caarlos0/tasktimer.svg)](https://starchart.cc/caarlos0/tasktimer)

[releases]: https://github.com/caarlos0/tasktimer/releases
