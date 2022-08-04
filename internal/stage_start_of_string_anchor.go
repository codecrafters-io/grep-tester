package internal

import (
	tester_utils "github.com/codecrafters-io/tester-utils"
)

func testStartOfStringAnchor(stageHarness *tester_utils.StageHarness) error {
	testCases := []TestCase{
		{
			Pattern:          "^log",
			Input:            "log",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "^log",
			Input:            "slog",
			ExpectedExitCode: 1,
		},
	}

	return RunTestCases(testCases, stageHarness)
}
