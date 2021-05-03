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

<img width="1288" alt="image" src="https://user-images.githubusercontent.com/245435/114648989-1679e400-9cb6-11eb-8752-460b0d5eb3fb.png">

You can just type a small description of what you're working on and press
<kbd>ENTER</kbd> to start timing.

At any time, press <kbd>ENTER</kbd> again to stop the
current timer or type a new task description and press <kbd>ENTER</kbd>
to stop the previous task and start the new one.

Each task will have its own timer, and the sum of all tasks will be displayed
in the header:

<img width="1288" alt="image" src="https://user-images.githubusercontent.com/245435/114649087-3b6e5700-9cb6-11eb-800f-5daaa4baac13.png">

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

<img width="1288" alt="image" src="https://user-images.githubusercontent.com/245435/114649175-622c8d80-9cb6-11eb-8de9-063ebf412f7f.png">

## Edit

Let's say you forgot the timer running... you can edit it using the edit command:

```sh
tt edit
```

<img width="1288" alt="image" src="https://user-images.githubusercontent.com/245435/114649253-86886a00-9cb6-11eb-8d41-b7895f012f57.png">

The project will be exporter to a JSON file and will open with your `$EDITOR`.
Once you close it, it will be imported over the old one.

You can also backup/edit/restore using `tt to-json` and `tt from-json`.

## My terminal is light

Gotcha you covered! TaskTimer automatically handles that thanks to the awesome [lipgloss](https://github.com/charmbracelet/lipgloss):

<img width="1288" alt="image" src="https://user-images.githubusercontent.com/245435/114649473-fac30d80-9cb6-11eb-96d7-44d0d626e9d1.png">

## Help

At any time, check `--help` to see the available options.

## Install

```sh
brew install caarlos0/tap/tt
```

Or use any of the other provided means in the [releases page][releases].

## FAQ

### Where are data and logs stored?

Depends on the OS, but you can see yours running:

```sh
tt paths
```

## Upgrades

### From 1.0.x to 1.1.x

Data was moved from `~/tasktimer` to user data and user logs dir according to
the OS.

To move, run:

```sh
tt paths
```

It will print something like this:

```
Database path: /Users/carlos/Library/Application Support/tasktimer/default.db
Log path:      /Users/carlos/Library/Logs/tasktimer/default.log
```

We only need to migrate the data, so:

```sh
rm -rf "/Users/carlos/Library/Application Support/tasktimer/"*.db # make sure its empty
cp -rf ~/tasktimer/*.db "/Users/carlos/Library/Application Support/tasktimer/" # copy data
rm -rf ~/tasktimer # delete old folder

```

## Stargazers over time

[![Stargazers over time](https://starchart.cc/caarlos0/tasktimer.svg)](https://starchart.cc/caarlos0/tasktimer)

[Badger]: https://github.com/dgraph-io/badger
[releases]:  https://github.com/caarlos0/tasktimer/releases

# Badges

[![Release](https://img.shields.io/github/release/caarlos0/tasktimer.svg?style=flat-square)](https://github.com/caarlos0/tasktimer/releases/latest)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](LICENSE.md)
[![Build](https://img.shields.io/github/workflow/status/caarlos0/tasktimer/build?style=flat-square)](https://github.com/caarlos0/tasktimer/actions?query=workflow%3Abuild)
[![Go Report Card](https://goreportcard.com/badge/github.com/caarlos0/tasktimer?style=flat-square)](https://goreportcard.com/report/github.com/caarlos0/tasktimer)
[![Godoc](https://godoc.org/github.com/caarlos0/tasktimer?status.svg&style=flat-square)](http://godoc.org/github.com/caarlos0/tasktimer)
[![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=flat-square)](https://github.com/goreleaser)
