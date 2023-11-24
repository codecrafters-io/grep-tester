package internal

import (
	tester_utils "github.com/codecrafters-io/tester-utils"
)

func testBrNested(stageHarness *tester_utils.StageHarness) error {
	testCases := []TestCase{
		// Base case
		{
			Pattern:          "('(cat) and \\2') is the same as \\1",
			Input:            "'cat and cat' is the same as 'cat and cat'",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "('(cat) and \\2') is the same as \\1",
			Input:            "'cat and cat' is the same as 'cat and dog'",
			ExpectedExitCode: 1,
		},
		// Integration with concepts from previous stages
		{
			Pattern:          "((\\w\\w\\w\\w) (\\di\\d\\d)) is doing \\2 \\3 times, and again \\1 times",
			Input:            "grep 101 is doing grep 101 times, and again grep 101 times",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "(([abc]+)-([def]+)) is \\1, not ([^xyz]+), \\2, or \\3",
			Input:            "abc-def is abc-def, not efg, abc, or def",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "^((\\w+) (\\w+)) is made of \\2 and \\3. love \\1$",
			Input:            "apple pie is made of apple and pie. love apple pie",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "'((how+dy) (he?y) there)' is made up of '\\2' and '\\3'. \\1",
			Input:            "'howwdy hey there' is made up of 'howwdy' and 'hey'. howwdy hey there",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "((c.t|d.g) and (f..h|b..d)), \\2 with \\3, \\1",
			Input:            "cat and fish, cat with fish, cat and fish",
			ExpectedExitCode: 0,
		},
	}

	return RunTestCases(testCases, stageHarness)
}
