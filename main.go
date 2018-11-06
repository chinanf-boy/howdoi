package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

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
		fmt.Printf(" date:%s", date)
		fmt.Printf(" commit:%s", commit)
		return
	}

	if err != nil {
		log.Fatalln(err)
	}

	cliChan := make(chan string)
	done := make(chan int)
	// pass howdoi.Cli
	go func() {
		n := 0
		for i := 0; i < 1; i++ {
			s := <-cliChan
			n, _ = strconv.Atoi(s)
			// fmt.Printf("form ChanHowdoi %s\n", s)
		}

		for i := 0; i < n; i++ {
			fmt.Println()
			fmt.Println(<-cliChan)

		}
		close(cliChan)
		done <- 0

	}()

	err = howdoi.ChanHowdoi(res, cliChan)

	if err != nil {
		log.Fatalln(err)
	}
	<-done
	close(done)
}
