package internal

import (
	tester_utils "github.com/codecrafters-io/tester-utils"
)

func testZeroOrOneQuantifier(stageHarness *tester_utils.StageHarness) error {
	testCases := []TestCase{
		{
			Pattern:          "ca?t",
			Input:            "cat",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "ca?t",
			Input:            "act",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "ca?t",
			Input:            "dog",
			ExpectedExitCode: 1,
		},
		{
			Pattern:          "ca?t",
			Input:            "cag",
			ExpectedExitCode: 1,
		},
	}

	return RunTestCases(testCases, stageHarness)
}
