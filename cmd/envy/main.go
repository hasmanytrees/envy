package main

import (
	"envy/internal/app/cmd"
	"os"
)

// envy init SHELL
// envy export
// envy gen

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
