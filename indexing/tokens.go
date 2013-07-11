package indexing

import (
	"bytes"
	"io"
	"strings"
	"unicode"
)

func Tokenize(r io.Reader) []string {
	buf := new(bytes.Buffer)
	buf.ReadFrom(r)

  return TokenizeString(buf.String())
}

func TokenizeString(s string) []string {
	fields := strings.Fields(s)
	var filtered []string

	for _, field := range fields {
		filtered = append(filtered, strings.ToLower(sanitize(field)))
	}
	return filtered
}

func sanitize(s string) string {
	buf := new(bytes.Buffer)
	for _, rne := range s {
		if !unicode.IsPunct(rne) {
			buf.WriteRune(rne)
		}
	}
	return buf.String()
}
