package howdoi

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
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
}

// Howdoi string
func Howdoi(res Cli) (string, error) {

	res.Query = []string{strings.Replace(strings.Join(res.Query[:], " "), "?", "", -1)}

	result, err := res.getInstructions()

	if err != nil {
		return "", err
	}
	fmt.Printf("%v", result)
	// userAgents[0] + searchUrls["bing"]
	return result, nil
}

func (clis Cli) getInstructions() (string, error) {
	var result string
	var err error
	links := clis.getLinks() // HERE
	for _, k := range links {
		fmt.Println(k)
	}

	return result, err
}

func (clis Cli) getLinks() []string {

	searchEngine := getEnv("HOWDOI_SEARCH_ENGINE", "bing")
	searchURL := getSearchURL(searchEngine)
	u, _ := url.Parse(fmt.Sprintf(searchURL, clis.Query[0], uRL))
	q := u.Query()
	u.RawQuery = q.Encode() //urlencode
	doc, engine := getResult(u.String(), searchEngine)
	return extractLinks(doc, engine)
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
		doc.Find(".l").Each(func(i int, s *goquery.Selection) {
			attr, exists := s.Attr("href")
			if exists == true {
				links = append(links, attr)
			}
		})
		doc.Find(".r a").Each(func(i int, s *goquery.Selection) {
			attr, exists := s.Attr("href")
			if exists == true {
				links = append(links, attr)
			}
		})
	}

	return links
}
func getResult(url string, engine string) (*goquery.Document, string) {

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return doc, engine
}
