package main

import (
	"os"

	"github.com/caarlos0/tasktimer/internal/cmd"
)

var version = "dev"

func main() {
	cmd.Execute(
		version,
		os.Exit,
		os.Args[1:],
	)
}
