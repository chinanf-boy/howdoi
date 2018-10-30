package howdoi

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/logrusorgru/aurora"
)

var (
	verifySslCertificate bool
	scheme               string
	uRL                  string
	starHeader           string
	userAgents           []string
	searchUrls           map[string]string
	answerHeader         string
	noAnswerMsg          string
	xdgCacheDir          string
	cacheDir             string
	cacheFile            string
	howdoiSession        string
	searchEngine         string
)

func init() {
	if len(os.Getenv("HOWDOI_DISABLE_SSL")) > 0 {
		scheme = "http://"
		verifySslCertificate = false
	} else {
		scheme = "https://"
		verifySslCertificate = true
	}
	uRL = getEnv("HOWDOI_URL", "stackoverflow.com")
	searchEngine = getEnv("HOWDOI_SEARCH_ENGINE", "bing")

	searchUrls = map[string]string{
		"bing":   scheme + "www.bing.com/search?q=%s site:%s",
		"google": scheme + "www.google.com/search?q=%s site:%s",
	}

	// some value ,  waiting to use
	userAgents = []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.7; rv:11.0) Gecko/20100101 Firefox/11.0",
	}
	starHeader = "\u2605"
	answerHeader = "%s  Answer from %s \n\n %s"
	noAnswerMsg = "< no answer given >"
}

// Howdoi string
func Howdoi(res Cli) ([]string, error) {

	res.Query = []string{strings.Replace(strings.Join(res.Query[:], " "), "?", "", -1)}

	result, err := res.getInstructions()

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (clis Cli) getInstructions() ([]string, error) {
	var err error
	links := clis.getLinks() // HERE
	var questionLinks []string
	if len(links) > 0 {
		questionLinks = clis.getQuestions(links, isQuestion)
		if searchEngine == "google" {
			questionLinks = cutURL(questionLinks)
		}
	}
	if len(questionLinks) > 0 {
		var n int
		answers := make([]string, 0)

		if clis.Num > len(questionLinks) {
			n = len(questionLinks)
		} else {
			n = clis.Num
		}
		for i := 0; i < n; i++ {
			var res string
			answer := clis.getAnswer(questionLinks[i])
			if len(answer) == 0 {
				res = noAnswerMsg
			} else if n > 1 {
				comeFrom := fmt.Sprintf(answerHeader,
					starHeader,
					questionLinks[i],
					strings.Join(answer, "\n"))

				res = comeFrom

			} else {
				res = strings.Join(answer, "\n")
			}
			answers = append(answers, res)
		}

		return answers, nil

	}
	err = errors.New(aurora.Red("no questions link").String())

	return nil, err
}

func (clis Cli) getAnswer(u string) []string {
	doc := getResult(u)
	return extractAnswer(doc)
}

func extractAnswer(doc *goquery.Document) []string {
	links := make([]string, 0)
	instructions := doc.Find(".answer").First().Find("pre")
	if instructions.Size() > 0 {
		instructions.Each(func(i int, s *goquery.Selection) {
			str := s.Text() // TODO: colored code with term
			links = append(links, str)
		})
	} else {
		links = append(links, doc.Find(".post-text").Eq(1).Text())
	}

	return links
}

func (clis Cli) getQuestions(links []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range links {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func (clis Cli) getLinks() []string {

	searchURL := getSearchURL(searchEngine)
	u, _ := url.Parse(fmt.Sprintf(searchURL, clis.Query[0], uRL))

	q := u.Query()
	u.RawQuery = q.Encode() //urlencode
	doc := getResult(u.String())
	return extractLinks(doc, searchEngine)
}

func getResult(u string) *goquery.Document {
	var res *http.Response
	var err error

	proxyIs := whichProxy()

	if proxyIs == SOCKS {
		httpClient := Socks5Client()
		res, err = httpClient.Get(u)
	} else {
		res, err = http.Get(u)
	}

	if err != nil {
		log.Fatalln(aurora.Red("请求失败:"+searchEngine), err)
	} else {
		defer res.Body.Close()
	}

	if res.StatusCode != 200 {
		log.Fatalln(aurora.Red("status code error:"), res.Request.URL, res.Status)
	} else {
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatalln(aurora.Red("goquery.NewDocumentFromReader error:"), err)
		}
		return doc
	}
	return nil
}
