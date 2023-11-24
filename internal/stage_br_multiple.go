package internal

import (
	tester_utils "github.com/codecrafters-io/tester-utils"
)

func testBrMultiple(stageHarness *tester_utils.StageHarness) error {
	testCases := []TestCase{
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
			Pattern:          "([abc]+)-([def]+) is \\1-\\2, not [^xyz]+",
			Input:            "abc-def is abc-def, not efg",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "^(\\w+) (\\w+), \\1 and \\2$",
			Input:            "apple pie, apple and pie",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "(how+dy) (he?y) there, \\1 \\2",
			Input:            "howwdy hey there, howwdy hey",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "(c.t|d.g) and (f..h|b..d), \\1 with \\2",
			Input:            "cat and fish, cat with fish",
			ExpectedExitCode: 0,
		},
	}

	return RunTestCases(testCases, stageHarness)
}
