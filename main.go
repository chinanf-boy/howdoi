package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
	howdoi "github.com/chinanf-boy/howdoi/howdoi"
)

func main() {
	const _version = "0.0.1"
	res, err := argsPar()

	if res.Version {
		fmt.Printf("Version:%s", _version)
		return
	}

	if err != "" {
		fmt.Printf(err)
		return
	}

	howdoi.Howdoi(res)
}

func argsPar() (howdoi.Cli, string) {
	parser := argparse.NewParser("howdoi", "cli to Ask the question")

	color := parser.Flag("c", "color", &argparse.Options{Required: false, Help: "colorful Output"})
	version := parser.Flag("v", "version", &argparse.Options{Required: false, Help: "version"})
	num := parser.Int("n", "num", &argparse.Options{Required: false, Help: "how many answer"})
	query := parser.List("q", "query", &argparse.Options{Required: true, Help: "query what"})

	// Parse input
	err := parser.Parse(os.Args)
	var errStr string
	if err != nil {
		errStr = parser.Usage(err)
	}

	res := howdoi.Cli{Color: *color, Num: int8(*num), Query: *query, Version: *version}

	return res, errStr
}
