package internal

import (
	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testPositiveCharacterGroups(stageHarness *test_case_harness.TestCaseHarness) error {
	testCases := []test_cases.StdinTestCase{
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

	return test_cases.RunStdinTestCases(testCases, stageHarness)
}
