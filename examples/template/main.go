package main

import (
	"fmt"
	"html/template"
	"os"
	"strings"

	. "github.com/accentdesign/ht"
	"golang.org/x/net/html"
)

func main() {
	// 1. Define a standard Go text/template (e.g., a legacy component)
	tpl := template.Must(template.New("userCard").Parse(`
<div class="card" style="border: 1px solid #ccc; padding: 1rem; margin-top: 1rem;">
	<h3>{{.Name}}</h3>
	<p>Role: <strong>{{.Role}}</strong></p>
</div>`))

	// 2. Build your AST natively with ht, embedding the template dynamically
	page := Document(
		Doctype("html"),
		Html(
			Lang("en"),
			Head(Title(Text("Template Interop Example"))),
			Body(
				H1(Text("Integrating text/template")),
				P(Text("Below is a template injected directly into the AST:")),

				// Call the helper to inject the template execution
				FromTemplate(tpl, "userCard", map[string]any{
					"Name": "Alice",
					"Role": "Administrator",
				}),

				FromTemplate(tpl, "userCard", map[string]any{
					"Name": "Bob",
					"Role": "Editor",
				}),
			),
		),
	)

	// 3. Render the final unified tree
	if err := html.Render(os.Stdout, page); err != nil {
		panic(err)
	}
}

// FromTemplate is a clean helper to execute a template and return it as a Raw HTML node.
func FromTemplate(tpl *template.Template, name string, data any) *html.Node {
	var b strings.Builder

	// Pre-allocate a reasonable buffer size to minimize allocations
	b.Grow(256)

	var err error
	if name == "" {
		err = tpl.Execute(&b, data)
	} else {
		err = tpl.ExecuteTemplate(&b, name, data)
	}

	if err != nil {
		// Return the error safely as a text node so it doesn't break rendering
		return Div(
			Class("error text-red-500"),
			Text(fmt.Sprintf("template error: %v", err)),
		)
	}

	// Wrap the output in a Raw node so it isn't HTML-escaped by the parent tree
	return Raw(b.String())
}
