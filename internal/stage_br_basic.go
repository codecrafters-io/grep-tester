package internal

import (
	tester_utils "github.com/codecrafters-io/tester-utils"
)

func testBrBasic(stageHarness *tester_utils.StageHarness) error {
	testCases := []TestCase{
		// Base case
		{
			Pattern:          "(cat) and \\1",
			Input:            "cat and cat",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "(cat) and \\1",
			Input:            "cat and dog",
			ExpectedExitCode: 1,
		},
		// Integration with concepts from previous stages
		{
			Pattern:          "(\\w\\w\\w\\w \\d\\d\\d) is doing \\1 times",
			Input:            "grep 101 is doing grep 101 times",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "([abcd]+) is \\1, not [^xyz]+",
			Input:            "abcd is abcd, not efg",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "^(\\w+) starts and ends with \\1$",
			Input:            "this starts and ends with this",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "once a (drea+mer), alwaysz? a \\1",
			Input:            "once a dreaaamer, always a dreaaamer",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "(b..s|c..e) here and \\1 there",
			Input:            "bugs here and bugs there",
			ExpectedExitCode: 0,
		},
	}

	return RunTestCases(testCases, stageHarness)
}
