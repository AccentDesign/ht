package main

import (
	"errors"
	"fmt"
	"os"

	. "github.com/accentdesign/ht"
	"golang.org/x/net/html"
)

func main() {
	node := Document(
		Doctype("html"),
		Html(
			Lang("en"),
			Head(
				Meta(Charset("utf-8")),
				Meta(Name("viewport"), Content("width=device-width", "initial-scale=1.0")),
				Title(Text("Page")),
				Script(Src("main.js"), Defer()),
				Link(Rel("stylesheet"), Href("style.css")),
			),
			Body(Class("body")),
		),
	)
	printNode(node)

	node = Div(Class("container"), P(Text("Hello, World!")))
	printNode(node)

	node = Fieldset(Class("fieldset"),
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

	node = Nav(
		Role("navigation"),
		Class("nav"),
		A(Href("/"), Text("Home")),
		A(Href("/about"), Text("About")),
	)
	printNode(node)

	node = Ul(
		Role("menu"),
		Class("menu"),
		Data("menu", "main"),
		Li(Text("Item 1")),
		Li(Text("Item 2")),
		Li(Text("Item 3")),
	)
	printNode(node)

	node = Div(If(true, Text("True")), If(false, Text("False")))
	printNode(node)

	node = Div(Class("card", "bg-base-100", "w-96", "shadow-sm"),
		Figure(Img(Src("https://img.daisyui.com/images/stock/photo-1606107557195-0e29a4b5b4aa.webp"), Alt("Card Image"))),
		Div(Class("card-body"),
			H2(Class("card-title"), Text("Card Title")),
			P(Text("A card component has a figure, a body part, and inside body there are title and actions parts")),
			Div(Class("card-actions justify-end"), Button(Class("btn", "btn-primary"), Text("Buy Now"))),
		),
	)
	printNode(node)

	str := "hello world"
	node = P(str)
	printNode(node)

	node = P(&str)
	printNode(node)

	node = P(1)
	printNode(node)

	node = P(true)
	printNode(node)

	node = P(errors.New("error message"))
	printNode(node)
}

func printNode(node *html.Node) {
	if err := html.Render(os.Stdout, node); err != nil {
		fmt.Println("Error rendering HTML:", err)
		return
	}
	fmt.Println()
}
