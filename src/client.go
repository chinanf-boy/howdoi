package howdoi

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"golang.org/x/net/proxy"
)

const (
	// SOCKS proxy got
	SOCKS = "socks"
	// HTTP proxy got
	HTTP = "http"
)

// Client get/post/...
type Client interface {
}

// Config env proxy
type Config struct {
	HTTPProxy  string
	HTTPSProxy string
	NoProxy    string
	CGI        bool
	ALLProxy   string
}

var config = Config{
	HTTPProxy:  getEnvAny("HTTP_PROXY", "http_proxy"),
	HTTPSProxy: getEnvAny("HTTPS_PROXY", "https_proxy"),
	NoProxy:    getEnvAny("NO_PROXY", "no_proxy"),
	CGI:        os.Getenv("REQUEST_METHOD") != "",
	ALLProxy:   getEnvAny("ALL_PROXY", "all_proxy"),
}

// Socks5Client > get httpClient with socks5
func Socks5Client() *http.Client {
	u, e := getURL()
	if e != nil {
		log.Fatalln("Proxy Env Set Error", e)
	}
	dialer, err := proxy.SOCKS5("tcp", fmt.Sprintf("%s:%s", u.Hostname(), u.Port()),
		nil,
		&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 10 * time.Second,
		},
	)
	if err != nil {
		log.Fatalln("get dialer error", dialer)
	}

	httpTransport := &http.Transport{Dial: dialer.Dial}
	httpClient := &http.Client{Transport: httpTransport}
	return httpClient
}

// HTTPClient > get httpClient with socks5

// GetProxis get http/s_proxy
func GetProxis() Config {
	return config
}

func getURL() (*url.URL, error) {
	var u *url.URL
	var e error
	if len(config.ALLProxy) > 0 {
		u, e = url.Parse(config.ALLProxy)
	} else if len(config.HTTPSProxy) > 0 {
		u, e = url.Parse(config.HTTPSProxy)

	} else if len(config.HTTPProxy) > 0 {
		u, e = url.Parse(config.HTTPProxy)
	}

	return u, e
}

func whichProxy() string {
	var result string
	if isSocks(config.ALLProxy) ||
		isSocks(config.HTTPProxy) ||
		isSocks(config.HTTPSProxy) {
		result = SOCKS
	} else {
		result = HTTP
	}
	return result
}

func isSocks(s string) bool {
	b := isRegexp(s, `^socks`)

	return b
}

func getEnvAny(names ...string) string {
	for _, n := range names {
		if val := os.Getenv(n); val != "" {
			return val
		}
	}
	return ""
}
