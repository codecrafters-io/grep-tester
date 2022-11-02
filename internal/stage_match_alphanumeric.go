package internal

import (
	tester_utils "github.com/codecrafters-io/tester-utils"
)

func testMatchAlphanumeric(stageHarness *tester_utils.StageHarness) error {
	testCases := []TestCase{
		{
			Pattern:          `\w`,
			Input:            "word",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          `\w`,
			Input:            "$!?",
			ExpectedExitCode: 1,
		},
	}

	return RunTestCases(testCases, stageHarness)
}
