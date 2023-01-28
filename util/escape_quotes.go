package util

import (
	"strings"
)

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func EscapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}
