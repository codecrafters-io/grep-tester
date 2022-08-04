package internal

import tester_utils "github.com/codecrafters-io/tester-utils"

func testPositiveCharacterGroups(stageHarness *tester_utils.StageHarness) error {
	testCases := []TestCase{
		{
			Pattern:          "[abcd]",
			Input:            "a",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "[abcd]",
			Input:            "efgh",
			ExpectedExitCode: 1,
		},
	}

	return RunTestCases(testCases, stageHarness)
}
