package internal

import (
	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testZeroOrOneQuantifier(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	testCaseCollection := test_cases.StdinTestCaseCollection{
		{
			Pattern: "ca?t",
			Input:   "cat",
		},
		{
			Pattern: "ca?t",
			Input:   "act",
		},
		{
			Pattern: "ca?t",
			Input:   "dog",
		},
		{
			Pattern: "ca?t",
			Input:   "cag",
		},
	}

	return testCaseCollection.Run(stageHarness)
}
