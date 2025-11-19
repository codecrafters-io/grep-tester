package assertions

import (
	"fmt"
	"slices"
	"strings"

	"github.com/codecrafters-io/tester-utils/executable"
	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/dustin/go-humanize/english"
)

type UnorderedLinesAssertion struct {
	ExpectedOutputLines []string
}

func (a UnorderedLinesAssertion) Run(result executable.ExecutableResult, logger *logger.Logger) error {
	actualOutput := strings.TrimSpace(string(result.Stdout))

	actualOutputLines := strings.FieldsFunc(actualOutput, func(r rune) bool {
		return r == '\n'
	})

	foundLines := []string{}
	missingLines := []string{}
	extraLines := []string{}

	for _, expectedLine := range a.ExpectedOutputLines {
		if slices.Contains(actualOutputLines, expectedLine) {
			foundLines = append(foundLines, expectedLine)
		} else {
			missingLines = append(missingLines, expectedLine)
		}
	}

	for _, actualLine := range actualOutputLines {
		if !slices.Contains(a.ExpectedOutputLines, actualLine) {
			extraLines = append(extraLines, actualLine)
		}
	}

	if len(missingLines) == 0 && len(extraLines) == 0 && len(foundLines) == len(a.ExpectedOutputLines) {
		if len(foundLines) == 0 {
			logger.Successf("✓ No output found")
		} else {
			logger.Successf("✓ Stdout contains %d expected line(s)", len(a.ExpectedOutputLines))
		}

		return nil
	}

	for _, line := range foundLines {
		logger.Successf("✓ Found line %q", line)
	}

	if len(missingLines) > 0 {
		logger.Errorf(
			"Expected %s in output, only found %s. Missing %s:",
			english.Plural(len(a.ExpectedOutputLines), "line", "lines"),
			english.Plural(len(foundLines), "matching line", "matching lines"),
			english.PluralWord(len(missingLines), "match", "matches"),
		)
		errorMessage := []string{}

		for _, line := range missingLines {
			errorMessage = append(errorMessage, fmt.Sprintf("  %q", line))
		}

		return fmt.Errorf("%s", strings.Join(errorMessage, "\n"))
	}

	if len(extraLines) > 0 {
		errorMessage := []string{}

		for _, line := range extraLines {
			errorMessage = append(errorMessage, fmt.Sprintf("  %q", line))
		}

		logger.Errorf(
			"Expected %s in output, found %d. Unexpected %s:",
			english.Plural(len(a.ExpectedOutputLines), "line", "lines"),
			len(actualOutputLines),
			english.PluralWord(len(extraLines), "line", "lines"),
		)

		return fmt.Errorf("%s", strings.Join(errorMessage, "\n"))
	}

	return nil
}
