package internal

import (
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testMatchDigit(stageHarness *test_case_harness.TestCaseHarness) error {
	testCases := []TestCase{
		{
			Pattern:          `\d`,
			Input:            "123",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          `\d`,
			Input:            "apple",
			ExpectedExitCode: 1,
		},
	}

	return RunTestCases(testCases, stageHarness)
}
