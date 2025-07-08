package internal

import (
	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testInit(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	testCases := test_cases.StdinTestCaseCollection{
		{
			Pattern:          "d",
			Input:            "dog",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "f",
			Input:            "dog",
			ExpectedExitCode: 1,
		},
	}

	return testCases.Run(stageHarness)
}
