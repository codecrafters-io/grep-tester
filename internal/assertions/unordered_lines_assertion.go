package assertions

import (
	"fmt"
	"slices"
	"strings"

	"github.com/codecrafters-io/grep-tester/internal/utils"
	"github.com/codecrafters-io/tester-utils/executable"
	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/dustin/go-humanize/english"
)

type UnorderedLinesAssertion struct {
	ExpectedOutputLines []string
}

func (a UnorderedLinesAssertion) Run(result executable.ExecutableResult, logger *logger.Logger) error {
	actualOutput := string(result.Stdout)
	actualOutputLines := utils.ProgramOutputToLines(actualOutput)

	foundLines := []string{}
	missingLines := []string{}
	extraLines := []string{}

	// Collect found and missing lines
	for _, expectedLine := range a.ExpectedOutputLines {
		if slices.Contains(actualOutputLines, expectedLine) {
			foundLines = append(foundLines, expectedLine)
		} else {
			missingLines = append(missingLines, expectedLine)
		}
	}

	// Collect extra lines
	for _, actualLine := range actualOutputLines {
		if !slices.Contains(a.ExpectedOutputLines, actualLine) {
			extraLines = append(extraLines, actualLine)
		}
	}

	// Success case
	if len(missingLines) == 0 && len(extraLines) == 0 && len(foundLines) == len(a.ExpectedOutputLines) {
		if len(foundLines) == 0 {
			logger.Successf("✓ No output found")
		} else {
			logger.Successf(
				"✓ Stdout contains %s",
				english.Plural(len(a.ExpectedOutputLines), "expected line", "expected lines"),
			)
		}

		return nil
	}

	// Failure case
	// Display all found lines first
	for _, line := range foundLines {
		logger.Successf("✓ Found line %s", utils.FormatLineForLogging(line))
	}

	// Prioritize errors related to missing lines
	if len(missingLines) > 0 {
		missingLinesErrorMessages := []string{}

		for _, missingLine := range missingLines {
			missingLinesErrorMessages = append(missingLinesErrorMessages, fmt.Sprintf("⨯ %s", utils.FormatLineForLogging(missingLine)))
		}

		return fmt.Errorf(
			"Expected %s in output, only found %s. Missing %s:\n%s",
			english.Plural(len(a.ExpectedOutputLines), "line", "lines"),
			english.Plural(len(foundLines), "matching line", "matching lines"),
			english.PluralWord(len(missingLines), "match", "matches"),
			strings.Join(missingLinesErrorMessages, "\n"),
		)
	}

	// Print errors related to extra lines at last
	if len(extraLines) > 0 {
		extraLineErrorMessages := []string{}

		for _, extraLine := range extraLines {
			extraLineErrorMessages = append(extraLineErrorMessages, fmt.Sprintf("⨯ %s", utils.FormatLineForLogging(extraLine)))
		}

		// Better formatting for no output case
		if len(a.ExpectedOutputLines) == 0 {
			return fmt.Errorf(
				"Expected no output, got %s:\n%s",
				english.Plural(len(extraLines), "line", "lines"),
				strings.Join(extraLineErrorMessages, "\n"),
			)
		}

		return fmt.Errorf(
			"Expected %s in output, found %d. Unexpected %s:\n%s",
			english.Plural(len(a.ExpectedOutputLines), "line", "lines"),
			len(actualOutputLines),
			english.PluralWord(len(extraLines), "line", "lines"),
			strings.Join(extraLineErrorMessages, "\n"),
		)
	}

	return nil
}
