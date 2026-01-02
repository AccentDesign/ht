package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	. "github.com/accentdesign/ht"
	"golang.org/x/net/html"
)

var (
	svgSun = `<svg class="swap-off h-8 w-8 fill-current" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
<path d="M5.64,17l-.71.71a1,1,0,0,0,0,1.41,1,1,0,0,0,1.41,0l.71-.71A1,1,0,0,0,5.64,17ZM5,12a1,1,0,0,0-1-1H3a1,1,0,0,0,0,2H4A1,1,0,0,0,5,12Zm7-7a1,1,0,0,0,1-1V3a1,1,0,0,0-2,0V4A1,1,0,0,0,12,5ZM5.64,7.05a1,1,0,0,0,.7.29,1,1,0,0,0,.71-.29,1,1,0,0,0,0-1.41l-.71-.71A1,1,0,0,0,4.93,6.34Zm12,.29a1,1,0,0,0,.7-.29l.71-.71a1,1,0,1,0-1.41-1.41L17,5.64a1,1,0,0,0,0,1.41A1,1,0,0,0,17.66,7.34ZM21,11H20a1,1,0,0,0,0,2h1a1,1,0,0,0,0-2Zm-9,8a1,1,0,0,0-1,1v1a1,1,0,0,0,2,0V20A1,1,0,0,0,12,19ZM18.36,17A1,1,0,0,0,17,18.36l.71.71a1,1,0,0,0,1.41,0,1,1,0,0,0,0-1.41ZM12,6.5A5.5,5.5,0,1,0,17.5,12,5.51,5.51,0,0,0,12,6.5Zm0,9A3.5,3.5,0,1,1,15.5,12,3.5,3.5,0,0,1,12,15.5Z" />
</svg>`
	svgMoon = `<svg class="swap-on h-8 w-8 fill-current" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
<path d="M21.64,13a1,1,0,0,0-1.05-.14,8.05,8.05,0,0,1-3.37.73A8.15,8.15,0,0,1,9.08,5.49a8.59,8.59,0,0,1,.25-2A1,1,0,0,0,8,2.36,10.14,10.14,0,1,0,22,14.05,1,1,0,0,0,21.64,13Zm-9.5,6.69A8.14,8.14,0,0,1,7.08,5.22v.27A10.15,10.15,0,0,0,17.22,15.63a9.79,9.79,0,0,0,2.1-.22A8.11,8.11,0,0,1,12.14,19.73Z" />
</svg>`
)

// This example serves a simple page that showcases a few common daisyUI components
// using the ht helpers. Run: go run ./examples/daisyui and open http://localhost:8080
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
					Title(Text("daisyUI")),
					Link(Href("https://cdn.jsdelivr.net/npm/daisyui@5"), Rel("stylesheet"), Type("text/css")),
					Script(Src("https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4")),
				),
				Body(Class("antialiased"),
					// Navbar
					Div(Class("navbar", "shadow-sm"),
						Div(Class("flex-1"),
							A(Class("btn", "btn-ghost", "text-xl"), Href("#"), Text("daisyUI")),
						),
						Div(Class("flex-none"),
							Label(Class("swap", "swap-rotate"),
								Input(Class("theme-controller"), Type("checkbox"), Value("dark")),
								Raw(svgMoon),
								Raw(svgSun),
							),
						),
					),

					// Content container
					Div(Class("px-6", "py-8", "flex", "flex-col", "gap-16"),
						// Buttons
						SectionBlock(
							"Buttons",
							Div(Class("flex", "flex-wrap", "gap-2"),
								Button(Class("btn"), Text("Default")),
								Button(Class("btn", "btn-primary"), Text("Primary")),
								Button(Class("btn", "btn-secondary"), Text("Secondary")),
								Button(Class("btn", "btn-accent"), Text("Accent")),
								Button(Class("btn", "btn-outline"), Text("Outline")),
								Button(Class("btn", "btn-soft"), Text("Soft")),
								Button(Class("btn", "btn-ghost"), Text("Ghost")),
								Button(Class("btn", "btn-link"), Text("Link")),
								Button(Class("btn", "btn-disabled"), Disabled(), Text("Disabled")),
							),
						),

						// Alerts
						SectionBlock(
							"Alerts",
							Div(Class("flex", "flex-col", "gap-2"),
								Div(Class("alert"),
									Span(Text("Neutral alert — basic usage.")),
								),
								Div(Class("alert", "alert-info"),
									Span(Text("Info alert — some extra context.")),
								),
								Div(Class("alert", "alert-success"),
									Span(Text("Success alert — everything worked!")),
								),
								Div(Class("alert", "alert-warning"),
									Span(Text("Warning alert — something looks off.")),
								),
								Div(Class("alert", "alert-error"),
									Span(Text("Error alert — action failed.")),
								),
							),
						),

						// Card
						SectionBlock(
							"Card",
							Div(Class("card", "bg-base-200", "shadow-xl", "w-full", "sm:w-96"),
								Figure(
									Img(
										Src("https://img.daisyui.com/images/stock/photo-1606107557195-0e29a4b5b4aa.webp"),
										Alt("Shoes"),
									),
								),
								Div(Class("card-body"),
									H2(Class("card-title"), Text("Shoes!")),
									P(Text("If a dog chews shoes whose shoes does he choose?")),
									Div(Class("card-actions justify-end"),
										Button(Class("btn", "btn-primary"), Text("Buy Now")),
									),
								),
							),
						),

						// Form controls
						SectionBlock(
							"Form",
							Div(Class("grid", "grid-cols-1", "sm:grid-cols-2", "gap-4"),
								Fieldset(Class("fieldset"),
									Label(Class("label"), Text("Username")),
									Input(Class("input", "w-full"), Type("text"), Placeholder("johndoe")),
									Label(Class("label"), Text("Required")),
								),
								Fieldset(Class("fieldset"),
									Label(Class("label"), Text("Password")),
									Input(Class("input", "w-full"), Type("password"), Placeholder("••••••••")),
									Label(Class("label"),
										A(Class("link", "link-hover"), Href("#"), Text("Forgot password?")),
									),
								),
								Fieldset(Class("fieldset"),
									Label(Class("label", "justify-start", "gap-4"),
										Input(Class("checkbox", "checkbox-primary"), Type("checkbox"), Checked()),
										Text("Accept terms and conditions"),
									),
								),
								Div(Class("sm:col-span-2"),
									Button(Class("btn", "btn-primary", "w-full", "sm:w-auto"), Text("Submit")),
								),
							),
						),

						// Badges
						SectionBlock(
							"Badges",
							Div(Class("flex", "flex-wrap", "gap-2"),
								Div(Class("badge"), Text("neutral")),
								Div(Class("badge", "badge-primary"), Text("primary")),
								Div(Class("badge", "badge-secondary"), Text("secondary")),
								Div(Class("badge", "badge-accent"), Text("accent")),
								Div(Class("badge", "badge-ghost"), Text("ghost")),
								Div(Class("badge", "badge-outline"), Text("outline")),
							),
						),

						// Stats
						SectionBlock(
							"Stats",
							Div(Class("stats", "stats-vertical", "lg:stats-horizontal", "shadow"),
								Div(Class("stat"),
									Div(Class("stat-title"), Text("Downloads")),
									Div(Class("stat-value"), Text("31K")),
									Div(Class("stat-desc"), Text("Jan 1st - Feb 1st")),
								),
								Div(Class("stat"),
									Div(Class("stat-title"), Text("New Users")),
									Div(Class("stat-value"), Text("4,200")),
									Div(Class("stat-desc"), Text("↗︎ 400 (22%)")),
								),
								Div(Class("stat"),
									Div(Class("stat-title"), Text("New Registers")),
									Div(Class("stat-value"), Text("1,200")),
									Div(Class("stat-desc"), Text("↘︎ 90 (14%)")),
								),
							),
						),
					),
				),
			),
		)

		if err := html.Render(w, page); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dur := time.Since(start)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		ms := float64(dur.Microseconds()) / 1000.0
		w.Header().Set("Server-Timing", fmt.Sprintf("render;desc=html.Render;dur=%.2f", ms))
		log.Printf("%s %s -> render=%v", r.Method, r.URL.Path, dur)
	})

	log.Println("Serving daisyUI example on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// SectionBlock is a tiny helper to render a titled section consistently in this example file.
func SectionBlock(title string, content *html.Node) *html.Node {
	return Div(Class("flex flex-col gap-3"),
		H2(Class("text-xl font-semibold"), Text(title)),
		content,
	)
}
