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

//go:embed fs.html
var uiHtml string

//go:embed static
var Static embed.FS

//go:embed static/images/favicon.ico
var Favicon []byte

var StatusTpl *template.Template

func init() {
	StatusTpl, _ = template.New(FSName).Funcs(FuncMap).Parse(uiHtml)
}
