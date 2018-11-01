package howdoi

import (
	"bytes"
	"log"

	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

func colorCode(s string, data Cli) string {
	if data.Color == false {
		return s
	}
	res := new(bytes.Buffer)
	lexer := lexers.Get("bash")
	if lexer == nil {
		lexer = lexers.Fallback
	}

	style := styles.Get(data.Theme)
	if style == nil {
		style = styles.Fallback
	}

	// [html json noop terminal terminal16m terminal256 tokens]
	formatter := formatters.Get("terminal")
	if formatter == nil {
		formatter = formatters.Fallback
	}
	iterator, err := lexer.Tokenise(nil, s)

	err = formatter.Format(res, style, iterator)

	if err != nil {
		log.Fatalln(err)
	}

	return res.String()
}
