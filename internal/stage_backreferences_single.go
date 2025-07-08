package internal

import (
	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testBackreferencesSingle(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	testCases := test_cases.StdinTestCaseCollection{
		// Base case
		{
			Pattern: "(cat) and \\1",
			Input:   "cat and cat",
		},
		{
			Pattern: "(cat) and \\1",
			Input:   "cat and dog",
		},
		// Integration with concepts from previous stages
		{
			Pattern: "(\\w\\w\\w\\w \\d\\d\\d) is doing \\1 times",
			Input:   "grep 101 is doing grep 101 times",
		},
		{
			Pattern: "(\\w\\w\\w \\d\\d\\d) is doing \\1 times",
			Input:   "$?! 101 is doing $?! 101 times",
		},
		{
			Pattern: "(\\w\\w\\w\\w \\d\\d\\d) is doing \\1 times",
			Input:   "grep yes is doing grep yes times",
		},
		{
			Pattern: "([abcd]+) is \\1, not [^xyz]+",
			Input:   "abcd is abcd, not efg",
		},
		{
			Pattern: "([abcd]+) is \\1, not [^xyz]+",
			Input:   "efgh is efgh, not efg",
		},
		{
			Pattern: "([abcd]+) is \\1, not [^xyz]+",
			Input:   "abcd is abcd, not xyz",
		},
		{
			Pattern: "^(\\w+) starts and ends with \\1$",
			Input:   "this starts and ends with this",
		},
		{
			Pattern: "^(this) starts and ends with \\1$",
			Input:   "that starts and ends with this",
		},
		{
			Pattern: "^(this) starts and ends with \\1$",
			Input:   "this starts and ends with this?",
		},
		{
			Pattern: "once a (drea+mer), alwaysz? a \\1",
			Input:   "once a dreaaamer, always a dreaaamer",
		},
		{
			Pattern: "once a (drea+mer), alwaysz? a \\1",
			Input:   "once a dremer, always a dreaaamer",
		},
		{
			Pattern: "once a (drea+mer), alwaysz? a \\1",
			Input:   "once a dreaaamer, alwayszzz a dreaaamer",
		},
		{
			Pattern: "(b..s|c..e) here and \\1 there",
			Input:   "bugs here and bugs there",
		},
		{
			Pattern: "(b..s|c..e) here and \\1 there",
			Input:   "bugz here and bugs there",
		},
	}

	return testCases.Run(stageHarness)
}
