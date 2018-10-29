package howdoi

import (
	"errors"
	"os"

	"github.com/akamensky/argparse"
)

// Cli args struct for cli
type Cli struct {
	Color   bool
	Num     int8
	Query   []string
	Version bool
}

// ArgsPar : get me parse OS.args with howdoi.Cli struct
func ArgsPar() (Cli, error) {
	parser := argparse.NewParser("howdoi", "cli to Ask the question")

	color := parser.Flag("c", "color", &argparse.Options{Required: false, Help: "colorful Output"})
	version := parser.Flag("v", "version", &argparse.Options{Required: false, Help: "version"})
	num := parser.Int("n", "num", &argparse.Options{Required: false, Help: "how many answer"})
	query := parser.List("q", "query", &argparse.Options{Required: true, Help: "query what"})

	// Parse input
	err := parser.Parse(os.Args)
	var errStr error
	if err != nil {
		errStr = errors.New(parser.Usage(err))
	}

	res := Cli{Color: *color, Num: int8(*num), Query: *query, Version: *version}

	return res, errStr
}
