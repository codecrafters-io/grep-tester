package internal

import "github.com/codecrafters-io/tester-utils/test_case_harness"

func testNegativeCharacterGroups(stageHarness *test_case_harness.TestCaseHarness) error {
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
		{
			Pattern:          "[^opq]",
			Input:            "orange",
			ExpectedExitCode: 0,
		},
	}

	return RunTestCases(testCases, stageHarness)
}
