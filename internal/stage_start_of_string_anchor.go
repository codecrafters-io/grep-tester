package internal

import (
	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testStartOfStringAnchor(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	testCaseCollection := test_cases.StdinTestCaseCollection{
		{
			Pattern: "^log",
			Input:   "log",
		},
		{
			Pattern: "^log",
			Input:   "slog",
		},
	}

	return testCaseCollection.Run(stageHarness)
}
