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

	// Assert each expected line in order
	for i, expectedLine := range a.ExpectedOutputLines {
		if len(actualOutputLines) <= i {
			missingLines := []string{}

			for j := i; j < len(a.ExpectedOutputLines); j++ {
				missingLine := escapeLine(a.ExpectedOutputLines[j])
				missingLines = append(missingLines, fmt.Sprintf("  %s", missingLine))
			}

			return fmt.Errorf("Expected %s, missing %s:\n%s",
				english.Plural(len(a.ExpectedOutputLines), "line", "lines"),
				english.Plural(len(a.ExpectedOutputLines)-i, "line", "lines"),
				strings.Join(missingLines, "\n"))
		}

		if actualOutputLines[i] != expectedLine {
			return fmt.Errorf("Expected line #%d to be %s, got %s", i+1, escapeLine(expectedLine), escapeLine(actualOutputLines[i]))
		}

		logger.Successf("✓ Found line %s", escapeLine(expectedLine))
	}

	// Check for extra lines after all expected lines
	if len(actualOutputLines) > len(a.ExpectedOutputLines) {
		extraLines := []string{}
		for i := len(a.ExpectedOutputLines); i < len(actualOutputLines); i++ {
			extraLines = append(extraLines, fmt.Sprintf("  %s", escapeLine(actualOutputLines[i])))
		}

		return fmt.Errorf("Expected %s, found %s:\n%s",
			english.Plural(len(a.ExpectedOutputLines), "line", "lines"),
			english.Plural(len(extraLines), "extra line", "extra lines"),
			strings.Join(extraLines, "\n"))
	}

	logger.Successf("✓ Stdout contains %s in order", english.Plural(len(a.ExpectedOutputLines), "expected line", "expected lines"))
	return nil
}
