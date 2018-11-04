package main

import (
	"fmt"
	"log"

	howdoi "github.com/chinanf-boy/howdoi/pkg"
	"github.com/logrusorgru/aurora"
)

var version string
var name = "howdoi-cli"

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
