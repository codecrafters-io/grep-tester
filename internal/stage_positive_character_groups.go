package internal

import (
	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testPositiveCharacterGroups(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	testCases := test_cases.StdinTestCaseCollection{
		{
			Pattern:          "[abcd]",
			Input:            "a",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "[abcd]",
			Input:            "efgh",
			ExpectedExitCode: 1,
		},
	}

	return testCases.Run(stageHarness)
}
