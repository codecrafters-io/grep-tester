package internal

import (
	tester_utils "github.com/codecrafters-io/tester-utils"
)

func testInit(stageHarness *tester_utils.StageHarness) error {
	testCases := []TestCase{
		{
			Pattern:          "d",
			Input:            "dog",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "f",
			Input:            "dog",
			ExpectedExitCode: 1,
		},
	}

	return RunTestCases(testCases, stageHarness)
}
