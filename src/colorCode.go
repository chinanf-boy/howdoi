package howdoi

import (
	"bytes"
	"log"

	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	debug "github.com/visionmedia/go-debug"
)

func colorCode(s string, data Cli) string {
	gLog := debug.Debug("colorCode")

	if data.Color == false {
		gLog("no need color")
		return s
	}
	res := new(bytes.Buffer)

	// lexer
	lexer := lexers.Get("bash")
	for _, v := range data.Tags {
		lexer = lexers.Get(v)
		if lexer != nil {
			gLog("lexer : %s", v)
			break
		}
	}
	if lexer == nil {
		gLog("lexer : %s", "noop")
		lexer = lexers.Fallback
	}

	// styles
	style := styles.Get(data.Theme)
	if style == nil {
		style = styles.Fallback
	}
	gLog("style : %s", style.Name)

	// formatter
	// [html json noop terminal terminal16m terminal256 tokens]
	f := "terminal"
	formatter := formatters.Get(f)
	if formatter == nil {
		formatter = formatters.Fallback
	}
	gLog("formatter : %s", f)

	iterator, err := lexer.Tokenise(nil, s)
	// done
	err = formatter.Format(res, style, iterator)

	if err != nil {
		log.Fatalln(err)
	}

	return res.String()
}
