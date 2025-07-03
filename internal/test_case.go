package internal

import (
	"fmt"
	"path"

	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

type TestCase struct {
	Pattern          string
	Input            string
	ExpectedExitCode int
}

func RunTestCases(testCases []TestCase, stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	executable := stageHarness.Executable

	for _, testCase := range testCases {
		logger.Infof("$ echo -n \"%s\" | ./%s -E \"%s\"", testCase.Input, path.Base(executable.Path), testCase.Pattern)
		result, err := executable.RunWithStdin([]byte(testCase.Input), "-E", testCase.Pattern)
		if err != nil {
			return err
		}

		if result.ExitCode != testCase.ExpectedExitCode {
			return fmt.Errorf("Expected exit code %v, got %v", testCase.ExpectedExitCode, result.ExitCode)
		}

		logger.Successf("✓ Received exit code %d.", testCase.ExpectedExitCode)
	}

	return nil
}
