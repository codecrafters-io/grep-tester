package internal

import (
	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/grep-tester/internal/utils"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testZeroOrOneQuantifier(stageHarness *test_case_harness.TestCaseHarness) error {
	utils.RelocateSystemGrep(stageHarness)

	testCaseCollection := test_cases.StdinTestCaseCollection{
		{
			Pattern:          "ca?t",
			Input:            "cat",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "ca?t",
			Input:            "act",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "ca?a?t",
			Input:            "cat",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "ca?t",
			Input:            "caat",
			ExpectedExitCode: 1,
		},
		{
			Pattern:          "ca?t",
			Input:            "cag",
			ExpectedExitCode: 1,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
