package internal

import (
	"fmt"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/grep-tester/internal/utils"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testHighlightingNeverOption(stageHarness *test_case_harness.TestCaseHarness) error {
	utils.RelocateSystemGrep(stageHarness)

	words := utils.RandomWordsWithoutSubstrings(2)
	animals := random.RandomElementsFromArray(utils.ANIMALS, 3)

	testCaseCollection := test_cases.HighlightingTestCaseCollection{
		{
			Pattern:          `\d`,
			Stdin:            "1" + animals[0],
			ExpectedExitCode: 0,
			HighlightingMode: utils.ColorNever,
		},
		{
			Pattern:          `\d`,
			Stdin:            animals[1],
			ExpectedExitCode: 1,
			HighlightingMode: utils.ColorNever,
		},
		{
			Pattern:          fmt.Sprintf("^%s", words[0]),
			Stdin:            words[0] + "_suffix",
			ExpectedExitCode: 0,
			HighlightingMode: utils.ColorNever,
		},
		{
			Pattern:          fmt.Sprintf("^%s", words[1]),
			Stdin:            "prefix_" + words[1],
			ExpectedExitCode: 1,
			HighlightingMode: utils.ColorNever,
		},
		{
			Pattern:          "do+g",
			Stdin:            "doooog",
			ExpectedExitCode: 0,
			HighlightingMode: utils.ColorNever,
		},
		{
			Pattern:          `cats and dogs`,
			Stdin:            "It's raining cats and dogs",
			ExpectedExitCode: 0,
			HighlightingMode: utils.ColorNever,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
