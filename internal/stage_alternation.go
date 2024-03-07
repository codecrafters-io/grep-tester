package internal

import (
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testAlternation(stageHarness *test_case_harness.TestCaseHarness) error {
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
