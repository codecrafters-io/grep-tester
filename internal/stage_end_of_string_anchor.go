package internal

import (
	tester_utils "github.com/codecrafters-io/tester-utils"
)

func testEndOfStringAnchor(stageHarness *tester_utils.StageHarness) error {
	testCases := []TestCase{
		{
			Pattern:          "cat$",
			Input:            "cat",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "cat$",
			Input:            "cats",
			ExpectedExitCode: 1,
		},
	}

	return RunTestCases(testCases, stageHarness)
}
