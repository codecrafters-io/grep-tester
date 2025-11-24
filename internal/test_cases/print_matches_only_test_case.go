package test_cases

import (
	"fmt"
	"path"
	"strings"

	"github.com/codecrafters-io/grep-tester/internal/assertions"
	"github.com/codecrafters-io/grep-tester/internal/grep"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

type PrintMatchesOnlyTestCase struct {
	Pattern          string
	InputLines       []string
	ExpectedExitCode int
}

type PrintMatchesOnlyTestCaseCollection []PrintMatchesOnlyTestCase

func (c PrintMatchesOnlyTestCaseCollection) Run(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	executable := stageHarness.Executable

	for _, testCase := range c {
		// Run executable and collect result
		allInputLines := strings.Join(testCase.InputLines, "\n")
		logger.Infof("$ echo -ne %q | ./%s -o -E '%s'", allInputLines, path.Base(executable.Path), testCase.Pattern)

		grepResult := grep.EmulateGrep([]string{"-o", "-E", testCase.Pattern}, []byte(allInputLines))
		actualResult, err := executable.RunWithStdin([]byte(allInputLines), "-o", "-E", testCase.Pattern)

		if err != nil {
			return err
		}

		// Compare against grep
		if testCase.ExpectedExitCode != grepResult.ExitCode {
			panic(fmt.Sprintf("CodeCrafters Internal Error: Expected exit code %v, grep returned %v", testCase.ExpectedExitCode, grepResult.ExitCode))
		}

		// Run assertions
		exitCodeAssertion := assertions.ExitCodeAssertion{
			ExpectedExitCode: testCase.ExpectedExitCode,
		}

		if err := exitCodeAssertion.Run(actualResult, logger); err != nil {
			return err
		}

		expectedOutput := strings.TrimSpace(string(grepResult.Stdout))

		expectedOutputLines := strings.FieldsFunc(expectedOutput, func(r rune) bool {
			return r == '\n'
		})

		orderedLinesAssertion := assertions.OrderedLinesAssertion{
			ExpectedOutputLines: expectedOutputLines,
		}

		if err := orderedLinesAssertion.Run(actualResult, logger); err != nil {
			return err
		}
	}

	return nil
}
