package internal

import (
	tester_utils "github.com/codecrafters-io/tester-utils"
)

func testMatchDigit(stageHarness *tester_utils.StageHarness) error {
	testCases := []TestCase{
		{
			Pattern:          `\d`,
			Input:            "123",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          `\d`,
			Input:            "apple",
			ExpectedExitCode: 1,
		},
	}

	return RunTestCases(testCases, stageHarness)
}
