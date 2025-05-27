package ht

import (
	"context"
	h "golang.org/x/net/html"
	"io"
)

// TemplComponent is a component that renders a h.Node.
// It implements the templ.Component interface, allowing it to be used
// in templ components.
type TemplComponent struct {
	Node *h.Node
}

func (c *TemplComponent) Render(ctx context.Context, w io.Writer) error {
	return h.Render(w, c.Node)
}

// Templ creates a new TemplComponent from a given Node.
func Templ(node *h.Node) *TemplComponent {
	return &TemplComponent{Node: node}
}
