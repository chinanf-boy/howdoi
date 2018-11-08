package howdoi

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

func newArg(s int) []string {
	// []string{"-q", "format date bash", "-c", "-C"}
	nArr := make([]string, 0)
	nArr = append(nArr, os.Args[0])
	args := []string{"format date bash"}

	for _, v := range args {
		nArr = append(nArr, "-q")
		nArr = append(nArr, v)
	}
	if s > 0 {
		nArr = append(nArr, "-c")
		nArr = append(nArr, "-T")
		nArr = append(nArr, "github")
		nArr = append(nArr, "-C")
		nArr = append(nArr, "-n")
		nArr = append(nArr, "3")
	}
	if s > 1 {
		nArr = append(nArr, "-R")
		nArr = append(nArr, "-D")
	}
	fmt.Println("new args", nArr)

	return nArr
}

func TestExampleVersion(t *testing.T) {

	// use Lib for howdoi, ArgsPar get the Cli struct
	nArr := make([]string, 0)
	exampleArgs := append(nArr, os.Args[0])
	exampleArgs = append(exampleArgs, "-v")

	res, err := ArgsPar(exampleArgs)

	if res.Version {
		return
	}

	if err != nil {
		fmt.Println("args parse fail", err)
		t.Fail()
	}
}

func TestExample(t *testing.T) {

	// use Lib for howdoi, ArgsPar get the Cli struct
	exampleArgs := newArg(1)

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
}

func TestExampleDebug(t *testing.T) {

	// use Lib for howdoi, ArgsPar get the Cli struct
	exampleArgs := newArg(2)

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
}

func TestChanExample(t *testing.T) {

	// use Lib for howdoi, ArgsPar get the Cli struct
	exampleArgs := newArg(1)

	res, err := ArgsPar(exampleArgs)

	if res.Version {
		return
	}

	cliChan := make(chan string)
	done := make(chan int)

	if err != nil {
		t.Fatal(err)
	}
	go func() {
		n := 0
		for i := 0; i < 1; i++ {
			s := <-cliChan
			n, _ = strconv.Atoi(s)
			// fmt.Printf("form ChanHowdoi %s\n", s)
		}

		for i := 0; i < n; i++ {
			<-cliChan

		}
		close(cliChan)
		done <- 0

	}()
	// pass Cli
	err = ChanHowdoi(res, cliChan)

	if err != nil {
		t.Fatal(err)
	}
	<-done
	close(done)
}
