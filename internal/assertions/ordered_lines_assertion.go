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

	actualOutputLines := strings.FieldsFunc(actualOutput, func(r rune) bool {
		return r == '\n'
	})

	// Assert each expected line in order
	for i, expectedLine := range a.ExpectedOutputLines {
		if len(actualOutputLines) <= i {
			missingLines := []string{}

			for j := i; j < len(a.ExpectedOutputLines); j++ {
				missingLine := (a.ExpectedOutputLines[j])
				missingLines = append(missingLines, fmt.Sprintf("  %q", missingLine))
			}

			return fmt.Errorf("Expected %s in total, missing %s:\n%s",
				english.Plural(len(a.ExpectedOutputLines), "line", "lines"),
				english.Plural(len(a.ExpectedOutputLines)-i, "line", "lines"),
				strings.Join(missingLines, "\n"))
		}

		if actualOutputLines[i] != expectedLine {
			return fmt.Errorf("Expected line #%d to be %q, got %q", i+1, expectedLine, actualOutputLines[i])
		}

		logger.Successf("✓ Found line %q", expectedLine)
	}

	// Check for extra lines after all expected lines
	if len(actualOutputLines) > len(a.ExpectedOutputLines) {
		extraLines := []string{}
		for i := len(a.ExpectedOutputLines); i < len(actualOutputLines); i++ {
			extraLines = append(extraLines, fmt.Sprintf("  %q", actualOutputLines[i]))
		}

		// Better formatting for no output case
		if len(a.ExpectedOutputLines) == 0 {
			return fmt.Errorf(
				"Expected no output, got %s:\n%s",
				english.Plural(len(extraLines), "line", "lines"),
				strings.Join(extraLines, "\n"),
			)
		}

		return fmt.Errorf("Expected %s in total, found %s:\n%s",
			english.Plural(len(a.ExpectedOutputLines), "line", "lines"),
			english.Plural(len(extraLines), "extra line", "extra lines"),
			strings.Join(extraLines, "\n"))
	}

	if len(a.ExpectedOutputLines) == 0 {
		logger.Successf("✓ Stdout contains no output")
	} else {
		logger.Successf("✓ Stdout contains %s in order", english.Plural(len(a.ExpectedOutputLines), "expected line", "expected lines"))
	}

	return nil
}
