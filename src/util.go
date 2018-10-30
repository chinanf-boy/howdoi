package howdoi

import (
	"fmt"
	"os"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/logrusorgru/aurora"
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
