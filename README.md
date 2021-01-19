# tasktimer

[![Release](https://img.shields.io/github/release/caarlos0/tasktimer.svg?style=flat-square)](https://github.com/caarlos0/tasktimer/releases/latest)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](LICENSE.md)
[![Build](https://img.shields.io/github/workflow/status/caarlos0/tasktimer/build?style=flat-square)](https://github.com/caarlos0/tasktimer/actions?query=workflow%3Abuild)
[![Go Report Card](https://goreportcard.com/badge/github.com/caarlos0/tasktimer?style=flat-square)](https://goreportcard.com/report/github.com/caarlos0/tasktimer)
[![Godoc](https://godoc.org/github.com/caarlos0/tasktimer?status.svg&style=flat-square)](http://godoc.org/github.com/caarlos0/tasktimer)
[![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=flat-square)](https://github.com/goreleaser)

Task Timer (`tt`) is a dead simple TUI task timer

## Usage

To get started, just run `tt`:

```sh
tt
```

You'll be presented with something like this:

<img width="1279" alt="image" src="https://user-images.githubusercontent.com/245435/104979833-f783f280-59e3-11eb-8c36-d6a086cd63fb.png">

You can just type a small description of what you're working on and press
<kbd>ENTER</kbd> to start timing.

At any time, press <kbd>ESC</kbd> or <kbd>ENTER</kbd> again to stop the
current timer or type a new task description and press <kbd>ENTER</kbd>
to stop the previous task and start the new one.

Each task will have its own timer, and the sum of all tasks will be displayed
in the header (along with a clock and the project name):

<img width="1279" alt="image" src="https://user-images.githubusercontent.com/245435/104979802-e804a980-59e3-11eb-8d79-51ac7b272c31.png">

At any time, press <kbd>CTRL</kbd>+<kbd>c</kbd> to stop the current
timer (if any) and exit.

## Report

You can extract a markdown file by running:

```sh
tt report
```

It will output the given project (via `-p PROJECT`) to `STDOUT`. You can
then save it to a file, pipe to another software or do whatever you like.

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
