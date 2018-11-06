package howdoi

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/logrusorgru/aurora"
	debug "github.com/visionmedia/go-debug"
)

// ChanHowdoi chan result
func ChanHowdoi(clis Cli, res chan<- string) error {
	clis.prePare()

	err := clis.getChanInstructions(res)

	return err
}

func (clis Cli) getChanInstructions(cliChan chan<- string) error {
	gLog := debug.Debug("getChanInstructions")
	gLog(red("starting"))

	var err error
	links := clis.getLinks()
	var questionLinks []string
	if len(links) > 0 {
		questionLinks = clis.getQuestions(links, isQuestion)

		gLog(gree(fmt.Sprintf("0.2. match questions links: %d", len(questionLinks))))
	}

	gLog("1. questions: %d", len(questionLinks))

	if len(questionLinks) > 0 {
		var n int
		answers := make([]string, 0)

		if clis.Num > len(questionLinks) { // user num o/r questionLinks len
			n = len(questionLinks)
		} else {
			n = clis.Num
		}
		gLog("1.1 truth Num: %d", n)

		cliChan <- strconv.Itoa(n)
		// checkd TODO: go func,
		var wg sync.WaitGroup

		gLog("1.2 full URLs: \n%v", sliceGoodFmt(questionLinks[:n]))
		for i := 0; i < n; i++ { // the bigger one
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				var res string
				answer := clis.getAnswer(questionLinks[i])
				if len(answer) == 0 { // no answer
					res = noAnswerMsg
				} else if n > 1 { // user want more answers
					comeFrom := fmt.Sprintf(answerHeader,
						starHeader,
						questionLinks[i],
						strings.Join(answer, "\n"))

					res = comeFrom

				} else { // one answer
					res = strings.Join(answer, "\n")
				}
				cliChan <- res // chan to cli
			}(i)
		}

		wg.Wait()

		gLog("2. answers: %v", string(len(answers)))

		return nil
	}

	err = errors.New(aurora.Red("howdoi fail").String())

	return err
}
