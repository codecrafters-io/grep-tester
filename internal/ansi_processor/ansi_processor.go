package ansi_processor

import (
	"github.com/codecrafters-io/grep-tester/internal/virtual_terminal"
)

// AnsiProcessor is used to process ansi escape codes (https://en.wikipedia.org/wiki/ANSI_escape_code)
// in a string and return only the printable text characters
type AnsiProcessor struct {
}

func NewAnsiProcessor() *AnsiProcessor {
	return &AnsiProcessor{}
}

func (a *AnsiProcessor) Evaluate(input string) string {
	virtualTerminal := virtual_terminal.NewStandardVT()
	virtualTerminal.Write([]byte(input))
	return virtualTerminal.RenderText()
}
