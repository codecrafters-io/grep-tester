package test_cases

import (
	"fmt"
	"path"

	"github.com/codecrafters-io/grep-tester/internal/grep"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

type StdinTestCase struct {
	Pattern string
	Input   string
}

type StdinTestCaseCollection []StdinTestCase

func (c StdinTestCaseCollection) Run(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	executable := stageHarness.Executable

	for _, testCase := range c {
		logger.Infof("$ echo -n \"%s\" | ./%s -E \"%s\"", testCase.Input, path.Base(executable.Path), testCase.Pattern)

		expectedResult := grep.EmulateGrep([]string{"-E", testCase.Pattern}, []byte(testCase.Input))
		actualResult, err := executable.RunWithStdin([]byte(testCase.Input), "-E", testCase.Pattern)
		if err != nil {
			return err
		}

		if actualResult.ExitCode != expectedResult.ExitCode {
			return fmt.Errorf("Expected exit code %v, got %v", expectedResult.ExitCode, actualResult.ExitCode)
		}

		logger.Successf("✓ Received exit code %d.", expectedResult.ExitCode)
	}

	return nil
}
