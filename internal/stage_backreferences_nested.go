package internal

import (
	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testBackreferencesNested(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	testCases := test_cases.StdinTestCaseCollection{
		// Base case
		{
			Pattern: "('(cat) and \\2') is the same as \\1",
			Input:   "'cat and cat' is the same as 'cat and cat'",
		},
		{
			Pattern: "('(cat) and \\2') is the same as \\1",
			Input:   "'cat and cat' is the same as 'cat and dog'",
		},
		// Integration with concepts from previous stages
		{
			Pattern: "((\\w\\w\\w\\w) (\\d\\d\\d)) is doing \\2 \\3 times, and again \\1 times",
			Input:   "grep 101 is doing grep 101 times, and again grep 101 times",
		},
		{
			Pattern: "((\\w\\w\\w) (\\d\\d\\d)) is doing \\2 \\3 times, and again \\1 times",
			Input:   "$?! 101 is doing $?! 101 times, and again $?! 101 times",
		},
		{
			Pattern: "((\\w\\w\\w\\w) (\\d\\d\\d)) is doing \\2 \\3 times, and again \\1 times",
			Input:   "grep yes is doing grep yes times, and again grep yes times",
		},
		{
			Pattern: "(([abc]+)-([def]+)) is \\1, not ([^xyz]+), \\2, or \\3",
			Input:   "abc-def is abc-def, not efg, abc, or def",
		},
		{
			Pattern: "(([abc]+)-([def]+)) is \\1, not ([^xyz]+), \\2, or \\3",
			Input:   "efg-hij is efg-hij, not klm, efg, or hij",
		},
		{
			Pattern: "(([abc]+)-([def]+)) is \\1, not ([^xyz]+), \\2, or \\3",
			Input:   "abc-def is abc-def, not xyz, abc, or def",
		},
		{
			Pattern: "^((\\w+) (\\w+)) is made of \\2 and \\3. love \\1$",
			Input:   "apple pie is made of apple and pie. love apple pie",
		},
		{
			Pattern: "^((apple) (\\w+)) is made of \\2 and \\3. love \\1$",
			Input:   "pineapple pie is made of apple and pie. love apple pie",
		},
		{
			Pattern: "^((\\w+) (pie)) is made of \\2 and \\3. love \\1$",
			Input:   "apple pie is made of apple and pie. love apple pies",
		},
		{
			Pattern: "'((how+dy) (he?y) there)' is made up of '\\2' and '\\3'. \\1",
			Input:   "'howwdy hey there' is made up of 'howwdy' and 'hey'. howwdy hey there",
		},
		{
			Pattern: "'((how+dy) (he?y) there)' is made up of '\\2' and '\\3'. \\1",
			Input:   "'hody hey there' is made up of 'hody' and 'hey'. hody hey there",
		},
		{
			Pattern: "'((how+dy) (he?y) there)' is made up of '\\2' and '\\3'. \\1",
			Input:   "'howwdy heeey there' is made up of 'howwdy' and 'heeey'. howwdy heeey there",
		},
		{
			Pattern: "((c.t|d.g) and (f..h|b..d)), \\2 with \\3, \\1",
			Input:   "cat and fish, cat with fish, cat and fish",
		},
		{
			Pattern: "((c.t|d.g) and (f..h|b..d)), \\2 with \\3, \\1",
			Input:   "bat and fish, bat with fish, bat and fish",
		},
	}

	return testCases.Run(stageHarness)
}
