package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"sync"
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

	messageChan := make(chan []byte, 10)
	b.newClients <- messageChan

	defer func() {
		b.defunctClients <- messageChan
	}()

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
			return
		}
	}
}

// ---------------------------------------------------------
// Metrics Engine
// ---------------------------------------------------------

func StartMetricsEngine(ctx context.Context, b *Broker) {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				var m runtime.MemStats
				runtime.ReadMemStats(&m)

				goroutines := runtime.NumGoroutine()
				memMB := int(m.Alloc / 1024 / 1024)
				heapObjects := int(m.HeapObjects)
				numGC := int(m.NumGC)

				node := renderMetrics(goroutines, memMB, heapObjects, numGC)

				var buf bytes.Buffer
				_ = html.Render(&buf, node)

				// Build SSE payload with proper multi-line data framing.
				// Each line of the HTML is prefixed with "data: " so the SSE
				// protocol reassembles them correctly on the client.
				var sb strings.Builder
				sb.WriteString("event: metrics_update\n")
				for _, line := range strings.Split(buf.String(), "\n") {
					sb.WriteString("data: ")
					sb.WriteString(line)
					sb.WriteString("\n")
				}
				sb.WriteString("\n")

				b.messages <- []byte(sb.String())

			case <-ctx.Done():
				return
			}
		}
	}()
}

// ---------------------------------------------------------
// UI Components
// ---------------------------------------------------------

func renderMetrics(goroutines, memMB, heapObjects, numGC int) *html.Node {
	return Div(Class("stats stats-vertical lg:stats-horizontal bg-base-100 shadow w-full"),
		Div(Class("stat"),
			Div(Class("stat-title"), Text("Goroutines")),
			Div(Class("stat-value"), Text(fmt.Sprintf("%d", goroutines))),
			Div(Class("stat-desc"), Text("Active Goroutines")),
		),
		Div(Class("stat"),
			Div(Class("stat-title"), Text("Memory Allocated")),
			Div(Class("stat-value text-info"), Text(fmt.Sprintf("%d MB", memMB))),
			Div(Class("stat-desc"), Text("Heap Allocation")),
		),
		Div(Class("stat"),
			Div(Class("stat-title"), Text("Heap Objects")),
			Div(Class("stat-value text-accent"), Text(fmt.Sprintf("%d", heapObjects))),
			Div(Class("stat-desc"), Text("Active objects")),
		),
		Div(Class("stat"),
			Div(Class("stat-title"), Text("GC Cycles")),
			Div(Class("stat-value text-secondary"), Text(fmt.Sprintf("%d", numGC))),
			Div(Class("stat-desc"), Text("Completed collections")),
		),
	)
}

// ---------------------------------------------------------
// Main Server
// ---------------------------------------------------------

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	broker := NewBroker()
	StartMetricsEngine(ctx, broker)

	mux := http.NewServeMux()

	var memorySink [][]byte
	var memoryMutex sync.Mutex

	mux.HandleFunc("POST /allocate", func(w http.ResponseWriter, r *http.Request) {
		memoryMutex.Lock()
		defer memoryMutex.Unlock()
		memorySink = append(memorySink, make([]byte, 5*1024*1024))
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("POST /gc", func(w http.ResponseWriter, r *http.Request) {
		runtime.GC()
		w.WriteHeader(http.StatusOK)
	})

	// Background worker to drain memory over time.
	go func() {
		for {
			select {
			case <-time.After(2 * time.Second):
				memoryMutex.Lock()
				if len(memorySink) > 0 {
					// Reslice to drop the oldest chunk; the GC will reclaim
					// the backing memory on the next collection cycle.
					memorySink[0] = nil
					memorySink = memorySink[1:]
					runtime.GC()
				}
				memoryMutex.Unlock()
			case <-ctx.Done():
				return
			}
		}
	}()

	mux.Handle("/events", broker)

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
							Class("text-center max-w-2xl mb-2"),
							H2(Class("text-2xl font-bold mb-4"), Text("Go Memory & Garbage Collection Demo")),
							P(Class("mb-4"), Text("Watch how Go handles memory! The Heap Objects count will slowly increase as the server continuously streams these live UI updates to you (allocating tiny amounts of memory). Go's Garbage Collector (GC) runs lazily to save CPU, so it waits until memory pressure builds up before cleaning them away.")),
							P(Class("text-sm opacity-80"), Text("Click 'Allocate 5MB Memory' to manually build pressure. A background worker will slowly drain this memory over time, triggering natural GC cycles. Or, click 'Force Garbage Collection' to trigger it manually right now!")),
						),
						Div(
							Class("w-full"),
							Attr("sse-swap", "metrics_update"),
							renderMetrics(0, 0, 0, 0),
						),
						Div(Class("flex gap-4 mt-4"),
							Button(
								Class("btn btn-primary"),
								Attr("hx-post", "/allocate"),
								Attr("hx-swap", "none"),
								Text("Allocate 5MB Memory"),
							),
							Button(
								Class("btn btn-secondary"),
								Attr("hx-post", "/gc"),
								Attr("hx-swap", "none"),
								Text("Force Garbage Collection"),
							),
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
