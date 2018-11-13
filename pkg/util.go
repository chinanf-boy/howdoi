package howdoi

import (
	"bytes"
	"fmt"
	"os"
	"regexp"

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
	if len(ls) > 0 {
		return ls
	}
	return links
}
func isQuestion(s string) bool {
	m := isRegexp(s, `questions/\d+/`)

	return m
}
func getSearchURL(s string) string {
	return getMapDef(searchUrls, s, searchUrls["bing"])
}

// UqineSlice remove same string with slice
func UqineSlice(elements []string) []string {
	ret := elements[:0]
	// 利用 struct{}{} 减少内存占用
	assist := map[string]struct{}{}
	for _, v := range elements {
		if _, ok := assist[v]; !ok {
			assist[v] = struct{}{}
			ret = append(ret, v)
		}
	}
	return ret
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

func sliceGoodFmt(arr []string) string {
	var buf bytes.Buffer

	for _, v := range arr {
		buf.WriteString(fmt.Sprintf("%v\n", v))
	}

	return buf.String()
}

// func structGoodFmt(strt interface{}) string {
// 	// gLog := debug.Debug("structGoodFmt")

// 	v := reflect.ValueOf(strt).Elem()
// 	ks := reflect.TypeOf(strt).Elem()

// 	var buf bytes.Buffer
// 	count := v.NumField()

// 	for i := 0; i < count; i++ {
// 		f := v.Field(i)
// 		k := ks.Field(i)

// 		buf.WriteString(fmt.Sprintf("%v : %v\n", k, f))
// 	}

// 	return buf.String()
// }
