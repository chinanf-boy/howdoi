package howdoi

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"

	userAgent "github.com/EDDYCJY/fake-useragent"
	"github.com/PuerkitoBio/goquery"
	"github.com/logrusorgru/aurora"
	debug "github.com/visionmedia/go-debug"
)

var (
	verifySslCertificate bool
	scheme               string
	uRL                  string
	starHeader           string
	searchUrls           map[string]string
	answerHeader         string
	noAnswerMsg          string
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

	// format output
	starHeader = "\u2605"
	answerHeader = "%s Answer from  " + aurora.Green("%s").String() + "\n\n%s"
	noAnswerMsg = "< no answer given >"

	// cache
	cacheDir = getEnv("HOWDOI_CACHE_DIR", "")
}

// Howdoi string
func Howdoi(cli Cli) ([]string, error) {
	if cli.Debug {
		debug.Enable("*")
	}
	cli.Query = []string{strings.Replace(strings.Join(cli.Query[:], " "), "?", "", -1)}

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

		gLog(gree(fmt.Sprintf("0.2. questions: %d", len(questionLinks))))

		if searchEngine == "google" {
			questionLinks = cutURL(questionLinks)
		}
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
		// TODO: go func
		var wg sync.WaitGroup
		wg.Add(n)

		for i := 0; i < n; i++ { // the bigger one

			go func(i int) {
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
				wg.Done()
			}(i)
		}

		wg.Wait()

		gLog("2. answers: %v", string(len(answers)))

		return answers, nil

	}
	err = errors.New(aurora.Red("no questions link").String())

	return nil, err
}

func (clis Cli) getAnswer(u string) []string {
	doc := clis.getResult(u)
	return clis.extractAnswer(doc)
}

func (clis *Cli) extractAnswer(doc *goquery.Document) []string {
	gLog := debug.Debug("extractAnswer")

	links := make([]string, 0)
	tags := make([]string, 0)
	// get tag , use by chroma lexer
	getTags := doc.Find(".post-tag")
	if getTags.Size() > 0 {
		getTags.Each(func(i int, s *goquery.Selection) {
			str := s.Text()
			tags = append(tags, str)
		})
	}
	gLog("got post-tag: %v", tags)

	instructions := doc.Find(".answer").First().Find("pre")
	if instructions.Size() > 0 {
		instructions.Each(func(i int, s *goquery.Selection) {
			str := s.Text()
			str = colorCode(str, &ChromaColor{
				Color: clis.Color,
				Tags:  tags,
				Theme: clis.Theme}) // use chroma, colorful code
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
	vsf = UqineSlice(vsf)
	return vsf
}

func (clis Cli) getLinks() []string {

	searchURL := getSearchURL(searchEngine)
	u, _ := url.Parse(fmt.Sprintf(searchURL, clis.Query[0], uRL))

	q := u.Query()
	u.RawQuery = q.Encode() //urlencode
	doc := clis.getResult(u.String())
	return extractLinks(doc, searchEngine)
}

func (clis Cli) getResult(u string) *goquery.Document {
	gLog := debug.Debug("getResult")
	gLog("0. get URL:%v", u)

	var resp *http.Response
	var err error

	cacheHandle := CacheHowdoi{cacheDir} // Get Cache
	cacheBoby, ok := cacheHandle.cached(u)
	// TODO ? clis.ReCache
	if ok && !clis.ReCache {
		// resp from Cache
		gLog(gree("0. Resq from Cache"))

		r := bufio.NewReader(bytes.NewReader(cacheBoby))
		resp, err = http.ReadResponse(r, nil)
		if err != nil {
			log.Fatal(err)
		}
	} else { // GET URL
		gLog(red("ReCache:%v"), clis.ReCache)
		gLog(cyan("0. Resq from GET URL"))
		var req *http.Request

		// User-Agent random
		req, err = http.NewRequest("GET", u, nil)
		if err != nil {
			log.Fatalln(err)
		}
		req.Header.Set("User-Agent", userAgent.Random())

		proxyIs := whichProxy()
		if proxyIs == SOCKS {
			httpClient := Socks5Client()
			resp, err = httpClient.Do(req)
		} else {
			client := &http.Client{}
			resp, err = client.Do(req)
		}

		if err != nil {
			log.Fatalln(aurora.Red("请求失败:"+searchEngine), err)
		} else {
			defer resp.Body.Close()
		}

		if resp.StatusCode != 200 { // no 200, can no Cache
			log.Fatalln(aurora.Red("status code error:"), resp.Request.URL, resp.Status)
		}

		// Keep Cache
		if clis.Cache {
			body, err := httputil.DumpResponse(resp, clis.Cache)

			if err != nil {
				log.Fatal(err)
			}

			CacheResq(u, body, cacheDir)
		}
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalln(aurora.Red("goquery.NewDocumentFromReader error:"), err)
	}
	return doc
}
