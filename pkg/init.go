package howdoi

import (
	"os"

	"github.com/logrusorgru/aurora"
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
	uRLRegexText         string
)

const allSearchEngine = "ALL"

func init() {
	if len(os.Getenv("HOWDOI_DISABLE_SSL")) > 0 {
		scheme = "http://"
		verifySslCertificate = false
	} else {
		scheme = "https://"
		verifySslCertificate = true
	}
	uRL = getEnv("HOWDOI_URL", "stackoverflow.com")
	searchEngine = getEnv("HOWDOI_SEARCH_ENGINE", allSearchEngine)
	uRLRegexText = getEnv("HOWDOI_URL_REGEX", `questions/\d+/`)

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
