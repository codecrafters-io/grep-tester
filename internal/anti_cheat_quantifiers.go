package internal

import (
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testQuantifiersAsAntiCheat(stageHarness *test_case_harness.TestCaseHarness) error {
	testCases := []AntiCheatTestCase{
		{
			Pattern:          "a{1,2}bc",
			Input:            "abc",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "a{2,3}bc",
			Input:            "abc",
			ExpectedExitCode: 1,
		},
	}

	return RunAntiCheatTestCases(testCases, stageHarness)
}
