package internal

import "github.com/codecrafters-io/tester-utils/test_case_harness"

func testPositiveCharacterGroups(stageHarness *test_case_harness.TestCaseHarness) error {
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
		{
			Pattern:          "[abcd]",
			Input:            "[]",
			ExpectedExitCode: 1,
		},
	}

	return RunTestCases(testCases, stageHarness)
}
