package internal

import (
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testStartOfStringAnchor(stageHarness *test_case_harness.TestCaseHarness) error {
	testCases := []TestCase{
		{
			Pattern:          "^log",
			Input:            "log",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "^log",
			Input:            "slog",
			ExpectedExitCode: 1,
		},
	}

	return RunTestCases(testCases, stageHarness)
}
