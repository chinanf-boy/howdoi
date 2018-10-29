package main

import (
	"fmt"
	"os"

	howdoi "github.com/chinanf-boy/howdoi/howdoi"
)

func main() {
	const _version = "0.0.1"

	// use Lib for howdoi, ArgsPar get the howdoi.Cli struct
	res, err := howdoi.ArgsPar()

	if res.Version {
		fmt.Printf("Version:%s", _version)
		return
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	// pass howdoi.Cli
	_, err = howdoi.Howdoi(res)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
