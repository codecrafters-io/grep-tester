package internal

import (
	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testWildcard(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	testCases := test_cases.StdinTestCaseCollection{
		{
			Pattern: "c.t",
			Input:   "cat",
		},
		{
			Pattern: "c.t",
			Input:   "car",
		},
		{
			Pattern: "g.+gol",
			Input:   "goøö0Ogol",
		},
		{
			Pattern: "g.+gol",
			Input:   "gol",
		},
	}

	return testCases.Run(stageHarness)
}
