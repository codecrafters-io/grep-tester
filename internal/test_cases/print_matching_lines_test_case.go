package test_cases

import (
	"fmt"
	"path"
	"strings"

	"github.com/codecrafters-io/grep-tester/internal/ansi_processor"
	"github.com/codecrafters-io/grep-tester/internal/grep"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
	"github.com/dustin/go-humanize/english"
)

type PrintMatchingLinesTestCase struct {
	Pattern             string
	InputLines          []string
	ExpectedExitCode    int
	ExpectedOutputLines []string
}

type PrintMatchingLinesTestCaseCollection []PrintMatchingLinesTestCase

func (c PrintMatchingLinesTestCaseCollection) Run(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	executable := stageHarness.Executable

	for _, testCase := range c {
		allInputLines := strings.Join(testCase.InputLines, "\n")

		logger.Infof("$ echo -ne %q | ./%s -E '%s'", allInputLines, path.Base(executable.Path), testCase.Pattern)

		expectedResult := grep.EmulateGrep([]string{"-E", testCase.Pattern}, []byte(allInputLines))
		actualResult, err := executable.RunWithStdin([]byte(allInputLines), "-E", testCase.Pattern)

		if err != nil {
			return err
		}

		// Compare exit codes
		if testCase.ExpectedExitCode != expectedResult.ExitCode {
			panic(fmt.Sprintf("CodeCrafters Internal Error: Expected exit code %v, grep returned %v", testCase.ExpectedExitCode, expectedResult.ExitCode))
		}

		if actualResult.ExitCode != testCase.ExpectedExitCode {
			return fmt.Errorf("Expected exit code %v, got %v.\nHint: %s", testCase.ExpectedExitCode, actualResult.ExitCode, getRegex101Link(testCase.Pattern, allInputLines))
		}

		logger.Successf("✓ Received exit code %d.", actualResult.ExitCode)

		// Compare stdout text
		// Grep will never produce highlighted result since it is piped to our tester
		// However, the tester still needs to be lenient in case the user has implemented highlighting on their own
		actualStdoutText := removeTrailingNewline(
			ansi_processor.NewAnsiProcessor().Evaluate(string(actualResult.Stdout)),
		)

		stdoutTextFromGrep := string(expectedResult.Stdout)
		expectedStdoutText := strings.Join(testCase.ExpectedOutputLines, "\n")

		// Compare against grep
		if stdoutTextFromGrep != expectedStdoutText {
			panic(fmt.Sprintf("Codecrafters Internal Error: Expected output text: %q, grep returned %q", expectedStdoutText, stdoutTextFromGrep))
		}

		// No output case
		if testCase.ExpectedOutputLines == nil {
			if actualStdoutText != "" {
				return fmt.Errorf("Expected no output, got %q", actualStdoutText)
			}
			logger.Successf("✓ Received no output")
			continue
		}

		actualStdoutLines := strings.Split(actualStdoutText, "\n")

		// Compare length
		if len(testCase.ExpectedOutputLines) != len(actualStdoutLines) {
			return fmt.Errorf(
				"Expected %s:\n%s\n\nGot %s:\n%s\n",
				english.Plural(len(testCase.ExpectedOutputLines), "line", "lines"),
				expectedStdoutText,
				english.Plural(len(actualStdoutLines), "line", "lines"),
				actualStdoutText,
			)
		}

		// Compare each line in the output
		for i, expectedOutputLine := range testCase.ExpectedOutputLines {
			if expectedOutputLine != actualStdoutLines[i] {
				return fmt.Errorf("Expected line #%d of output to be %q, got %q", i+1, expectedOutputLine, actualStdoutLines[i])
			}
			logger.Successf("✓ Received line %q in output", expectedOutputLine)
		}
	}

	return nil
}
