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

	"github.com/PuerkitoBio/goquery"
	"github.com/logrusorgru/aurora"
	debug "github.com/visionmedia/go-debug"
)

func (clis Cli) getAnswer(u string) []string {
	doc, err := clis.getResult(u)
	if err != nil {
		return []string{}
	}
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

	// that action, do that for goole links
	vsf = cutURL(vsf)

	return vsf
}

func (clis Cli) getLinks() []string {
	gLog := debug.Debug("getLinks")

	var links []string

	linksChan := make(chan []string)

	// ALL engine or User select one
	var finalEngine map[string]string
	if searchEngine != allSearchEngine {
		finalEngine = map[string]string{
			searchEngine: getSearchURL(searchEngine)}
	} else {
		finalEngine = searchUrls
	}

	var winner string
	for engine, searchURL := range finalEngine {

		go func(searchURL, engine string) {
			u, _ := url.Parse(fmt.Sprintf(searchURL, clis.Query[0], uRL))
			q := u.Query()
			u.RawQuery = q.Encode() //urlencode
			doc, err := clis.getResult(u.String())

			if err == nil {
				linksChan <- extractLinks(doc, engine)
				winner = engine
			} else {
				linksChan <- nil
			}

		}(searchURL, engine)

	}

	rest := 0
	for i := 0; i < len(finalEngine); i++ {
		rest = i
		if res := <-linksChan; res != nil {
			links = res // just the most fasest and Right one
			gLog("this winner is %s", red(winner))
			break
		}
	}

	go func(rest int) {
		for i := 0; i < rest; i++ {
			<-linksChan // free the last chan
		}
		close(linksChan)
	}(len(finalEngine) - rest)

	return links
}

func (clis Cli) getResult(u string) (doc *goquery.Document, reqErr error) {
	gLog := debug.Debug("getResult")
	gLog("0. get URL:%v", u)

	defer func() { // hold the error msg
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				reqErr = errors.New(x)
			case error:
				reqErr = x
			default:
				reqErr = errors.New("Unknown panic")
			}
			// invalidate rep
			doc = nil
		}
	}()

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
			log.Panicln(err)
		}
	} else { // GET URL
		gLog(red("ReCache:%v"), clis.ReCache)
		gLog(cyan("0. Resq from GET URL"))
		var req *http.Request

		// User-Agent random
		req, err = http.NewRequest("GET", u, nil)
		if err != nil {
			log.Panicln(err)
		}
		req.Header.Set("User-Agent", getRandomUA())

		proxyIs := whichProxy()
		if proxyIs == SOCKS {
			httpClient := Socks5Client()
			resp, err = httpClient.Do(req)
		} else {
			client := &http.Client{}
			resp, err = client.Do(req)
		}

		if err != nil {
			log.Panicln(aurora.Red("请求失败:"), err)
		} else {
			defer resp.Body.Close()
		}

		if resp.StatusCode != 200 { // no 200, can no Cache
			log.Panicln(aurora.Red("status code error:"), resp.Request.URL, resp.Status)
		}

		// Keep Cache
		if clis.Cache {
			body, err := httputil.DumpResponse(resp, clis.Cache)

			if err != nil {
				log.Panicln(err)
			}

			CacheResq(u, body, cacheDir)
		}
	}

	doc, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Panicln(aurora.Red("goquery.NewDocumentFromReader error:"), err)
	}

	return
}
