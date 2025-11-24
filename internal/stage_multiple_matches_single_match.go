package internal

import (
	"fmt"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/grep-tester/internal/utils"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testMultipleMatchesSingleMatch(stageHarness *test_case_harness.TestCaseHarness) error {
	utils.RelocateSystemGrep(stageHarness)

	words := utils.RandomWordsWithoutSubstrings(3)

	testCaseCollection := test_cases.PrintMatchesOnlyTestCaseCollection{
		{
			Pattern:          `\d`,
			InputLines:       []string{"only1digit"},
			ExpectedExitCode: 0,
		},
		{
			Pattern:          `\d`,
			InputLines:       []string{"cherry"},
			ExpectedExitCode: 1,
		},
		{
			Pattern:          fmt.Sprintf("^%s", words[0]),
			InputLines:       []string{words[0] + "_suffix"},
			ExpectedExitCode: 0,
		},
		{
			Pattern:          fmt.Sprintf("^%s", words[0]),
			InputLines:       []string{"prefix_" + words[0]},
			ExpectedExitCode: 1,
		},
		{
			Pattern:          "ca?t",
			InputLines:       []string{"cat"},
			ExpectedExitCode: 0,
		},
		{
			Pattern:          `^I see \d+ (cat|dog)s?$`,
			InputLines:       []string{"I see 42 dogs"},
			ExpectedExitCode: 0,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
