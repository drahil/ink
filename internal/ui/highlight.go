package ui

import (
	"bytes"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

type lineHighlighter func(string) string

func (e EditorPane) lineHighlighter() lineHighlighter {
	lexer := lexers.Match(e.Path)
	if lexer == nil {
		lexer = lexers.Analyse(e.Content)
	}
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)

	formatter := formatters.Get("terminal16m")
	if formatter == nil {
		formatter = formatters.Fallback
	}

	style := styles.Get("monokai")
	if style == nil {
		style = styles.Fallback
	}

	return func(line string) string {
		if line == "" {
			return ""
		}

		iterator, err := lexer.Tokenise(nil, line)
		if err != nil {
			return line
		}

		var b bytes.Buffer
		if err := formatter.Format(&b, style, iterator); err != nil {
			return line
		}

		return strings.TrimRight(b.String(), "\n")
	}
}
