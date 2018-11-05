package main

import (
	"fmt"
	"os"
	"testing"

	howdoi "github.com/chinanf-boy/howdoi/pkg"
	"github.com/logrusorgru/aurora"
)

func TestExample(t *testing.T) {

	// use Lib for howdoi, ArgsPar get the howdoi.Cli struct

	exampleArgs := append(os.Args, "-q")
	exampleArgs = append(exampleArgs, "format date bash")

	res, err := howdoi.ArgsPar(exampleArgs)

	if res.Version {
		fmt.Printf(aurora.Green(name + ", version:" + version).String())
		fmt.Printf("date:%s", date)
		fmt.Printf("commit:%s", commit)
		return
	}

	if err != nil {
		fmt.Println("args parse fail", err)
		t.Fail()
	}

	// pass howdoi.Cli
	_, err = howdoi.Howdoi(res)

	if err != nil {
		fmt.Println("howdoi fail", err)

		t.Fail()
	}

	// for _, v := range result {
	// 	fmt.Println()
	// 	fmt.Println(v)
	// }
}
