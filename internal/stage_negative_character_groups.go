package internal

import tester_utils "github.com/codecrafters-io/tester-utils"

func testNegativeCharacterGroups(stageHarness *tester_utils.StageHarness) error {
	testCases := []TestCase{
		{
			Pattern:          "[^xyz]",
			Input:            "apple",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "[^anb]",
			Input:            "banana",
			ExpectedExitCode: 1,
		},
	}

	return RunTestCases(testCases, stageHarness)
}
