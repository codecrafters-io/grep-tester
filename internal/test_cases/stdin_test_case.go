package test_cases

import (
	"fmt"
	"path"

	"github.com/codecrafters-io/grep-tester/internal/grep"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

type StdinTestCase struct {
	Pattern          string
	Input            string
	ExpectedExitCode int
}

type StdinTestCaseCollection []StdinTestCase

func (c StdinTestCaseCollection) Run(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	executable := stageHarness.Executable

	for _, testCase := range c {
		logger.Infof("$ echo -n '%s' | ./%s -E '%s'", testCase.Input, path.Base(executable.Path), testCase.Pattern)

		expectedResult := grep.EmulateGrep([]string{"-E", testCase.Pattern}, []byte(testCase.Input))
		actualResult, err := executable.RunWithStdin([]byte(testCase.Input), "-E", testCase.Pattern)
		if err != nil {
			return err
		}

		if testCase.ExpectedExitCode != expectedResult.ExitCode {
			panic(fmt.Sprintf("CodeCrafters Internal Error: Expected exit code %v, grep returned %v", testCase.ExpectedExitCode, expectedResult.ExitCode))
		}
		if actualResult.ExitCode != testCase.ExpectedExitCode {
			return fmt.Errorf("Expected exit code %v, got %v", testCase.ExpectedExitCode, actualResult.ExitCode)
		}

		logger.Successf("✓ Received exit code %d.", actualResult.ExitCode)
	}

	return nil
}
