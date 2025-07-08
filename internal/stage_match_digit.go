package internal

import (
	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testMatchDigit(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	testCases := test_cases.StdinTestCaseCollection{
		{
			Pattern: `\d`,
			Input:   "123",
		},
		{
			Pattern: `\d`,
			Input:   "apple",
		},
	}

	return testCases.Run(stageHarness)
}
