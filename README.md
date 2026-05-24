# ht — An HTML AST Builder for Go

`ht` is a hyper-minimal, functional HTML builder for Go. It provides a clean, type-safe API for constructing standard `golang.org/x/net/html` AST nodes.

> [!WARNING]
> **DO NOT `go get` THIS PACKAGE!**
> 
> This is not a dependency to add to your `go.mod`. It is a minimal (~300 lines) foundation designed to be **copied, stolen, and modified**. 
> 
> Copy `nodes.go` and `attrs.go` directly into your project. Tweak them. Add your own project-specific helpers (like custom SVG components or specific JavaScript bindings). **Own your HTML builder.**

## Why use this over `html/template`?

1. **Functional Components**: Build HTML exactly how you build Go code. Components are just standard Go functions that return `*html.Node`. No more wrestling with template context or inheritance hierarchies.
2. **Standard Library AST**: This isn't a custom virtual DOM. It constructs the exact same `*html.Node` AST that the Go standard library uses to parse HTML. `html.Render` simply serializes it.
3. **Resilient Rendering**: In standard templates, an error halfway through rendering crashes the HTTP response and leaves broken HTML. With `ht`, AST building and network serialization are separate. If a component fails to build, you can catch the error and safely return an "Error Node" (like a fallback UI) without breaking the rest of your layout.
4. **First-Class HTMX & Alpine Support**: Easily create hypermedia-driven SPAs with included attribute helpers (`HxPost`, `HxSwapOob`, `XData`, `XOn`) that feel native to Go.

## Examples

We provide some examples in the `examples/` directory to show how to use `ht`:

1. **`search/`**: A basic movie search engine using HTMX, SSE, and DaisyUI. Shows how to use partials, out-of-band swaps (`hx-swap-oob`), and local component state.
2. **`sysmon/`**: A basic, real-time fake system monitor using HTMX and SSE.
3. **`template/`**: Demonstrates how to seamlessly interop with `text/template` or `html/template` code. Shows how to safely execute templates into `Raw` nodes and embed them directly into an `ht` AST.
4. **`todo/`**: A fully functional, highly interactive Todo application utilizing HTMX, Alpine.js, and DaisyUI. Shows how to use partials, out-of-band swaps (`hx-swap-oob`), and local component state.

## Usage

```go
package main

import (
    "os"
    "golang.org/x/net/html"
    
    // Import your local, copied package
    . "yourproject/internal/ht" 
)

func main() {
    page := Document(
        Doctype("html"),
        Html(
            Lang("en"),
            Head(
                Title(Text("My App")),
                Script(Src("https://unpkg.com/htmx.org@2.0.4")),
            ),
            Body(
                Class("bg-base-300"),
                H1(Class("text-2xl font-bold"), Text("Hello, World")),
                Button(
                    Class("btn btn-primary"),
                    HxPost("/clicked"),
                    Text("Click Me"),
                ),
            ),
        ),
    )
    
    _ = html.Render(os.Stdout, page)
}
```

## Important Notes

- **Naming Conflicts**: Some attribute helpers are suffixed with `Attr` (e.g., `LabelAttr`, `StyleAttr`, `TitleAttr`) to avoid naming conflicts with the HTML element constructors (`Label`, `Style`, `Title`).
- **Raw HTML**: Use `Raw("<br>")` to inject unescaped HTML strings. Only pass trusted content to `Raw`. Use `Text("Hello")` for regular strings; it will be automatically escaped by the renderer.
- **Node Detachment**: If you pass an existing `*html.Node` as a child, it is appended using standard `node.AppendChild` semantics. The child node MUST be detached (`Parent == nil`, `PrevSibling == nil`, `NextSibling == nil`) or the Go standard library will panic. 