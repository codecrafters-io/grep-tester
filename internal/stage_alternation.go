package internal

import (
	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testAlternation(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	testCaseCollection := test_cases.StdinTestCaseCollection{
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
			Pattern:          "a (cat|dog)",
			Input:            "a cog",
			ExpectedExitCode: 1,
		},
		{
			Pattern:          "^I see \\d+ (cat|dog)s?$",
			Input:            "I see 1 cat",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "^I see \\d+ (cat|dog)s?$",
			Input:            "I see 42 dogs",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "^I see \\d+ (cat|dog)s?$",
			Input:            "I see a cat",
			ExpectedExitCode: 1,
		},
		{
			Pattern:          "^I see \\d+ (cat|dog)s?$",
			Input:            "I see 2 dog3",
			ExpectedExitCode: 1,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
