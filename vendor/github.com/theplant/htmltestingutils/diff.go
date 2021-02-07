package htmltestingutils

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/andybalholm/cascadia"
	"github.com/pmezard/go-difflib/difflib"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"github.com/yosssi/gohtml"
)

func PrettyHtmlDiff(actual io.Reader, actualCssSelector string, expected string) (r string) {
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, actual)

	sel, err := cascadia.Compile(actualCssSelector)
	if err != nil {
		panic(err)
	}

	n, err := html.Parse(buf)
	if err != nil {
		panic(err)
	}
	mn := sel.MatchFirst(n)
	if mn == nil {
		panic(fmt.Sprintf("css selector '%s' not found in html:\n%s", actualCssSelector, buf.String()))
	}
	selBuf := bytes.NewBuffer(nil)
	html.Render(selBuf, mn)

	factual := trimLinesAndFormat(selBuf.String())
	fexpected := trimLinesAndFormat(expected)
	if fexpected != factual {
		diff := difflib.UnifiedDiff{
			A:        difflib.SplitLines(fexpected),
			B:        difflib.SplitLines(factual),
			FromFile: "Expected",
			ToFile:   "Actual",
			Context:  3,
		}
		r, _ = difflib.GetUnifiedDiffString(diff)
	}

	return
}

func trimLinesAndFormat(content string) string {
	trimmedBuf := bytes.NewBuffer(nil)
	lines := strings.Split(content, "\n")
	for _, l := range lines {
		trimmedBuf.WriteString(strings.TrimSpace(l) + "\n")
	}

	ns, err := html.ParseFragment(trimmedBuf, &html.Node{Data: "body", Type: html.ElementNode, DataAtom: atom.Body})
	if err != nil {
		panic(err)
	}

	renderBuf := bytes.NewBuffer(nil)
	for _, n := range ns {
		html.Render(renderBuf, n)
	}
	return gohtml.Format(renderBuf.String())
}
