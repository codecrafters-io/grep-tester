package internal

import (
	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testOneOrMoreQuantifier(stageHarness *test_case_harness.TestCaseHarness) error {
	MoveGrepToTemp(stageHarness, stageHarness.Logger)

	testCases := []test_cases.StdinTestCase{
		{
			Pattern:          "ca+t",
			Input:            "cat",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "ca+at",
			Input:            "caaats",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "ca+t",
			Input:            "act",
			ExpectedExitCode: 1,
		},
		{
			Pattern:          "ca+t",
			Input:            "ca",
			ExpectedExitCode: 1,
		},
	}

	return test_cases.RunStdinTestCases(testCases, stageHarness)
}
