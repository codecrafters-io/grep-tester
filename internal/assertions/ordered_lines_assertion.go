package assertions

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/tester-utils/executable"
	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/dustin/go-humanize/english"
)

type OrderedLinesAssertion struct {
	ExpectedOutputLines []string
}

func (a OrderedLinesAssertion) Run(result executable.ExecutableResult, logger *logger.Logger) error {
	actualOutput := strings.TrimSpace(string(result.Stdout))
	actualOutputLines := strings.Split(actualOutput, "\n")

	if actualOutput == "" {
		actualOutputLines = []string{}
	}

	// Assert each expected line in order
	for i, expectedLine := range a.ExpectedOutputLines {
		if i < len(actualOutputLines) {
			if actualOutputLines[i] != expectedLine {
				// I'll remove this comment later:
				// Should we check for occurence of ANSI escape sequences in actualOutputLines[i]
				// And when the expected string matches the actual string but differ in highlighting,
				// hint that the output should not contain highlighting
				// Or do we leave it this way given that highlighting is out of scope for this extension
				// And just display the escpae codes like we're currently doing (See internal/test_helpers/fixtures/printing_matches/highlighted)
				return fmt.Errorf("Expected line #%d to be %q, got %q", i+1, expectedLine, actualOutputLines[i])
			}

			logger.Successf("✓ Found line %q", expectedLine)
		} else {
			// We've run out of actual output lines
			missingLines := []string{}
			for j := i; j < len(a.ExpectedOutputLines); j++ {
				missingLines = append(missingLines, fmt.Sprintf("  %q", a.ExpectedOutputLines[j]))
			}
			return fmt.Errorf("Expected %s, missing %s:\n%s",
				english.Plural(len(a.ExpectedOutputLines), "line", "lines"),
				english.Plural(len(a.ExpectedOutputLines)-i, "line", "lines"),
				strings.Join(missingLines, "\n"))
		}
	}

	// Check for extra lines after all expected lines
	if len(actualOutputLines) > len(a.ExpectedOutputLines) {
		extraLines := []string{}
		for i := len(a.ExpectedOutputLines); i < len(actualOutputLines); i++ {
			extraLines = append(extraLines, fmt.Sprintf("  %q", actualOutputLines[i]))
		}

		return fmt.Errorf("Expected %s, found %s:\n%s",
			english.Plural(len(a.ExpectedOutputLines), "line", "lines"),
			english.Plural(len(extraLines), "extra line", "extra lines"),
			strings.Join(extraLines, "\n"))
	}

	logger.Successf("✓ Stdout contains %s in order", english.Plural(len(a.ExpectedOutputLines), "expected line", "expected lines"))
	return nil
}
