package internal

import (
	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testBackreferencesMultiple(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	testCaseCollection := test_cases.StdinTestCaseCollection{
		// Base case
		{
			Pattern:          "(\\d+) (\\w+) squares and \\1 \\2 circles",
			Input:            "3 red squares and 3 red circles",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "(\\d+) (\\w+) squares and \\1 \\2 circles",
			Input:            "3 red squares and 4 red circles",
			ExpectedExitCode: 1,
		},
		// Integration with concepts from previous stages
		{
			Pattern:          "(\\w\\w\\w\\w) (\\d\\d\\d) is doing \\1 \\2 times",
			Input:            "grep 101 is doing grep 101 times",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "(\\w\\w\\w) (\\d\\d\\d) is doing \\1 \\2 times",
			Input:            "$?! 101 is doing $?! 101 times",
			ExpectedExitCode: 1,
		},
		{
			Pattern:          "(\\w\\w\\w\\w) (\\d\\d\\d) is doing \\1 \\2 times",
			Input:            "grep yes is doing grep yes times",
			ExpectedExitCode: 1,
		},
		{
			Pattern:          "([abc]+)-([def]+) is \\1-\\2, not [^xyz]+",
			Input:            "abc-def is abc-def, not efg",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "([abc]+)-([def]+) is \\1-\\2, not [^xyz]+",
			Input:            "efg-hij is efg-hij, not efg",
			ExpectedExitCode: 1,
		},
		{
			Pattern:          "([abc]+)-([def]+) is \\1-\\2, not [^xyz]+",
			Input:            "abc-def is abc-def, not xyz",
			ExpectedExitCode: 1,
		},
		{
			Pattern:          "^(\\w+) (\\w+), \\1 and \\2$",
			Input:            "apple pie, apple and pie",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "^(apple) (\\w+), \\1 and \\2$",
			Input:            "pineapple pie, pineapple and pie",
			ExpectedExitCode: 1,
		},
		{
			Pattern:          "^(\\w+) (pie), \\1 and \\2$",
			Input:            "apple pie, apple and pies",
			ExpectedExitCode: 1,
		},
		{
			Pattern:          "(how+dy) (he?y) there, \\1 \\2",
			Input:            "howwdy hey there, howwdy hey",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "(how+dy) (he?y) there, \\1 \\2",
			Input:            "hody hey there, howwdy hey",
			ExpectedExitCode: 1,
		},
		{
			Pattern:          "(how+dy) (he?y) there, \\1 \\2",
			Input:            "howwdy heeey there, howwdy heeey",
			ExpectedExitCode: 1,
		},
		{
			Pattern:          "(c.t|d.g) and (f..h|b..d), \\1 with \\2",
			Input:            "cat and fish, cat with fish",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "(c.t|d.g) and (f..h|b..d), \\1 with \\2",
			Input:            "bat and fish, cat with fish",
			ExpectedExitCode: 1,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
