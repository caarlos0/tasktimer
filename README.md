<p align="center">
  <img alt="header image" src="https://raw.githubusercontent.com/caarlos0/tasktimer/main/static/undraw_Dev_focus_re_6iwt.svg" height="250" />
  <h1 align="center">tasktimer</h1>
  <p align="center">Task Timer (<code>tt</code>) is a dead simple TUI task timer.</p>
</p>

---

## Usage

To get started, just run `tt`:

```sh
tt
```

You'll be presented with something like this:

<img width="1312" alt="image" src="https://user-images.githubusercontent.com/245435/132955599-312d5240-995a-41bc-b53d-d9cf714fd2b1.png">

You can just type a small description of what you're working on and press
<kbd>ENTER</kbd> to start timing.

At any time, press <kbd>ENTER</kbd> again to stop the
current timer or type a new task description and press <kbd>ENTER</kbd>
to stop the previous task and start the new one.

Each task will have its own timer, and the sum of all tasks will be displayed
in the header:

<img width="1312" alt="image" src="https://user-images.githubusercontent.com/245435/132955639-dea42092-c48a-478c-bbe1-e29fbf343c3c.png">

At any time, press <kbd>CTRL</kbd>+<kbd>c</kbd> to stop the current
timer (if any) and exit.

You can also press <kbd>ESC</kbd> stop the current task and blur the input
field and navigate  around a long list of tasks using the
arrow keys/page up/page down/etc.

## Report

You can extract a markdown file by running:

```sh
tt report
```

It will output the given project (via `-p PROJECT`) to `STDOUT`. You can
then save it to a file, pipe to another software or do whatever you like:

<img width="1312" alt="image" src="https://user-images.githubusercontent.com/245435/132955650-a2b0cfd1-eb38-4ecb-9116-20ca815fe01a.png">

## Edit

Let's say you forgot the timer running... you can edit it using the edit command:

```sh
tt edit
```

<img width="1312" alt="image" src="https://user-images.githubusercontent.com/245435/132955661-1349a06d-9c24-45ee-97a3-583ad8e066c9.png">

The project will be exporter to a JSON file and will open with your `$EDITOR`.
Once you close it, it will be imported over the old one.

You can also backup/edit/restore using `tt to-json` and `tt from-json`.

## Help

At any time, check `--help` to see the available options.

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

Download the pre-compiled binaries from the [releases page][releases] or clone the repo build from source.

## Stargazers over time

[![Stargazers over time](https://starchart.cc/caarlos0/clone-org.svg)](https://starchart.cc/caarlos0/clone-org)

## FAQ

### Where are data and logs stored?

Depends on the OS, but you can see yours running:

```sh
tt paths
```

## Stargazers over time

[![Stargazers over time](https://starchart.cc/caarlos0/tasktimer.svg)](https://starchart.cc/caarlos0/tasktimer)

[Badger]: https://github.com/dgraph-io/badger
[releases]:  https://github.com/caarlos0/tasktimer/releases

# Badges

[![Release](https://img.shields.io/github/release/caarlos0/tasktimer.svg?style=for-the-badge)](https://github.com/caarlos0/tasktimer/releases/latest)

[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge)](LICENSE.md)

[![Build](https://img.shields.io/github/workflow/status/caarlos0/tasktimer/build?style=for-the-badge)](https://github.com/caarlos0/tasktimer/actions?query=workflow%3Abuild)

[![Go Report Card](https://goreportcard.com/badge/github.com/caarlos0/tasktimer?style=for-the-badge)](https://goreportcard.com/report/github.com/caarlos0/tasktimer)

[![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=for-the-badge)](https://github.com/goreleaser)
