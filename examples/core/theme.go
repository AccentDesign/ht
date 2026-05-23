package core

import (
	. "github.com/accentdesign/ht"
	h "golang.org/x/net/html"
)

// ThemeSwitcher returns a daisyUI-based theme switcher dropdown component.
// It displays a dropdown with all available themes and handles theme synchronization.
func ThemeSwitcher() *h.Node {
	xData := `{
		theme: localStorage.getItem('theme') || 'light',
		themes: [
			{value: 'light', label: '☀️ Light'},
			{value: 'dark', label: '🌙 Dark'},
			{value: 'cupcake', label: '🧁 Cupcake'},
			{value: 'bumblebee', label: '🐝 Bumblebee'},
			{value: 'emerald', label: '❇️ Emerald'},
			{value: 'corporate', label: '💼 Corporate'},
			{value: 'synthwave', label: '🌌 Synthwave'},
			{value: 'retro', label: '📻 Retro'},
			{value: 'cyberpunk', label: '🤖 Cyberpunk'},
			{value: 'valentine', label: '💕 Valentine'},
			{value: 'halloween', label: '🎃 Halloween'},
			{value: 'garden', label: '🌷 Garden'},
			{value: 'forest', label: '🌲 Forest'},
			{value: 'aqua', label: '🌊 Aqua'},
			{value: 'lofi', label: '🎧 Lofi'},
			{value: 'pastel', label: '🌸 Pastel'},
			{value: 'fantasy', label: '🧙 Fantasy'},
			{value: 'wireframe', label: '📐 Wireframe'},
			{value: 'black', label: '🖤 Black'},
			{value: 'luxury', label: '💎 Luxury'},
			{value: 'dracula', label: '🧛 Dracula'},
			{value: 'cmyk', label: '🎨 CMYK'},
			{value: 'autumn', label: '🍁 Autumn'},
			{value: 'business', label: '📈 Business'},
			{value: 'acid', label: '🧪 Acid'},
			{value: 'lemonade', label: '🍋 Lemonade'},
			{value: 'night', label: '🌃 Night'},
			{value: 'coffee', label: '☕ Coffee'},
			{value: 'winter', label: '❄️ Winter'},
			{value: 'dim', label: '🔅 Dim'},
			{value: 'nord', label: '🏔️ Nord'},
			{value: 'sunset', label: '🌇 Sunset'},
			{value: 'caramellatte', label: '🍮 Caramellatte'},
			{value: 'abyss', label: '🌑 Abyss'},
			{value: 'silk', label: '🪷 Silk'}
		]
	}`

	return Div(
		Class("dropdown dropdown-end"),
		X("data", xData),
		X("init", "document.documentElement.setAttribute('data-theme', theme); $watch('theme', val => { localStorage.setItem('theme', val); document.documentElement.setAttribute('data-theme', val) })"),
		Div(
			Class("btn btn-ghost rounded-field"),
			Tabindex("0"),
			Role("button"),
			Text("Theme"),
		),
		Ul(
			Class("menu dropdown-content bg-base-200 rounded-box z-1 mt-4 w-56 p-2 shadow-sm h-[30.5rem] flex-nowrap overflow-y-auto *:w-full"),
			Template(
				X("for", "t in themes"),
				Li(
					Input(
						Type("radio"),
						Name("theme-dropdown"),
						Class("theme-controller btn btn-sm btn-ghost justify-start"),
						XBind("aria-label", "t.label"),
						XBind("value", "t.value"),
						X("model", "theme"),
					),
				),
			),
		),
	)
}
