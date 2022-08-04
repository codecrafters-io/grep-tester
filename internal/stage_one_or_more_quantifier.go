package internal

import (
	tester_utils "github.com/codecrafters-io/tester-utils"
)

func testOneOrMoreQuantifier(stageHarness *tester_utils.StageHarness) error {
	testCases := []TestCase{
		{
			Pattern:          "ca+t",
			Input:            "caaats",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "ca+t",
			Input:            "cat",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "ca+t",
			Input:            "act",
			ExpectedExitCode: 1,
		},
	}

	return RunTestCases(testCases, stageHarness)
}
