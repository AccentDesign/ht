# HTML Node Builder in Go

This project provides a set of functions to create and manipulate HTML nodes programmatically using Go. It simplifies the process of generating HTML documents by offering a clean and intuitive API.

## Features

- Generate HTML elements with ease.
- Support for attributes, text, and raw HTML.
- Render HTML nodes to standard output or files.

## Installation

To use this library, add it to your Go project:

```bash
go get github.com/accentdesign/ht
```

## Usage

Here is an example of how to use the library:

```go
package main

import (
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
				Script(Src("main.js")),
				Link(Rel("stylesheet"), Href("style.css")),
			),
			Body(Class("body")),
		),
	)
	_ = html.Render(os.Stdout, node)
}
```

## Note:

Some attribute helpers are suffixed with `Attr` (e.g., LabelAttr, StyleAttr, TitleAttr)
to avoid naming conflicts with element constructors (Label, Style, Title).
Use these helpers to set the corresponding attribute on an element.