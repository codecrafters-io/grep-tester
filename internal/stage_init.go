package internal

import (
	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testInit(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	testCases := test_cases.StdinTestCaseCollection{
		{
			Pattern: "d",
			Input:   "dog",
		},
		{
			Pattern: "f",
			Input:   "dog",
		},
	}

	return testCases.Run(stageHarness)
}
