package howdoi

import (
	"os"
	"testing"
)

func TestCientUtil(t *testing.T) {
	var config = Config{
		HTTPProxy:  getEnvAny("HTTP_PROXY", "http_proxy"),
		HTTPSProxy: getEnvAny("HTTPS_PROXY", "https_proxy"),
		NoProxy:    getEnvAny("NO_PROXY", "no_proxy"),
		CGI:        os.Getenv("REQUEST_METHOD") != "",
		ALLProxy:   getEnvAny("ALL_PROXY", "all_proxy"),
	}
	a := "NOTHINGNOTHINGNOTHING"
	res := getEnvAny(a, a)

	if res != "" {
		t.Fail()
	}
	_, err := getURL(config)

	if err != nil {
		t.Fail()
	}
}
