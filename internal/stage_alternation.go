package internal

import (
	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testAlternation(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	testCases := test_cases.StdinTestCaseCollection{
		{
			Pattern:          "a (cat|dog)",
			Input:            "a cat",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "a (cat|dog)",
			Input:            "a cow",
			ExpectedExitCode: 1,
		},
		{
			Pattern:          "^I see (\\d (cat|dog|cow)s?(, | and )?)+$",
			Input:            "I see 1 cat, 2 dogs and 3 cows",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "^I see (\\d (cat|dog|cow)(, | and )?)+$",
			Input:            "I see 1 cat, 2 dogs and 3 cows",
			ExpectedExitCode: 1,
		},
	}

	return testCases.Run(stageHarness)
}
