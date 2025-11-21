package internal

import (
	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/grep-tester/internal/utils"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testOneOrMoreQuantifier(stageHarness *test_case_harness.TestCaseHarness) error {
	utils.RelocateSystemGrep(stageHarness)

	testCaseCollection := test_cases.StdinTestCaseCollection{
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
		{
			Pattern:          `^abc_\d+_xyz$`,
			Input:            "abc_123_xyz",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          `^abc_\d+_xyz$`,
			Input:            "abc_rst_xyz",
			ExpectedExitCode: 1,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
