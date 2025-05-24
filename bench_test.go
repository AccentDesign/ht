package ht

import (
	"bytes"
	"golang.org/x/net/html"
	"testing"
)

func BenchmarkDiv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Div(Text("Hello, World!"))
	}
}

func BenchmarkText(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Text("Benchmarking text node")
	}
}

func BenchmarkRaw(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Raw("<h1>Raw HTML</h1>")
	}
}

func BenchmarkComplexNode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Div(
			Div(Text("Nested")),
			Text("Sibling"),
			Raw("<span>Raw</span>"),
		)
	}
}

func BenchmarkDivRender(b *testing.B) {
	for i := 0; i < b.N; i++ {
		node := Div(Text("Hello, World!"))
		var buf bytes.Buffer
		_ = html.Render(&buf, node)
	}
}

func BenchmarkComplexNodeRender(b *testing.B) {
	for i := 0; i < b.N; i++ {
		node := Div(
			Div(Text("Nested")),
			Text("Sibling"),
			Raw("<span>Raw</span>"),
		)
		var buf bytes.Buffer
		_ = html.Render(&buf, node)
	}
}
