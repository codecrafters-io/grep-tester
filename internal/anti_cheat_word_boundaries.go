package internal

import (
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testWordBoundariesAsAntiCheat(stageHarness *test_case_harness.TestCaseHarness) error {
	testCaseCollection := AntiCheatTestCaseCollection{
		{
			Pattern:          `\bcat\b`,
			Input:            "cat",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          `\bcat\b`,
			Input:            "category",
			ExpectedExitCode: 1,
		},
		{
			Pattern:          `\Bcat\B`,
			Input:            "concatenate",
			ExpectedExitCode: 0,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
