package main

import (
	"fmt"
	. "github.com/accentdesign/ht"
	"golang.org/x/net/html"
	"os"
)

func main() {
	node := Document(
		Doctype("html"),
		Html(
			Lang("en"),
			Head(
				Meta(Charset("utf-8")),
				Meta(Name("viewport"), Content("width=device-width, initial-scale=1.0")),
				Title(Text("Page")),
				Script(Src("main.js")),
				Link(Rel("stylesheet"), Href("style.css")),
			),
			Body(Class("body")),
		),
	)
	printNode(node)

	node = Div(Class("container"), P(Text("Hello, World!")))
	printNode(node)

	node = Div(
		Class("field"),
		Label(Class("label"), For("name"), Text("Name")),
		Input(Class("input"), Id("name"), Type("text"), Name("name"), Placeholder("name"), Value("John Doe")),
		P(Class("help"), Text("Please enter your name.")),
	)
	printNode(node)

	node = Comment("some awesome comment")
	printNode(node)

	node = Text("<h1>Hello, World!</h1>")
	printNode(node)

	node = Raw("<h1>Hello, World!</h1>")
	printNode(node)

	node = Div(
		H1(Text("Header 1")),
		H2(Text("Header 2")),
		H3(Text("Header 3")),
		H4(Text("Header 4")),
		H5(Text("Header 5")),
	)
	printNode(node)
}

func printNode(node *html.Node) {
	if err := html.Render(os.Stdout, node); err != nil {
		fmt.Println("Error rendering HTML:", err)
		return
	}
	fmt.Println()
}
