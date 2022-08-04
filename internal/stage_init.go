package internal

import (
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
			Input:            "do not mention the letter after e",
			ExpectedExitCode: 1,
		},
	}

	for _, testCase := range testCases {
		logger.Infof("$ echo \"%s\" | ./your_grep.sh \"%s\"", testCase.Input, testCase.Pattern)
		result, err := executable.Run("abcd")
		if err != nil {
			return err
		}
	}

	return nil
}
