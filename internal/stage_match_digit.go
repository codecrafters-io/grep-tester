package internal

import (
	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testMatchDigit(stageHarness *test_case_harness.TestCaseHarness) error {
	MoveGrepToTemp(stageHarness, stageHarness.Logger)

	testCases := test_cases.StdinTestCases{
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

	return testCases.Run(stageHarness)
}
