package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"

	"github.com/creativerezz/aio/internal/cli"
)

func main() {
	err := cli.Cli(version)
	if err != nil && !flags.WroteHelp(err) {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
