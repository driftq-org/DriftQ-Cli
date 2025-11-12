package main

import (
	"fmt"
	"os"

	"github.com/BehnamAxo/driftq-cli/pkg/cli"
)

func main() {
	if err := cli.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
