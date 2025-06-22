package helper

import (
	"html"
	"strings"
)

func SanitizeInput(input string) string {
	return strings.TrimSpace(html.EscapeString(input))
}