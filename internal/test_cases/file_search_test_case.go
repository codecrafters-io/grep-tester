package test_cases

import (
	"fmt"
	"path"
	"strings"

	"github.com/codecrafters-io/grep-tester/internal/assertions"
	"github.com/codecrafters-io/grep-tester/internal/grep"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

type FileSearchTestCase struct {
	Pattern                   string
	FilePaths                 []string
	ExpectedExitCode          int
	ShouldEnableRecursiveFlag bool
}

type FileSearchTestCaseCollection []FileSearchTestCase

func (c FileSearchTestCaseCollection) Run(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	executable := stageHarness.Executable

	for _, testCase := range c {
		args := []string{}
		if testCase.ShouldEnableRecursiveFlag {
			args = append(args, "-r")
		}
		args = append(args, "-E", testCase.Pattern)
		args = append(args, testCase.FilePaths...)

		logger.Infof("$ ./%s %s", path.Base(executable.Path), strings.Join(args, " "))

		expectedResult := grep.EmulateGrep(args, []byte{})
		actualResult, err := executable.Run(args...)
		if err != nil {
			return err
		}

		if testCase.ExpectedExitCode != expectedResult.ExitCode {
			panic(fmt.Sprintf("CodeCrafters Internal Error: Expected exit code %v, grep returned %v", testCase.ExpectedExitCode, expectedResult.ExitCode))
		}

		exitCodeAssertion := assertions.ExitCodeAssertion{
			ExpectedExitCode: testCase.ExpectedExitCode,
		}

		if err := exitCodeAssertion.Run(actualResult, logger); err != nil {
			return err
		}

		expectedOutput := strings.TrimSpace(string(expectedResult.Stdout))
		expectedOutputLines := strings.Split(expectedOutput, "\n")

		unorderedLinesAssertion := assertions.UnorderedLinesAssertion{
			ExpectedOutputLines: expectedOutputLines,
		}

		if err := unorderedLinesAssertion.Run(actualResult, logger); err != nil {
			return err
		}

	}

	return nil
}
