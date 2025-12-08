package grep

import "strings"

// linesToProgramOutput converts lines produced by emulated grep to the actual output string
func linesToProgramOutput(lines []string, appendNewline bool) string {
	result := strings.Join(lines, "\n")

	if appendNewline {
		result += "\n"
	}

	return result
}
