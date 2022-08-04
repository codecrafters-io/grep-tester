package internal

import (
	tester_utils "github.com/codecrafters-io/tester-utils"
)

func testMatchDigit(stageHarness *tester_utils.StageHarness) error {
	testCases := []TestCase{
		{
			Pattern:          "\\d",
			Input:            "contains a number: 1",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "\\d",
			Input:            "does not contain a number",
			ExpectedExitCode: 1,
		},
	}

	return RunTestCases(testCases, stageHarness)
}
