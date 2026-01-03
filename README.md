# HTML Node Builder in Go

This project provides a set of functions to create and manipulate HTML nodes programmatically using Go. It simplifies the process of generating HTML documents by offering a clean and intuitive API.

## Features

- Generate HTML elements with ease.
- Support for attributes, text, and raw HTML.
- Render HTML nodes to standard output or files.

## How to use

This package is intended to be copied directly into your project (not installed via `go get`).
Copy the `ht` folder into your repository and import it locally.

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
                Script(Src("main.js"), Defer()),
                Link(Rel("stylesheet"), Href("style.css")),
            ),
            Body(Class("body")),
        ),
    )
    _ = html.Render(os.Stdout, node)
}
```

## Notes

Some attribute helpers are suffixed with `Attr` (e.g., LabelAttr, StyleAttr, TitleAttr)
to avoid naming conflicts with element constructors (Label, Style, Title).
Use these helpers to set the corresponding attribute on an element.

`Raw` injects unescaped HTML. Only pass trusted content to `Raw`. Use `Text` for normal text; it will be escaped by the renderer.

### Contract for passing *html.Node as children

When you pass an existing `*html.Node` (from `golang.org/x/net/html`) as an argument to `Element(...)` (or any element helper like `Div(...)`, `Span(...)`, etc.), it is appended using `node.AppendChild` semantics. That means the child node MUST be detached before you pass it in:

- The child must have `Parent == nil`, `PrevSibling == nil`, and `NextSibling == nil`.
- If it already belongs to another tree, detach it first, e.g. `child.Parent.RemoveChild(child)`, then pass it.
- Alternatively, clone the node if you want to keep the original where it is.

If you pass an attached node, `AppendChild` will panic (this is behavior from the `html` package).