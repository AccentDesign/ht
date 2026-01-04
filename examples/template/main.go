package main

import (
	"bytes"
	"html/template"
	"os"

	. "github.com/accentdesign/ht"
	"golang.org/x/net/html"
)

func main() {
	partial := template.Must(template.New("partial").Parse(`<strong>Hello, {{.Name}}</strong>`))

	node := Div(
		Class("prose"),
		P("Before partial"),
		RawFromExec(func(buf *bytes.Buffer) error {
			return partial.Execute(buf, map[string]any{"Name": "Gopher"})
		}),
		P("After partial"),
	)

	if err := html.Render(os.Stdout, node); err != nil {
		panic(err)
	}
}

func RawFromExec(runner func(*bytes.Buffer) error) *html.Node {
	var b bytes.Buffer
	if err := runner(&b); err != nil {
		return Text(err.Error())
	}
	return Raw(b.String())
}
