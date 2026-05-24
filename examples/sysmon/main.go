package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	. "github.com/accentdesign/ht"
	"github.com/accentdesign/ht/examples/core"
	"golang.org/x/net/html"
)

// ---------------------------------------------------------
// SSE Broker
// ---------------------------------------------------------

type Broker struct {
	clients        map[chan []byte]bool
	newClients     chan chan []byte
	defunctClients chan chan []byte
	messages       chan []byte
}

func NewBroker() *Broker {
	b := &Broker{
		clients:        make(map[chan []byte]bool),
		newClients:     make(chan chan []byte),
		defunctClients: make(chan chan []byte),
		messages:       make(chan []byte),
	}
	go b.start()
	return b
}

func (b *Broker) start() {
	for {
		select {
		case s := <-b.newClients:
			b.clients[s] = true
		case s := <-b.defunctClients:
			delete(b.clients, s)
			close(s)
		case msg := <-b.messages:
			// Broadcast to all connected clients
			for s := range b.clients {
				select {
				case s <- msg:
				default:
					// Drop message if client buffer is full to prevent blocking
				}
			}
		}
	}
}

func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Create a channel for this client
	messageChan := make(chan []byte, 10)
	b.newClients <- messageChan

	// Ensure cleanup when the client disconnects
	defer func() {
		b.defunctClients <- messageChan
	}()

	// Set required SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ctx := r.Context()
	for {
		select {
		case msg := <-messageChan:
			w.Write(msg)
			f.Flush()
		case <-ctx.Done():
			// Browser disconnected
			return
		}
	}
}

// ---------------------------------------------------------
// Metrics Engine
// ---------------------------------------------------------

func StartMetricsEngine(b *Broker) {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			// Simulate gathering system metrics
			cpu := rand.Intn(100)
			ram := rand.Intn(16000)

			// Render the updated UI component
			node := renderMetrics(cpu, ram)
			var buf bytes.Buffer
			_ = html.Render(&buf, node)

			// HTMX SSE requires data to be formatted cleanly.
			// Newlines inside the data must be prefixed with "data: " or stripped.
			// We strip them for simplicity.
			htmlStr := strings.ReplaceAll(buf.String(), "\n", " ")
			msg := fmt.Sprintf("event: metrics_update\ndata: %s\n\n", htmlStr)

			b.messages <- []byte(msg)
		}
	}()
}

// ---------------------------------------------------------
// UI Components
// ---------------------------------------------------------

func renderMetrics(cpu, ram int) *html.Node {
	// Simple color coding based on load
	cpuColor := "text-success"
	if cpu > 80 {
		cpuColor = "text-error"
	} else if cpu > 50 {
		cpuColor = "text-warning"
	}

	return Div(Class("stats bg-base-100 shadow w-full grid-cols-2"),
		Div(Class("stat"),
			Div(Class("stat-title"), Text("CPU Usage")),
			Div(Class("stat-value", cpuColor), Text(fmt.Sprintf("%d %%", cpu))),
			Div(Class("stat-desc"), Text("Current CPU usage")),
		),
		Div(Class("stat"),
			Div(Class("stat-title"), Text("RAM Usage")),
			Div(Class("stat-value text-info"), Text(fmt.Sprintf("%d MB", ram))),
			Div(Class("stat-desc"), Text("Out of 16000 MB")),
		),
	)
}

// ---------------------------------------------------------
// Main Server
// ---------------------------------------------------------

func main() {
	// 1. Initialize our Broker
	broker := NewBroker()

	// 2. Start our background worker that generates metrics
	StartMetricsEngine(broker)

	mux := http.NewServeMux()

	// 3. Register our SSE endpoint
	mux.Handle("/events", broker)

	// 4. Render the Shell Page
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		start := time.Now()

		page := Document(
			Doctype("html"),
			Html(
				Lang("en"),
				Head(
					Meta(Charset("utf-8")),
					Meta(Name("viewport"), Content("width=device-width", "initial-scale=1.0")),
					Title(Text("ht - System Monitor (SSE)")),
					Link(Href("https://cdn.jsdelivr.net/npm/daisyui@5"), Rel("stylesheet"), Type("text/css")),
					Link(Href("https://cdn.jsdelivr.net/npm/daisyui@5/themes.css"), Rel("stylesheet"), Type("text/css")),
					Script(Src("https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4")),
					Script(Src("https://unpkg.com/htmx.org@2.0.4")),
					Script(Src("https://unpkg.com/htmx-ext-sse@2.2.2/sse.js")),
					Script(Src("https://unpkg.com/alpinejs@3.x.x/dist/cdn.min.js"), Defer()),
					Script(Raw("document.documentElement.setAttribute('data-theme', localStorage.getItem('theme') || 'light')")),
				),
				Body(
					Class("antialiased bg-base-300 min-h-screen text-base-content"),
					Div(
						Class("navbar bg-base-100 shadow-sm mb-8"),
						Div(Class("flex-1"),
							A(Class("btn btn-ghost text-xl"), Href("/"), Text("HTMX Real-Time System Monitor")),
						),
						Div(Class("flex-none gap-2"),
							core.ThemeSwitcher(),
						),
					),
					Div(
						Class("p-4 flex flex-col gap-8 max-w-4xl mx-auto items-center"),
						Attr("hx-ext", "sse"),
						Attr("sse-connect", "/events"),
						Div(
							Class("w-full"),
							Attr("sse-swap", "metrics_update"),
							renderMetrics(0, 0),
						),
					),
				),
			),
		)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := html.Render(w, page); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("%s %s -> render=%v", r.Method, r.URL.Path, time.Since(start))
	})

	log.Println("Serving System Monitor on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
