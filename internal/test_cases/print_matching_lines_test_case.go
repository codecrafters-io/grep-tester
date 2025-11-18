package test_cases

import (
	"fmt"
	"path"
	"strings"

	"github.com/codecrafters-io/grep-tester/internal/assertions"
	"github.com/codecrafters-io/grep-tester/internal/grep"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
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
		// Run executable and collect result
		allInputLines := strings.Join(testCase.InputLines, "\n")
		logger.Infof("$ echo -ne %q | ./%s -E '%s'", allInputLines, path.Base(executable.Path), testCase.Pattern)

		grepResult := grep.EmulateGrep([]string{"-E", testCase.Pattern}, []byte(allInputLines))
		actualResult, err := executable.RunWithStdin([]byte(allInputLines), "-E", testCase.Pattern)

		if err != nil {
			return err
		}

		// Compare against grep

		if testCase.ExpectedExitCode != grepResult.ExitCode {
			panic(fmt.Sprintf("CodeCrafters Internal Error: Expected exit code %v, grep returned %v", testCase.ExpectedExitCode, grepResult.ExitCode))
		}

		expectedOutput := strings.Join(testCase.ExpectedOutputLines, "\n")

		if expectedOutput != string(grepResult.Stdout) {
			panic(fmt.Sprintf("Codecrafters Internal Error: Expected stdout: %q, grep's stdout: %q", expectedOutput, grepResult.Stdout))
		}

		// Run assertions

		exitCodeAssertion := assertions.ExitCodeAssertion{
			ExpectedExitCode: testCase.ExpectedExitCode,
		}

		if err := exitCodeAssertion.Run(actualResult, logger); err != nil {
			return err
		}

		orderedLinesAssertion := assertions.OrderedLinesAssertion{
			ExpectedOutputLines: testCase.ExpectedOutputLines,
		}

		if err := orderedLinesAssertion.Run(actualResult, logger); err != nil {
			return err
		}
	}

	return nil
}
