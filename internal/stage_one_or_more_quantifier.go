package internal

import (
	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testOneOrMoreQuantifier(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	testCases := test_cases.StdinTestCaseCollection{
		{
			Pattern: "ca+t",
			Input:   "cat",
		},
		{
			Pattern: "ca+at",
			Input:   "caaats",
		},
		{
			Pattern: "ca+t",
			Input:   "act",
		},
		{
			Pattern: "ca+t",
			Input:   "ca",
		},
	}

	return testCases.Run(stageHarness)
}
