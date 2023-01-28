package ui

import (
	"embed"
	_ "embed"
	"github.com/dustin/go-humanize"
	"html/template"
	"net/url"
	"strings"
)

func printUrlPath(path ...string) string {
	str := strings.Join(path, "")
	return (&url.URL{Path: str}).String()
}

var FuncMap = template.FuncMap{
	"humanizeIBytes": humanize.IBytes,
	"printUrlPath":   printUrlPath,
}

//go:embed ui.html
var uiHtml string

//go:embed set_rules.html
var setRulesHtml string

//go:embed static
var Static embed.FS

//go:embed static/images/favicon.ico
var Favicon []byte

var StatusTpl *template.Template

func init() {
	tpl, _ := template.New(UiName).Funcs(FuncMap).Parse(uiHtml)
	StatusTpl = template.Must(tpl.New(SetRulesName).Parse(setRulesHtml))
}
