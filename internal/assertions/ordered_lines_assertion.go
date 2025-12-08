package assertions

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/grep-tester/internal/utils"
	"github.com/codecrafters-io/tester-utils/executable"
	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/dustin/go-humanize/english"
)

type OrderedLinesAssertion struct {
	ExpectedOutputLines []string
}

func (a OrderedLinesAssertion) Run(result executable.ExecutableResult, logger *logger.Logger) error {
	actualOutput := string(result.Stdout)
	actualOutputLines := utils.ProgramOutputToLines(actualOutput)

	// Assert each expected line in order
	for i, expectedLine := range a.ExpectedOutputLines {
		if len(actualOutputLines) <= i {
			missingLinesErrorMessages := []string{}

			for j := i; j < len(a.ExpectedOutputLines); j++ {
				missingLine := (a.ExpectedOutputLines[j])
				missingLinesErrorMessages = append(missingLinesErrorMessages, fmt.Sprintf("⨯ %s", utils.FormatLineForLogging(missingLine)))
			}

			return fmt.Errorf("Expected %s in total, missing %s:\n%s",
				english.Plural(len(a.ExpectedOutputLines), "line", "lines"),
				english.Plural(len(a.ExpectedOutputLines)-i, "line", "lines"),
				strings.Join(missingLinesErrorMessages, "\n"))
		}

		if actualOutputLines[i] != expectedLine {
			return fmt.Errorf("Expected line #%d to be %s, got %s", i+1,
				utils.FormatLineForLogging(expectedLine), utils.FormatLineForLogging(actualOutputLines[i]))
		}

		logger.Successf("✓ Found line %s", utils.FormatLineForLogging(expectedLine))
	}

	// Check for extra lines after all expected lines
	if len(actualOutputLines) > len(a.ExpectedOutputLines) {
		extraLinesErrorMessages := []string{}
		for i := len(a.ExpectedOutputLines); i < len(actualOutputLines); i++ {
			extraLinesErrorMessages = append(extraLinesErrorMessages, fmt.Sprintf("⨯ %s", utils.FormatLineForLogging(actualOutputLines[i])))
		}

		// Better formatting for no output case
		if len(a.ExpectedOutputLines) == 0 {
			return fmt.Errorf(
				"Expected no output, got %s:\n%s",
				english.Plural(len(extraLinesErrorMessages), "line", "lines"),
				strings.Join(extraLinesErrorMessages, "\n"),
			)
		}

		return fmt.Errorf("Expected %s in total, found %s:\n%s",
			english.Plural(len(a.ExpectedOutputLines), "line", "lines"),
			english.Plural(len(extraLinesErrorMessages), "extra line", "extra lines"),
			strings.Join(extraLinesErrorMessages, "\n"))
	}

	if len(a.ExpectedOutputLines) == 0 {
		logger.Successf("✓ No output found")
	}

	if len(a.ExpectedOutputLines) > 1 {
		logger.Successf("✓ Stdout contains %s in order", english.Plural(len(a.ExpectedOutputLines), "expected line", "expected lines"))
	}

	return nil
}
