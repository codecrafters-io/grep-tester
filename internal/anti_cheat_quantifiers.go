package internal

import (
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testQuantifiersAsAntiCheat(stageHarness *test_case_harness.TestCaseHarness) error {
	testCases := AntiCheatTestCaseCollection{
		{
			Pattern: "a{1,2}bc",
			Input:   "abc",
		},
		{
			Pattern: "a{2,3}bc",
			Input:   "abc",
		},
	}

	return testCases.Run(stageHarness)
}
