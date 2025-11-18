package assertions

import (
	"fmt"
	"slices"
	"strings"

	"github.com/codecrafters-io/tester-utils/executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

type UnorderedLinesAssertion struct {
	ExpectedOutputLines []string
}

func (a UnorderedLinesAssertion) Run(result executable.ExecutableResult, logger *logger.Logger) error {
	actualOutput := strings.TrimSpace(string(result.Stdout))

	actualOutputLines := strings.Split(actualOutput, "\n")

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
		logger.Successf("✓ Stdout contains %d expected line(s)", len(a.ExpectedOutputLines))
	} else {
		for _, line := range foundLines {
			logger.Successf("✓ Found line '%s'", line)
		}

		if len(missingLines) > 0 {
			logger.Infof("Expected %d line(s) in output, only found %d matching line(s). Missing match(es):", len(a.ExpectedOutputLines), len(foundLines))
			errorMessage := []string{}
			for _, line := range missingLines {
				errorMessage = append(errorMessage, fmt.Sprintf("⨯ Line not found: \"%s\"", line))
			}
			return fmt.Errorf("%s", strings.Join(errorMessage, "\n"))
		}

		if len(extraLines) > 0 {
			logger.Infof("Expected %d line(s) in output, found %d. Unexpected line(s):", len(a.ExpectedOutputLines), len(actualOutputLines))
			errorMessage := []string{}
			for _, line := range extraLines {
				errorMessage = append(errorMessage, fmt.Sprintf("⨯ Extra line found: \"%s\"", line))
			}
			return fmt.Errorf("%s", strings.Join(errorMessage, "\n"))
		}
	}

	return nil
}
