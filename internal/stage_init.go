package internal

import (
	tester_utils "github.com/codecrafters-io/tester-utils"
)

func testInit(stageHarness *tester_utils.StageHarness) error {
	testCases := []TestCase{
		{
			Pattern:          "d",
			Input:            "this input contains the character d",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "f",
			Input:            "does not include the character",
			ExpectedExitCode: 1,
		},
	}

	return RunTestCases(testCases, stageHarness)
}
