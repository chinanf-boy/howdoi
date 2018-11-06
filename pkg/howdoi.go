package howdoi

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/logrusorgru/aurora"
	debug "github.com/visionmedia/go-debug"
)

// Howdoi string
func Howdoi(cli Cli) ([]string, error) {

	cli.prePare()

	result, err := cli.getInstructions()

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (clis Cli) getInstructions() ([]string, error) {
	gLog := debug.Debug("getInstructions")

	var err error
	links := clis.getLinks() // HERE
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
		// checkd TODO: go func,
		var wg sync.WaitGroup

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
				answers = append(answers, res) // add answer result

			}(i)
		}

		wg.Wait()

		gLog("2. answers: %v", string(len(answers)))

		return answers, nil
	}

	err = errors.New(aurora.Red("howdoi fail").String())

	return nil, err
}
