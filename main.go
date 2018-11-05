package main

import (
	"fmt"
	"log"
	"os"

	howdoi "github.com/chinanf-boy/howdoi/pkg"
	"github.com/logrusorgru/aurora"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	name    = "howdoi-cli"
)

func main() {

	// use Lib for howdoi, ArgsPar get the howdoi.Cli struct
	res, err := howdoi.ArgsPar(os.Args)

	if res.Version {
		fmt.Printf(aurora.Green(name + ", version:" + version).String())
		fmt.Printf("date:%s", date)
		fmt.Printf("commit:%s", commit)
		return
	}

	if err != nil {
		log.Fatalln(err)
	}

	// pass howdoi.Cli
	result, err := howdoi.Howdoi(res)

	if err != nil {
		log.Fatalln(err)
	}

	for _, v := range result {
		fmt.Println()
		fmt.Println(v)
	}
}
