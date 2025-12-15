package internal

import (
	"fmt"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/grep-tester/internal/utils"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testHighlightingAlwaysSingleMatch(stageHarness *test_case_harness.TestCaseHarness) error {
	utils.RelocateSystemGrep(stageHarness)

	words := utils.RandomWordsWithoutSubstrings(2)
	fruits := random.RandomElementsFromArray(utils.FRUITS, 3)

	testCaseCollection := test_cases.HighlightingTestCaseCollection{
		{
			Pattern:          `\d`,
			InputLines:       []string{"only1digit"},
			ExpectedExitCode: 0,
			ColorMode:        utils.ColorAlways,
		},
		{
			Pattern:          `\d`,
			InputLines:       []string{fruits[1]},
			ExpectedExitCode: 1,
			ColorMode:        utils.ColorAlways,
		},
		{
			Pattern:          fmt.Sprintf("^%s", words[0]),
			InputLines:       []string{words[0] + "_suffix"},
			ExpectedExitCode: 0,
			ColorMode:        utils.ColorAlways,
		},
		{
			Pattern:          fmt.Sprintf("^%s", words[1]),
			InputLines:       []string{"prefix_" + words[1]},
			ExpectedExitCode: 1,
			ColorMode:        utils.ColorAlways,
		},
		{
			Pattern:          "ca+t",
			InputLines:       []string{"caaaat"},
			ExpectedExitCode: 0,
			ColorMode:        utils.ColorAlways,
		},
		{
			Pattern:          `^I see \d+ (cat|dog)s?$`,
			InputLines:       []string{"I see 42 dogs"},
			ExpectedExitCode: 0,
			ColorMode:        utils.ColorAlways,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
