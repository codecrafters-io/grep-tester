package internal

import (
	"fmt"
	tester_utils "github.com/codecrafters-io/tester-utils"
)

type TestCase struct {
	Name             string
	Pattern          string
	Input            string
	ExpectedExitCode int
}

func testInit(stageHarness *tester_utils.StageHarness) error {
	logger := stageHarness.Logger
	executable := stageHarness.Executable

	testCases := []TestCase{
		{
			Name:             "character is present in string",
			Pattern:          "d",
			Input:            "this input contains the character d",
			ExpectedExitCode: 0,
		},
		{
			Name:             "character is not present in string",
			Pattern:          "f",
			Input:            "does not include the character",
			ExpectedExitCode: 1,
		},
	}

	for _, testCase := range testCases {
		logger.Infof("$ echo \"%s\" | ./your_grep.sh \"%s\"", testCase.Input, testCase.Pattern)
		result, err := executable.RunWithStdin([]byte(testCase.Input), testCase.Pattern)
		if err != nil {
			return err
		}

		if result.ExitCode != testCase.ExpectedExitCode {
			return fmt.Errorf("expected exit code %v, got %v", testCase.ExpectedExitCode, result.ExitCode)
		}

		logger.Successf("✓ Received exit code %d.", testCase.ExpectedExitCode)
	}

	return nil
}
