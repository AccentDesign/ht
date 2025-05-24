package ht

import (
	"golang.org/x/net/html"
)

func Attr(k, v string) html.Attribute { return html.Attribute{Key: k, Val: v} }

// Common attributes

func Action(v string) html.Attribute       { return Attr("action", v) }
func Alt(v string) html.Attribute          { return Attr("alt", v) }
func Autocomplete(v string) html.Attribute { return Attr("autocomplete", v) }
func Charset(v string) html.Attribute      { return Attr("charset", v) }
func Class(v string) html.Attribute        { return Attr("class", v) }
func Content(v string) html.Attribute      { return Attr("content", v) }
func Download(v string) html.Attribute     { return Attr("download", v) }
func Enctype(v string) html.Attribute      { return Attr("enctype", v) }
func For(v string) html.Attribute          { return Attr("for", v) }
func Height(v string) html.Attribute       { return Attr("height", v) }
func Href(v string) html.Attribute         { return Attr("href", v) }
func Id(v string) html.Attribute           { return Attr("id", v) }
func LabelAttr(v string) html.Attribute    { return Attr("label", v) }
func Lang(v string) html.Attribute         { return Attr("lang", v) }
func Max(v string) html.Attribute          { return Attr("max", v) }
func Method(v string) html.Attribute       { return Attr("method", v) }
func Min(v string) html.Attribute          { return Attr("min", v) }
func Name(v string) html.Attribute         { return Attr("name", v) }
func Pattern(v string) html.Attribute      { return Attr("pattern", v) }
func Placeholder(v string) html.Attribute  { return Attr("placeholder", v) }
func Rel(v string) html.Attribute          { return Attr("rel", v) }
func Size(v string) html.Attribute         { return Attr("size", v) }
func Src(v string) html.Attribute          { return Attr("src", v) }
func Step(v string) html.Attribute         { return Attr("step", v) }
func StyleAttr(v string) html.Attribute    { return Attr("style", v) }
func Target(v string) html.Attribute       { return Attr("target", v) }
func TitleAttr(v string) html.Attribute    { return Attr("title", v) }
func Type(v string) html.Attribute         { return Attr("type", v) }
func Value(v string) html.Attribute        { return Attr("value", v) }
func Width(v string) html.Attribute        { return Attr("width", v) }

// Boolean attributes

func Autofocus() html.Attribute { return Attr("autofocus", "") }
func Checked() html.Attribute   { return Attr("checked", "") }
func Disabled() html.Attribute  { return Attr("disabled", "") }
func Hidden() html.Attribute    { return Attr("hidden", "") }
func Multiple() html.Attribute  { return Attr("multiple", "") }
func Readonly() html.Attribute  { return Attr("readonly", "") }
func Required() html.Attribute  { return Attr("required", "") }
func Selected() html.Attribute  { return Attr("selected", "") }

// htmx attributes

func HxBoost(v string) html.Attribute      { return Attr("hx-boost", v) }
func HxConfirm(v string) html.Attribute    { return Attr("hx-confirm", v) }
func HxDelete(v string) html.Attribute     { return Attr("hx-delete", v) }
func HxGet(v string) html.Attribute        { return Attr("hx-get", v) }
func HxPost(v string) html.Attribute       { return Attr("hx-post", v) }
func HxPushUrl(v string) html.Attribute    { return Attr("hx-push-url", v) }
func HxPut(v string) html.Attribute        { return Attr("hx-put", v) }
func HxPatch(v string) html.Attribute      { return Attr("hx-patch", v) }
func HxReplaceUrl(v string) html.Attribute { return Attr("hx-replace-url", v) }
func HxSelect(v string) html.Attribute     { return Attr("hx-select", v) }
func HxSelectOob(v string) html.Attribute  { return Attr("hx-select-oob", v) }
func HxSwap(v string) html.Attribute       { return Attr("hx-swap", v) }
func HxSwapOob(v string) html.Attribute    { return Attr("hx-swap-oob", v) }
func HxTarget(v string) html.Attribute     { return Attr("hx-target", v) }
func HxTrigger(v string) html.Attribute    { return Attr("hx-trigger", v) }
func HxVals(v string) html.Attribute       { return Attr("hx-vals", v) }
