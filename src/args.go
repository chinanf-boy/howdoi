package howdoi

import (
	"errors"
	"os"

	"github.com/akamensky/argparse"
)

// Cli args struct for cli
type Cli struct {
	Color   bool
	Num     int
	Query   []string
	Version bool
	Debug   bool
	Theme   string
	Cache   bool
}

// ArgsPar : get me parse OS.args with howdoi.Cli struct
func ArgsPar() (Cli, error) {
	parser := argparse.NewParser("howdoi", "cli to Ask the question")

	color := parser.Flag("c", "color", &argparse.Options{Required: false, Help: "colorful Output", Default: false})
	version := parser.Flag("v", "version", &argparse.Options{Required: false, Help: "version"})
	num := parser.Int("n", "num", &argparse.Options{Required: false, Help: "how many answer", Default: 1})
	query := parser.List("q", "query", &argparse.Options{Required: true, Help: "query what"})
	debug := parser.Flag("D", "debug", &argparse.Options{Required: false, Help: "debug *"})
	theme := parser.String("T", "theme", &argparse.Options{Required: false, Help: "chrome styles", Default: "pygments"})
	cache := parser.Flag("C", "cache", &argparse.Options{Required: false, Help: "cache [path]", Default: false})

	// Parse input
	err := parser.Parse(os.Args)
	var errStr error
	if err != nil {
		errStr = errors.New(parser.Usage(err))
	}

	res := Cli{
		Color:   *color,
		Num:     *num,
		Query:   *query,
		Version: *version,
		Debug:   *debug,
		Theme:   *theme,
		Cache:   *cache}

	return res, errStr
}
