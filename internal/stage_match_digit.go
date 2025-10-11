package internal

import (
	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testMatchDigit(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	testCaseCollection := test_cases.StdinTestCaseCollection{
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
		{
			Pattern:          `\d`,
			Input:            "abc_0_xyz",
			ExpectedExitCode: 0,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
