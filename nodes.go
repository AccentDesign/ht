package ht

import (
	"fmt"
	"strings"

	h "golang.org/x/net/html"
	a "golang.org/x/net/html/atom"
)

var mergeAttrMap = map[string]string{
	"class":   " ",
	"content": ", ",
}

func mergeAttr(oldVal, newVal, joiner string) string {
	if oldVal == "" {
		return strings.TrimSpace(newVal)
	}
	if newVal == "" {
		return oldVal
	}
	seen := make(map[string]struct{}, 8)
	out := make([]string, 0, 8)
	for _, s := range strings.Fields(oldVal) {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			out = append(out, s)
		}
	}
	for _, s := range strings.Fields(newVal) {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			out = append(out, s)
		}
	}
	return strings.Join(out, joiner)
}

func Comment(data string) *h.Node {
	return &h.Node{Type: h.CommentNode, Data: data}
}

func Doctype(data string) *h.Node {
	return &h.Node{Type: h.DoctypeNode, Data: data}
}

func Document(children ...*h.Node) *h.Node {
	node := &h.Node{Type: h.DocumentNode}
	for _, child := range children {
		node.AppendChild(child)
	}
	return node
}

func Element(tag a.Atom, args ...interface{}) *h.Node {
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
					if joiner, ok := mergeAttrMap[attr.Key]; ok {
						node.Attr[i].Val = mergeAttr(attr.Val, v.Val, joiner)
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
		default:
			fmt.Printf("unknown argument type in Node: %T\n", v)
		}
	}
	return node
}

func Raw(data string) *h.Node {
	return &h.Node{Type: h.RawNode, Data: data}
}

func Text(data string) *h.Node {
	return &h.Node{Type: h.TextNode, Data: data}
}

func A(args ...interface{}) *h.Node          { return Element(a.A, args...) }
func Abbr(args ...interface{}) *h.Node       { return Element(a.Abbr, args...) }
func Address(args ...interface{}) *h.Node    { return Element(a.Address, args...) }
func Article(args ...interface{}) *h.Node    { return Element(a.Article, args...) }
func Aside(args ...interface{}) *h.Node      { return Element(a.Aside, args...) }
func B(args ...interface{}) *h.Node          { return Element(a.B, args...) }
func Blockquote(args ...interface{}) *h.Node { return Element(a.Blockquote, args...) }
func Body(args ...interface{}) *h.Node       { return Element(a.Body, args...) }
func Button(args ...interface{}) *h.Node     { return Element(a.Button, args...) }
func Br(args ...interface{}) *h.Node         { return Element(a.Br, args...) }
func Caption(args ...interface{}) *h.Node    { return Element(a.Caption, args...) }
func Cite(args ...interface{}) *h.Node       { return Element(a.Cite, args...) }
func Code(args ...interface{}) *h.Node       { return Element(a.Code, args...) }
func Col(args ...interface{}) *h.Node        { return Element(a.Col, args...) }
func Colgroup(args ...interface{}) *h.Node   { return Element(a.Colgroup, args...) }
func Dd(args ...interface{}) *h.Node         { return Element(a.Dd, args...) }
func Details(args ...interface{}) *h.Node    { return Element(a.Details, args...) }
func Dialog(args ...interface{}) *h.Node     { return Element(a.Dialog, args...) }
func Div(args ...interface{}) *h.Node        { return Element(a.Div, args...) }
func Dl(args ...interface{}) *h.Node         { return Element(a.Dl, args...) }
func Dt(args ...interface{}) *h.Node         { return Element(a.Dt, args...) }
func Em(args ...interface{}) *h.Node         { return Element(a.Em, args...) }
func Fieldset(args ...interface{}) *h.Node   { return Element(a.Fieldset, args...) }
func Figcaption(args ...interface{}) *h.Node { return Element(a.Figcaption, args...) }
func Figure(args ...interface{}) *h.Node     { return Element(a.Figure, args...) }
func Footer(args ...interface{}) *h.Node     { return Element(a.Footer, args...) }
func Form(args ...interface{}) *h.Node       { return Element(a.Form, args...) }
func H1(args ...interface{}) *h.Node         { return Element(a.H1, args...) }
func H2(args ...interface{}) *h.Node         { return Element(a.H2, args...) }
func H3(args ...interface{}) *h.Node         { return Element(a.H3, args...) }
func H4(args ...interface{}) *h.Node         { return Element(a.H4, args...) }
func H5(args ...interface{}) *h.Node         { return Element(a.H5, args...) }
func Head(args ...interface{}) *h.Node       { return Element(a.Head, args...) }
func Header(args ...interface{}) *h.Node     { return Element(a.Header, args...) }
func Hr(args ...interface{}) *h.Node         { return Element(a.Hr, args...) }
func Html(args ...interface{}) *h.Node       { return Element(a.Html, args...) }
func I(args ...interface{}) *h.Node          { return Element(a.I, args...) }
func Img(args ...interface{}) *h.Node        { return Element(a.Img, args...) }
func Input(args ...interface{}) *h.Node      { return Element(a.Input, args...) }
func Label(args ...interface{}) *h.Node      { return Element(a.Label, args...) }
func Legend(args ...interface{}) *h.Node     { return Element(a.Legend, args...) }
func Li(args ...interface{}) *h.Node         { return Element(a.Li, args...) }
func Link(args ...interface{}) *h.Node       { return Element(a.Link, args...) }
func Main(args ...interface{}) *h.Node       { return Element(a.Main, args...) }
func Mark(args ...interface{}) *h.Node       { return Element(a.Mark, args...) }
func Meta(args ...interface{}) *h.Node       { return Element(a.Meta, args...) }
func Nav(args ...interface{}) *h.Node        { return Element(a.Nav, args...) }
func Ol(args ...interface{}) *h.Node         { return Element(a.Ol, args...) }
func Optgroup(args ...interface{}) *h.Node   { return Element(a.Optgroup, args...) }
func Option(args ...interface{}) *h.Node     { return Element(a.Option, args...) }
func P(args ...interface{}) *h.Node          { return Element(a.P, args...) }
func Pre(args ...interface{}) *h.Node        { return Element(a.Pre, args...) }
func Script(args ...interface{}) *h.Node     { return Element(a.Script, args...) }
func Section(args ...interface{}) *h.Node    { return Element(a.Section, args...) }
func Select(args ...interface{}) *h.Node     { return Element(a.Select, args...) }
func Small(args ...interface{}) *h.Node      { return Element(a.Small, args...) }
func Span(args ...interface{}) *h.Node       { return Element(a.Span, args...) }
func Strong(args ...interface{}) *h.Node     { return Element(a.Strong, args...) }
func Style(args ...interface{}) *h.Node      { return Element(a.Style, args...) }
func Sub(args ...interface{}) *h.Node        { return Element(a.Sub, args...) }
func Summary(args ...interface{}) *h.Node    { return Element(a.Summary, args...) }
func Sup(args ...interface{}) *h.Node        { return Element(a.Sup, args...) }
func Table(args ...interface{}) *h.Node      { return Element(a.Table, args...) }
func Tbody(args ...interface{}) *h.Node      { return Element(a.Tbody, args...) }
func Td(args ...interface{}) *h.Node         { return Element(a.Td, args...) }
func Textarea(args ...interface{}) *h.Node   { return Element(a.Textarea, args...) }
func Tfoot(args ...interface{}) *h.Node      { return Element(a.Tfoot, args...) }
func Th(args ...interface{}) *h.Node         { return Element(a.Th, args...) }
func Thead(args ...interface{}) *h.Node      { return Element(a.Thead, args...) }
func Title(args ...interface{}) *h.Node      { return Element(a.Title, args...) }
func Tr(args ...interface{}) *h.Node         { return Element(a.Tr, args...) }
func Ul(args ...interface{}) *h.Node         { return Element(a.Ul, args...) }

// conditional helpers

func If(cond bool, node *h.Node) *h.Node {
	if cond {
		return node
	}
	return nil
}
