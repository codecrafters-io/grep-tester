package internal

import (
	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testStartOfStringAnchor(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	testCases := test_cases.StdinTestCaseCollection{
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

	return testCases.Run(stageHarness)
}
