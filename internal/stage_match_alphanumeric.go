package internal

import (
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testMatchAlphanumeric(stageHarness *test_case_harness.TestCaseHarness) error {
	testCases := []TestCase{
		{
			Pattern:          `\w`,
			Input:            "word",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          `\w`,
			Input:            "$!?",
			ExpectedExitCode: 1,
		},
	}

	return RunTestCases(testCases, stageHarness)
}
