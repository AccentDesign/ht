package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	. "github.com/accentdesign/ht"
	"golang.org/x/net/html"
)

type Movie struct {
	ID       int
	Title    string
	Director string
	Year     int
}

type App struct {
	Prefix string
	movies []*Movie
}

func NewApp(prefix string) *App {
	return &App{
		Prefix: prefix,
		movies: []*Movie{
			{1, "The Shawshank Redemption", "Frank Darabont", 1994},
			{2, "The Godfather", "Francis Ford Coppola", 1972},
			{3, "The Dark Knight", "Christopher Nolan", 2008},
			{4, "The Godfather Part II", "Francis Ford Coppola", 1974},
			{5, "12 Angry Men", "Sidney Lumet", 1957},
			{6, "Schindler's List", "Steven Spielberg", 1993},
			{7, "The Lord of the Rings: The Return of the King", "Peter Jackson", 2003},
			{8, "Pulp Fiction", "Quentin Tarantino", 1994},
			{9, "The Lord of the Rings: The Fellowship of the Ring", "Peter Jackson", 2001},
			{10, "The Good, the Bad and the Ugly", "Sergio Leone", 1966},
		},
	}
}

func (a *App) prefixID() string {
	return strings.Trim(strings.ReplaceAll(a.Prefix, "/", "-"), "-")
}

func (a *App) resultsID() string {
	return "results-" + a.prefixID()
}

func (a *App) indicatorID() string {
	return "indicator-" + a.prefixID()
}

func (a *App) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET "+a.Prefix+"/results", a.handleResults)
}

func (a *App) Render() *html.Node {
	return Div(Class("card bg-base-100 shadow-xl w-full"),
		Div(Class("card-body"),
			Div(Class("flex flex-col gap-4"),
				H2(Class("card-title"), Text("Movie Database")),
				a.renderSearchBox(),
			),
			a.renderTable(a.movies),
		),
	)
}

func (a *App) renderSearchBox() *html.Node {
	return Label(
		Class("input w-full"),
		Input(
			Type("search"),
			Name("q"),
			Placeholder("Search movies by title or director..."),
			Class("grow"),
			HxGet(a.Prefix+"/results"),
			HxTrigger("input changed delay:500ms, search"),
			HxTarget("#"+a.resultsID()),
			HxIndicator("#"+a.indicatorID()),
		),
		// The indicator is hidden by default via the .htmx-indicator class
		Span(
			Id(a.indicatorID()),
			Class("htmx-indicator loading loading-spinner loading-sm text-primary"),
		),
	)
}

func (a *App) renderTable(movies []*Movie) *html.Node {
	return Div(Class("overflow-x-auto mt-4"),
		Table(Class("table table-zebra w-full"),
			Thead(
				Tr(
					Th(Text("Title")),
					Th(Text("Director")),
					Th(Text("Year")),
				),
			),
			Tbody(
				Id(a.resultsID()),
				Apply(Fragment(), a.renderRows(movies)),
			),
		),
	)
}

func (a *App) renderRows(movies []*Movie) []*html.Node {
	var rows []*html.Node
	for _, m := range movies {
		rows = append(rows, Tr(
			Td(Text(m.Title)),
			Td(Text(m.Director)),
			Td(Text(fmt.Sprintf("%d", m.Year))),
		))
	}
	if len(rows) == 0 {
		return []*html.Node{
			Tr(
				Td(Colspan("3"), Class("text-center text-base-content/50 py-4"), Text("No movies found matching your search.")),
			),
		}
	}
	return rows
}

func (a *App) handleResults(w http.ResponseWriter, r *http.Request) {
	// Add an artificial delay to show off the loading indicator nicely
	time.Sleep(300 * time.Millisecond)

	q := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))

	var results []*Movie
	for _, m := range a.movies {
		if q == "" || strings.Contains(strings.ToLower(m.Title), q) || strings.Contains(strings.ToLower(m.Director), q) {
			results = append(results, m)
		}
	}

	// Render just the rows to update the tbody
	fragment := Apply(Fragment(), a.renderRows(results))
	_ = html.Render(w, fragment)
}

func main() {
	movieApp := NewApp("/movies")

	mux := http.NewServeMux()
	movieApp.RegisterRoutes(mux)

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
					Title(Text("ht - Active Search Example")),
					Link(Href("https://cdn.jsdelivr.net/npm/daisyui@5"), Rel("stylesheet"), Type("text/css")),
					Script(Src("https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4")),
					Script(Src("https://unpkg.com/htmx.org@2.0.4")),
				),
				Body(
					Class("antialiased bg-base-300 min-h-screen text-base-content"),
					Div(
						Class("navbar bg-base-100 shadow-sm mb-8"),
						Div(Class("flex-1"),
							A(Class("btn btn-ghost text-xl"), Href("/"), Text("HTMX Active Search")),
						),
					),
					Div(
						Class("p-4 flex flex-col gap-8 max-w-4xl mx-auto items-start"),
						movieApp.Render(),
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

	log.Println("Serving Active Search example on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
