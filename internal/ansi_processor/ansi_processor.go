package ansi_processor

import (
	"github.com/charmbracelet/x/ansi"
)

// AnsiProcessor is used to process ansi escape codes (https://en.wikipedia.org/wiki/ANSI_escape_code)
// in a string and return only the printable text characters
type AnsiProcessor struct {
}

func NewAnsiProcessor() *AnsiProcessor {
	return &AnsiProcessor{}
}

func (a *AnsiProcessor) StripAnsi(input string) string {
	return ansi.Strip(input)
}
