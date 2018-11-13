package howdoi

import (
	"bytes"
	"log"
	"os/user"
	"path/filepath"

	"github.com/gregjones/httpcache/diskcache"
	debug "github.com/visionmedia/go-debug"
)

var (
	cacheDir string // HOWDOI_CACHE_DIR
)

// CacheHowdoi keep cacheDir
type CacheHowdoi struct {
	cacheDir string
}

// CacheResq cache Resq
func CacheResq(key string, value []byte, dir string) {
	gLog := debug.Debug("CacheResq")
	gLog("Cache URL:%v %d", key, len(value))
	
	C := CacheHowdoi{dir}
	
	_, ok := C.cached(key)

	if ok {
		return
	}

	C.cacheKey(key, value)

}

func (c CacheHowdoi) cached(key string) ([]byte, bool) {
	absDir := c.cacheAbsDir()

	cache := diskcache.New(absDir)

	body, ok := cache.Get(key)

	return body, ok
}

func (c CacheHowdoi) cacheKey(key string, value []byte) {
	absDir := c.cacheAbsDir()

	cache := diskcache.New(absDir)

	cache.Set(key, value)

	retVal, ok := cache.Get(key)
	if !ok {
		log.Fatal("could not retrieve an element we just added")
	}
	if !bytes.Equal(retVal, value) {
		log.Fatal("retrieved a different resque than what we put in")
	}
}

func (c CacheHowdoi) cacheAbsDir() string {
	dir := c.cacheDir
	if dir == "" {
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		h := usr.HomeDir
		dir = filepath.Join(h, ".howdoi-cache")
	}

	C, err := filepath.Abs(dir)

	if err != nil {
		log.Fatalf("cacheDir no Abs / some: %v", err)
	}

	if C != dir {
		c.setDir(C)
	}

	return C
}

func (c *CacheHowdoi) setDir(s string) {
	c.cacheDir = s
}
