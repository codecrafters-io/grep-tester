package test_cases

import (
	"fmt"
	"path"

	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

type StdinTestCase struct {
	Pattern          string
	Input            string
	ExpectedExitCode int
}

func (testCase StdinTestCase) Run(stageHarness *test_case_harness.TestCaseHarness) error {
	return runStdinTestCases([]StdinTestCase{testCase}, stageHarness)
}

type StdinTestCases []StdinTestCase

func (testCases StdinTestCases) Run(stageHarness *test_case_harness.TestCaseHarness) error {
	return runStdinTestCases(testCases, stageHarness)
}

func runStdinTestCases(testCases []StdinTestCase, stageHarness *test_case_harness.TestCaseHarness) error {
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
