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

	// If there's no output, we need to correct the behavior of strings.Split
	if actualOutput == "" {
		actualOutputLines = []string{}
	}

	// Track which expected lines we've found in order
	foundLinesCount := 0
	extraLines := []string{}

	for _, actualLine := range actualOutputLines {
		if foundLinesCount < len(a.ExpectedOutputLines) && actualLine == a.ExpectedOutputLines[foundLinesCount] {
			logger.Successf("✓ Found line '%s'", actualLine)
			foundLinesCount++
		} else {
			extraLines = append(extraLines, actualLine)
		}
	}

	// All expected lines were found: There are no missing lines
	if foundLinesCount == len(a.ExpectedOutputLines) {
		// No extra lines were found
		if len(extraLines) == 0 {
			logger.Successf(
				"✓ Stdout contains %s",
				english.Plural(len(a.ExpectedOutputLines), "expected line", "expected lines"),
			)
			return nil
		}

		// Extra lines were found in addition to expected lines
		logger.Infof(
			"Expected %s in output, found %d total. Unexpected %s: ",
			english.Plural(len(a.ExpectedOutputLines), "line", "lines"),
			len(actualOutputLines),
			english.Plural(len(extraLines), "line", "lines"),
		)

		errorMessage := []string{}
		for _, extraLine := range extraLines {
			errorMessage = append(errorMessage, fmt.Sprintf("⨯ Extra line found: \"%s\"", extraLine))
		}

		return fmt.Errorf("%s", strings.Join(errorMessage, "\n"))
	}

	// There are missing lines
	logger.Infof(
		"Expected %s in output in order, only found %s. Missing %s:",
		english.Plural(len(a.ExpectedOutputLines), "line", "lines"),
		english.Plural(foundLinesCount, "line", "lines"),
		english.PluralWord(len(a.ExpectedOutputLines)-foundLinesCount, "line", "lines"),
	)

	errorMessage := []string{}
	for i := foundLinesCount; i < len(a.ExpectedOutputLines); i++ {
		errorMessage = append(errorMessage, fmt.Sprintf("⨯ Line not found: \"%s\"", a.ExpectedOutputLines[foundLinesCount]))
	}

	if len(extraLines) > 0 {
		errorMessage = append(errorMessage, fmt.Sprintf("\nExtra %s encountered:", english.PluralWord(len(extraLines), "line", "lines")))
		for _, extraLine := range extraLines {
			errorMessage = append(errorMessage, fmt.Sprintf("⨯ Extra line found: \"%s\"", extraLine))
		}
	}
	return fmt.Errorf("%s", strings.Join(errorMessage, "\n"))
}
