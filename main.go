package main

import (
	"fmt"
	"log"

	howdoi "github.com/chinanf-boy/howdoi/src"
	"github.com/logrusorgru/aurora"
)

const (
	version = "0.3.0"
	name    = "howdoi-cli"
)

func main() {

	// use Lib for howdoi, ArgsPar get the howdoi.Cli struct
	res, err := howdoi.ArgsPar()

	if res.Version {
		fmt.Printf(aurora.Green(name + ", version:" + version).String())
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
