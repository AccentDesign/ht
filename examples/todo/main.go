package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	. "github.com/accentdesign/ht"
	"golang.org/x/net/html"
)

// Item represents a single todo entry.
type Item struct {
	ID   int
	Text string
	Done bool
}

// App encapsulates the state and routing for a Todo component.
type App struct {
	// Prefix is the route prefix (e.g., "/work-todos").
	Prefix string

	mu     sync.RWMutex
	items  []*Item
	nextID int
}

// NewApp returns a new Todo App component.
func NewApp(prefix string) *App {
	return &App{
		Prefix: prefix,
		items:  make([]*Item, 0),
		nextID: 1,
	}
}

// prefixID is a helper to generate unique HTML IDs per instance.
func (a *App) prefixID() string {
	return strings.Trim(strings.ReplaceAll(a.Prefix, "/", "-"), "-")
}

// counterID returns the unique DOM ID for this component's counter.
func (a *App) counterID() string {
	return "todo-counter-" + a.prefixID()
}

// formID returns the unique DOM ID for this component's form.
func (a *App) formID() string {
	return "todo-form-" + a.prefixID()
}

// listID returns the unique DOM ID for this component's list container.
func (a *App) listID() string {
	return "todo-list-" + a.prefixID()
}

// RegisterRoutes registers the HTMX endpoints to the provided multiplexer.
func (a *App) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST "+a.Prefix+"/add", a.handleAdd)
	mux.HandleFunc("PUT "+a.Prefix+"/{id}/toggle", a.handleToggle)
	mux.HandleFunc("DELETE "+a.Prefix+"/{id}", a.handleDelete)
}

// Render returns the full HTML node for the initial render (List + Form).
func (a *App) Render() *html.Node {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return Div(Class("card bg-base-100 shadow-xl w-full"),
		Div(Class("card-body"),
			Div(Class("flex gap-4 items-center"),
				H2(Class("card-title capitalize"), Text(strings.Trim(a.Prefix, "/")+" Tasks")),
				a.renderCounter(),
			),
			a.renderList(),
			a.renderForm(""), // Initial empty form
		),
	)
}

// renderCounter returns the HTML node containing the task counter.
func (a *App) renderCounter() *html.Node {
	count := len(a.items)
	label := "task"
	if count != 1 {
		label = "tasks"
	}

	return Span(
		Id(a.counterID()),
		Class("badge badge-soft"),
		Text(strconv.Itoa(count)+" "+label),
	)
}

// renderList returns the HTML node containing all current items.
func (a *App) renderList() *html.Node {
	var rows []*html.Node
	for _, item := range a.items {
		rows = append(rows, a.renderRow(item))
	}

	return Div(
		Id(a.listID()),
		Class("flex flex-col gap-2 mt-4"),
		Apply(Fragment(), rows), // Merge all rows as children
	)
}

// renderRow returns the HTML node for a single item (Partial).
func (a *App) renderRow(item *Item) *html.Node {
	rowID := fmt.Sprintf("todo-row-%s-%d", a.prefixID(), item.ID)

	return Div(
		Id(rowID),
		Class("todo-row flex items-center gap-3 p-3 bg-base-200 rounded-box shadow-sm transition-all"),
		Input(
			Type("checkbox"),
			Class("checkbox checkbox-primary"),
			If(item.Done, Checked()),
			HxPut(fmt.Sprintf("%s/%d/toggle", a.Prefix, item.ID)),
			HxTarget("closest div.todo-row"),
			HxSwap("outerHTML"),
		),
		Span(
			Class("flex-1 transition-all"),
			If(item.Done, Class("line-through text-base-content/50")),
			Text(item.Text),
		),
		Button(
			Class("btn btn-ghost btn-sm text-error"),
			HxDelete(fmt.Sprintf("%s/%d", a.Prefix, item.ID)),
			HxTarget("closest div.todo-row"),
			HxSwap("outerHTML"),
			Text("Delete"),
		),
	)
}

// renderForm returns the add item form, optionally displaying an error message (Partial).
func (a *App) renderForm(errMsg string) *html.Node {
	return Form(
		Id(a.formID()),
		Class("mt-4 flex flex-col gap-2"),
		X("data", "{ loading: false }"),
		XOn("submit", "loading = true"),
		HxOn("htmx:after-request", "loading = false"),
		HxPost(a.Prefix+"/add"),
		HxTarget("#"+a.listID()),
		HxSwap("beforeend"),

		Div(Class("flex gap-2"),
			Input(
				Type("text"),
				Name("text"),
				Placeholder("What needs to be done?"),
				Class("input flex-1"),
				If(errMsg != "", Class("input-error")),
				Required(),
				Autofocus(),
			),
			Button(
				Type("submit"),
				Class("btn btn-primary"),
				XBind("disabled", "loading"),
				Text("Add"),
			),
		),
		If(errMsg != "",
			Span(Class("text-sm text-error"), Text(errMsg)),
		),
	)
}

// ---------------------------------------------------------
// HTTP Handlers
// ---------------------------------------------------------

func (a *App) handleAdd(w http.ResponseWriter, r *http.Request) {
	text := strings.TrimSpace(r.FormValue("text"))
	if text == "" {
		// Render form with error, swap out the existing form Out Of Band
		form := Apply(a.renderForm("Task cannot be empty"), HxSwapOob("true"))
		_ = html.Render(w, form)
		return
	}

	a.mu.Lock()
	newItem := &Item{
		ID:   a.nextID,
		Text: text,
		Done: false,
	}
	a.nextID++
	a.items = append(a.items, newItem)
	a.mu.Unlock()

	// Return the newly created row (appended to list)
	_ = html.Render(w, a.renderRow(newItem))
	// AND return a cleared form to replace the old form (Out Of Band swap)
	_ = html.Render(w, Apply(a.renderForm(""), HxSwapOob("true")))
	// AND return the updated counter (Out Of Band swap)
	_ = html.Render(w, Apply(a.renderCounter(), HxSwapOob("true")))
}

func (a *App) handleToggle(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	a.mu.Lock()
	var updatedItem *Item
	for _, item := range a.items {
		if item.ID == id {
			item.Done = !item.Done
			updatedItem = item
			break
		}
	}
	a.mu.Unlock()

	if updatedItem == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	// Re-render just the updated row
	_ = html.Render(w, a.renderRow(updatedItem))
}

func (a *App) handleDelete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	a.mu.Lock()
	for i, item := range a.items {
		if item.ID == id {
			// Remove from slice
			a.items = append(a.items[:i], a.items[i+1:]...)
			break
		}
	}
	a.mu.Unlock()

	// Return the updated counter (Out Of Band swap)
	_ = html.Render(w, Apply(a.renderCounter(), HxSwapOob("true")))
}

// ---------------------------------------------------------
// Main Server Example
// ---------------------------------------------------------

func main() {
	// Initialize two independent instances
	workTodos := NewApp("/work")
	homeTodos := NewApp("/home")

	// Pre-fill some tasks
	workTodos.items = append(workTodos.items, &Item{ID: workTodos.nextID, Text: "Review PR #42", Done: false})
	workTodos.nextID++
	homeTodos.items = append(homeTodos.items, &Item{ID: homeTodos.nextID, Text: "Buy groceries", Done: true})
	homeTodos.nextID++

	mux := http.NewServeMux()

	// Register their API endpoints
	workTodos.RegisterRoutes(mux)
	homeTodos.RegisterRoutes(mux)

	// Render the main page
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
					Title(Text("ht - Todo Example")),
					Link(Href("https://cdn.jsdelivr.net/npm/daisyui@5"), Rel("stylesheet"), Type("text/css")),
					Script(Src("https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4")),
					Script(Src("https://unpkg.com/htmx.org@2.0.4")),
					Script(Src("https://unpkg.com/alpinejs@3.x.x/dist/cdn.min.js"), Defer()),
				),
				Body(
					Class("antialiased bg-base-300 min-h-screen text-base-content"),
					Div(
						Class("navbar bg-base-100 shadow-sm mb-8"),
						Div(Class("flex-1"),
							A(Class("btn btn-ghost text-xl"), Href("/"), Text("HTMX + daisyUI Todo")),
						),
					),
					Div(
						Class("p-4 flex flex-col lg:flex-row gap-8 max-w-7xl mx-auto items-start"),
						workTodos.Render(),
						homeTodos.Render(),
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

	log.Println("Serving Todo example on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
