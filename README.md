# tasktimer

[![Release](https://img.shields.io/github/release/caarlos0/tasktimer.svg?style=flat-square)](https://github.com/caarlos0/tasktimer/releases/latest)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](LICENSE.md)
[![Build](https://img.shields.io/github/checks-status/caarlos0/tasktimer/master?style=flat-square)](https://github.com/caarlos0/tasktimer/actions)
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

If ran with no additional flags, it will create a `default.md` in the
current working directory.

## Help

At any time, check `--help` to see the available options.

## FAQ

### Where is my data stored?

Data is stored in `~/tasktimer/{projectname}.db`.

[Badger][] is used as database, which means you can open it with the badger
CLI if you want to.

### Are there any logs?

Yes, they are written to `~/tasktimer/{projectname}.log`.

[Badger]: https://github.com/dgraph-io/badger
