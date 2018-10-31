package howdoi

import (
	"fmt"
	"os"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/logrusorgru/aurora"
	debug "github.com/visionmedia/go-debug"
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getMapDef(m map[string]string, key, fallback string) string {
	if value := m[key]; len(value) > 0 {
		return value
	}
	return fallback
}

func isRegexp(s string, reg string) bool {
	r := regexp.MustCompile(reg)

	m := r.MatchString(s)

	return m
}

func cutURL(links []string) []string {
	ls := make([]string, 0)

	for _, v := range links {
		if isRegexp(v, `^/url\?q=`) {
			ls = append(ls, v[7:])
		} else {
			ls = append(ls, v[:])
		}
	}
	return ls
}
func isQuestion(s string) bool {
	m := isRegexp(s, `questions/\d+/`)

	return m
}
func getSearchURL(s string) string {
	return getMapDef(searchUrls, s, searchUrls["bing"])
}
func extractLinks(doc *goquery.Document, engine string) []string {
	gLog := debug.Debug("extractLinks")

	var links []string
	if engine == "bing" {
		doc.Find(".b_algo h2 a").Each(func(i int, s *goquery.Selection) {
			attr, exists := s.Attr("href")
			if exists == true && isQuestion(attr) {
				links = append(links, attr)
			}
		})
	} else {
		one := doc.Find("a")
		if one.Size() > 0 {
			one.Each(func(i int, s *goquery.Selection) {
				attr, exists := s.Attr("href")

				if exists == true && isQuestion(attr) {
					links = append(links, attr)
				}
			})
		}
	}

	gLog("extract link %d", len(links))

	// Cache what you got
	// if len(links) == 0 {
	// 	s, _ := doc.Html()
	// 	gLog("page Hava text number %d", len(s))
	// 	f, _ := os.Create("./index.html")
	// 	n, _ := f.WriteString(s)

	// 	redLog(string(n))

	// }
	return links
}

// UqineSlice remove same string with slice
func UqineSlice(elements []string) []string {
	encountered := map[string]bool{}

	// Create a map of all unique elements.
	for v := range elements {
		encountered[elements[v]] = true
	}

	// Place all keys from the map into a slice.
	result := []string{}
	for key := range encountered {
		result = append(result, key)
	}
	return result
}

func redLog(s string) {
	fmt.Println(aurora.Red(s))
}

func red(s string) string {
	return aurora.Red(s).String()
}

func gree(s string) string {
	return aurora.Green(s).String()
}

func cyan(s string) string {
	return aurora.Cyan(s).String()
}

func format(f string, s ...interface{}) string {

	return fmt.Sprintf(f, s...)
}
