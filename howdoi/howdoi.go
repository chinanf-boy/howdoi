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

	userAgents = []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.7; rv:11.0) Gecko/20100101 Firefox/11.0",
	}

	searchUrls = map[string]string{
		"bing":   scheme + "www.bing.com/search?q=%s site:%s",
		"google": scheme + "www.google.com/search?q=%s site:%s",
	}
	starHeader = "\u2605"
	answerHeader = "{2}  Answer from {0} {2}\n\n{1}"
	noAnswerMsg = "< no answer given >"
	searchEngine = getEnv("HOWDOI_SEARCH_ENGINE", "google")
}

// Howdoi string
func Howdoi(res Cli) ([]string, error) {

	res.Query = []string{strings.Replace(strings.Join(res.Query[:], " "), "?", "", -1)}

	result, err := res.getInstructions()

	if err != nil {
		return nil, err
	}
	fmt.Printf("%v", result)
	// userAgents[0] + searchUrls["bing"]
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
		for _, k := range questionLinks {
			fmt.Println(k)
		}
		n := clis.Num
		for n > 0 {
			// HERE
			n--
		}
	} else {
		err = errors.New("no questions link")
	}

	return []string{}, err
}

func cutURL(links []string) []string {
	ls := make([]string, 0)

	for _, v := range links {
		if isRegexp(v, `^/url\?q=`) {
			ls = append(ls, v[7:])
		}
	}
	return ls
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

func isQuestion(s string) bool {
	m := isRegexp(s, `questions/\d+/`)

	return m
}

func (clis Cli) getLinks() []string {

	searchURL := getSearchURL(searchEngine)
	u, _ := url.Parse(fmt.Sprintf(searchURL, clis.Query[0], uRL))

	q := u.Query()
	u.RawQuery = q.Encode() //urlencode
	doc, err := getResult(u.String())

	if err != nil {

	}
	return extractLinks(doc, searchEngine)
}

func getSearchURL(s string) string {
	return getMapDef(searchUrls, s, searchUrls["bing"])
}

func extractLinks(doc *goquery.Document, engine string) []string {
	var links []string
	if engine == "bing" {
		doc.Find(".b_algo h2 a").Each(func(i int, s *goquery.Selection) {
			attr, exists := s.Attr("href")
			if exists == true {
				links = append(links, attr)
			}
		})
	} else {
		// doc.Find(".l").Each(func(i int, s *goquery.Selection) {
		// 	attr, exists := s.Attr("href")
		// 	if exists == true {
		// 		links = append(links, attr)
		// 	}
		// })
		doc.Find(".r a").Each(func(i int, s *goquery.Selection) {
			attr, exists := s.Attr("href")
			if exists == true {
				links = append(links, attr)
			}
		})
	}

	return links
}
func getResult(u string) (*goquery.Document, error) {
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
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	} else {
		doc, err := goquery.NewDocumentFromReader(res.Body)
		return doc, err
	}
	return nil, err
}
