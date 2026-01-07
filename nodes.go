package ht

import (
	"fmt"
	"strings"

	h "golang.org/x/net/html"
	a "golang.org/x/net/html/atom"
)

type delimiter struct {
	split string
	join  string
}

var mergeDelims = map[string]delimiter{
	"class":   {split: " ", join: " "},
	"content": {split: ",", join: ", "},
}

// mergeAttr combines two attribute strings (oldVal and newVal) using a delimiter, ensuring unique,
// trimmed components.
func mergeAttr(oldVal, newVal string, d delimiter) string {
	if oldVal == "" {
		return strings.TrimSpace(newVal)
	}
	if newVal == "" {
		return strings.TrimSpace(oldVal)
	}

	f := func(r rune) bool {
		return strings.ContainsRune(d.split, r)
	}

	fields := strings.FieldsFunc(oldVal+d.join+newVal, f)
	seen := make(map[string]bool, len(fields))
	parts := make([]string, 0, len(fields))

	for _, part := range fields {
		t := strings.TrimSpace(part)
		if t != "" && !seen[t] {
			seen[t] = true
			parts = append(parts, t)
		}
	}

	return strings.Join(parts, d.join)
}

// Comment creates and returns a new comment node with the provided data.
func Comment(data string) *h.Node {
	return &h.Node{Type: h.CommentNode, Data: data}
}

// Doctype creates a new node representing a document type with the specified data.
func Doctype(data string) *h.Node {
	return &h.Node{Type: h.DoctypeNode, Data: data}
}

// Document creates a new document node and appends the provided child nodes to it.
func Document(children ...*h.Node) *h.Node {
	node := &h.Node{Type: h.DocumentNode}
	for _, child := range children {
		node.AppendChild(child)
	}
	return node
}

// Element constructs an HTML element node with the given tag and variadic args.
//
// Supported arg types:
//   - h.Attribute: added or merged into the element's attributes. When a key
//     has a join strategy (see mergeAttrMap), values are combined instead of replaced.
//   - *h.Node: appended as a child of the created element.
//     CONTRACT: The passed node must be detached — i.e. n.Parent == nil,
//     n.PrevSibling == nil, and n.NextSibling == nil. This mirrors
//     golang.org/x/net/html.Node.AppendChild, which will panic if the child
//     already has a parent or siblings. Detach the node from its current
//     parent (e.g. parent.RemoveChild(n)) before passing it here, or clone it
//     if you need to keep the original in place.
//   - string, *string, fmt.Stringer, error, or any other type: coerced to text
//     via Text(...).
func Element(tag a.Atom, args ...any) *h.Node {
	node := &h.Node{Type: h.ElementNode, DataAtom: tag, Data: tag.String()}
	for _, arg := range args {
		switch v := arg.(type) {
		case *h.Node:
			if v != nil {
				node.AppendChild(v)
			}
		case h.Attribute:
			found := false
			for i, attr := range node.Attr {
				if attr.Key == v.Key {
					if d, ok := mergeDelims[attr.Key]; ok {
						node.Attr[i].Val = mergeAttr(attr.Val, v.Val, d)
					} else {
						node.Attr[i].Val = v.Val
					}
					found = true
					break
				}
			}
			if !found {
				node.Attr = append(node.Attr, v)
			}
		case string:
			node.AppendChild(Text(v))
		case *string:
			if v != nil {
				node.AppendChild(Text(*v))
			}
		case fmt.Stringer:
			if v != nil {
				node.AppendChild(Text(v.String()))
			}
		case error:
			if v != nil {
				node.AppendChild(Text(v.Error()))
			}
		default:
			// Coerce any other types to string content without side effects.
			node.AppendChild(Text(fmt.Sprint(v)))
		}
	}
	return node
}

// Raw creates a node with raw HTML content, bypassing any HTML escaping for the supplied input string.
func Raw(data string) *h.Node {
	return &h.Node{Type: h.RawNode, Data: data}
}

// Text creates a text node with the specified string content and returns a pointer to the node.
func Text(data string) *h.Node {
	return &h.Node{Type: h.TextNode, Data: data}
}

func A(args ...any) *h.Node          { return Element(a.A, args...) }
func Abbr(args ...any) *h.Node       { return Element(a.Abbr, args...) }
func Address(args ...any) *h.Node    { return Element(a.Address, args...) }
func Article(args ...any) *h.Node    { return Element(a.Article, args...) }
func Aside(args ...any) *h.Node      { return Element(a.Aside, args...) }
func B(args ...any) *h.Node          { return Element(a.B, args...) }
func Blockquote(args ...any) *h.Node { return Element(a.Blockquote, args...) }
func Body(args ...any) *h.Node       { return Element(a.Body, args...) }
func Button(args ...any) *h.Node     { return Element(a.Button, args...) }
func Br(args ...any) *h.Node         { return Element(a.Br, args...) }
func Caption(args ...any) *h.Node    { return Element(a.Caption, args...) }
func Cite(args ...any) *h.Node       { return Element(a.Cite, args...) }
func Code(args ...any) *h.Node       { return Element(a.Code, args...) }
func Col(args ...any) *h.Node        { return Element(a.Col, args...) }
func Colgroup(args ...any) *h.Node   { return Element(a.Colgroup, args...) }
func Dd(args ...any) *h.Node         { return Element(a.Dd, args...) }
func Details(args ...any) *h.Node    { return Element(a.Details, args...) }
func Dialog(args ...any) *h.Node     { return Element(a.Dialog, args...) }
func Div(args ...any) *h.Node        { return Element(a.Div, args...) }
func Dl(args ...any) *h.Node         { return Element(a.Dl, args...) }
func Dt(args ...any) *h.Node         { return Element(a.Dt, args...) }
func Em(args ...any) *h.Node         { return Element(a.Em, args...) }
func Fieldset(args ...any) *h.Node   { return Element(a.Fieldset, args...) }
func Figcaption(args ...any) *h.Node { return Element(a.Figcaption, args...) }
func Figure(args ...any) *h.Node     { return Element(a.Figure, args...) }
func Footer(args ...any) *h.Node     { return Element(a.Footer, args...) }
func Form(args ...any) *h.Node       { return Element(a.Form, args...) }
func H1(args ...any) *h.Node         { return Element(a.H1, args...) }
func H2(args ...any) *h.Node         { return Element(a.H2, args...) }
func H3(args ...any) *h.Node         { return Element(a.H3, args...) }
func H4(args ...any) *h.Node         { return Element(a.H4, args...) }
func H5(args ...any) *h.Node         { return Element(a.H5, args...) }
func Head(args ...any) *h.Node       { return Element(a.Head, args...) }
func Header(args ...any) *h.Node     { return Element(a.Header, args...) }
func Hr(args ...any) *h.Node         { return Element(a.Hr, args...) }
func Html(args ...any) *h.Node       { return Element(a.Html, args...) }
func I(args ...any) *h.Node          { return Element(a.I, args...) }
func Img(args ...any) *h.Node        { return Element(a.Img, args...) }
func Input(args ...any) *h.Node      { return Element(a.Input, args...) }
func Label(args ...any) *h.Node      { return Element(a.Label, args...) }
func Legend(args ...any) *h.Node     { return Element(a.Legend, args...) }
func Li(args ...any) *h.Node         { return Element(a.Li, args...) }
func Link(args ...any) *h.Node       { return Element(a.Link, args...) }
func Main(args ...any) *h.Node       { return Element(a.Main, args...) }
func Mark(args ...any) *h.Node       { return Element(a.Mark, args...) }
func Meta(args ...any) *h.Node       { return Element(a.Meta, args...) }
func Nav(args ...any) *h.Node        { return Element(a.Nav, args...) }
func Ol(args ...any) *h.Node         { return Element(a.Ol, args...) }
func Optgroup(args ...any) *h.Node   { return Element(a.Optgroup, args...) }
func Option(args ...any) *h.Node     { return Element(a.Option, args...) }
func P(args ...any) *h.Node          { return Element(a.P, args...) }
func Pre(args ...any) *h.Node        { return Element(a.Pre, args...) }
func Script(args ...any) *h.Node     { return Element(a.Script, args...) }
func Section(args ...any) *h.Node    { return Element(a.Section, args...) }
func Select(args ...any) *h.Node     { return Element(a.Select, args...) }
func Small(args ...any) *h.Node      { return Element(a.Small, args...) }
func Span(args ...any) *h.Node       { return Element(a.Span, args...) }
func Strong(args ...any) *h.Node     { return Element(a.Strong, args...) }
func Style(args ...any) *h.Node      { return Element(a.Style, args...) }
func Sub(args ...any) *h.Node        { return Element(a.Sub, args...) }
func Summary(args ...any) *h.Node    { return Element(a.Summary, args...) }
func Sup(args ...any) *h.Node        { return Element(a.Sup, args...) }
func Table(args ...any) *h.Node      { return Element(a.Table, args...) }
func Tbody(args ...any) *h.Node      { return Element(a.Tbody, args...) }
func Td(args ...any) *h.Node         { return Element(a.Td, args...) }
func Textarea(args ...any) *h.Node   { return Element(a.Textarea, args...) }
func Tfoot(args ...any) *h.Node      { return Element(a.Tfoot, args...) }
func Th(args ...any) *h.Node         { return Element(a.Th, args...) }
func Thead(args ...any) *h.Node      { return Element(a.Thead, args...) }
func Title(args ...any) *h.Node      { return Element(a.Title, args...) }
func Tr(args ...any) *h.Node         { return Element(a.Tr, args...) }
func Ul(args ...any) *h.Node         { return Element(a.Ul, args...) }

// If returns the provided node if the condition is true; otherwise, it returns nil.
func If(cond bool, node *h.Node) *h.Node {
	if cond {
		return node
	}
	return nil
}
