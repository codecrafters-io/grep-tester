package internal

import (
	tester_utils "github.com/codecrafters-io/tester-utils"
)

func testWildcard(stageHarness *tester_utils.StageHarness) error {
	testCases := []TestCase{
		{
			Pattern:          "c.t",
			Input:            "cat",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "c.t",
			Input:            "cot",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "c.t",
			Input:            "car",
			ExpectedExitCode: 1,
		},
	}

	return RunTestCases(testCases, stageHarness)
}
