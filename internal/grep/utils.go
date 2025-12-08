package grep

import "strings"

// linesToProgramOutput converts lines produced by emulated grep to the actual output string
func linesToProgramOutput(lines []string) string {
	result := strings.Join(lines, "\n")

	// Add newline in case of non-empty output
	if result != "" {
		result += "\n"
	}

	return result
}
