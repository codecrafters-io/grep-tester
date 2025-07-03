package internal

import (
	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testStartOfStringAnchor(stageHarness *test_case_harness.TestCaseHarness) error {
	MoveGrepToTemp(stageHarness, stageHarness.Logger)

	testCases := []test_cases.StdinTestCase{
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

	return test_cases.RunStdinTestCases(testCases, stageHarness)
}
