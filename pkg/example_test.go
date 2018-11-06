package howdoi

import (
	"fmt"
	"os"
	"testing"
)

func TestExample(t *testing.T) {

	// use Lib for howdoi, ArgsPar get the Cli struct

	exampleArgs := append(os.Args, "-q")
	exampleArgs = append(exampleArgs, "format date bash")

	res, err := ArgsPar(exampleArgs)

	if res.Version {
		return
	}

	if err != nil {
		fmt.Println("args parse fail", err)
		t.Fail()
	}

	// pass Cli
	_, err = Howdoi(res)

	if err != nil {
		fmt.Println("howdoi fail", err)

		t.Fail()
	}

	// for _, v := range result {
	// 	fmt.Println()
	// 	fmt.Println(v)
	// }
}
