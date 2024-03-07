package internal

import (
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testWildcard(stageHarness *test_case_harness.TestCaseHarness) error {
	testCases := []TestCase{
		{
			Pattern:          "c.t",
			Input:            "cat",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "c.t",
			Input:            "cot",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "c.t",
			Input:            "car",
			ExpectedExitCode: 1,
		},
	}

	return RunTestCases(testCases, stageHarness)
}
