package internal

import (
	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testBackreferencesSingle(stageHarness *test_case_harness.TestCaseHarness) error {
	MoveGrepToTemp(stageHarness, stageHarness.Logger)

	testCases := test_cases.StdinTestCases{
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
			Pattern:          "(\\w\\w\\w \\d\\d\\d) is doing \\1 times",
			Input:            "$?! 101 is doing $?! 101 times",
			ExpectedExitCode: 1,
		},
		{
			Pattern:          "(\\w\\w\\w\\w \\d\\d\\d) is doing \\1 times",
			Input:            "grep yes is doing grep yes times",
			ExpectedExitCode: 1,
		},
		{
			Pattern:          "([abcd]+) is \\1, not [^xyz]+",
			Input:            "abcd is abcd, not efg",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "([abcd]+) is \\1, not [^xyz]+",
			Input:            "efgh is efgh, not efg",
			ExpectedExitCode: 1,
		},
		{
			Pattern:          "([abcd]+) is \\1, not [^xyz]+",
			Input:            "abcd is abcd, not xyz",
			ExpectedExitCode: 1,
		},
		{
			Pattern:          "^(\\w+) starts and ends with \\1$",
			Input:            "this starts and ends with this",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "^(this) starts and ends with \\1$",
			Input:            "that starts and ends with this",
			ExpectedExitCode: 1,
		},
		{
			Pattern:          "^(this) starts and ends with \\1$",
			Input:            "this starts and ends with this?",
			ExpectedExitCode: 1,
		},
		{
			Pattern:          "once a (drea+mer), alwaysz? a \\1",
			Input:            "once a dreaaamer, always a dreaaamer",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "once a (drea+mer), alwaysz? a \\1",
			Input:            "once a dremer, always a dreaaamer",
			ExpectedExitCode: 1,
		},
		{
			Pattern:          "once a (drea+mer), alwaysz? a \\1",
			Input:            "once a dreaaamer, alwayszzz a dreaaamer",
			ExpectedExitCode: 1,
		},
		{
			Pattern:          "(b..s|c..e) here and \\1 there",
			Input:            "bugs here and bugs there",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "(b..s|c..e) here and \\1 there",
			Input:            "bugz here and bugs there",
			ExpectedExitCode: 1,
		},
	}

	return testCases.Run(stageHarness)
}
