package internal

import (
	tester_utils "github.com/codecrafters-io/tester-utils"
)

func testAlternation(stageHarness *tester_utils.StageHarness) error {
	testCases := []TestCase{
		{
			Pattern:          "a (cat|dog)",
			Input:            "a cat",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "a (cat|dog)",
			Input:            "a dog",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "a (cat|dog)",
			Input:            "a cow",
			ExpectedExitCode: 1,
		},
	}

	return RunTestCases(testCases, stageHarness)
}
