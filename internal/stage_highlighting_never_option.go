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
			InputLines:       []string{"1" + animals[0]},
			ExpectedExitCode: 0,
			ColorMode:        utils.ColorNever,
		},
		{
			Pattern:          `\d`,
			InputLines:       []string{animals[1]},
			ExpectedExitCode: 1,
			ColorMode:        utils.ColorNever,
		},
		{
			Pattern:          fmt.Sprintf("^%s", words[0]),
			InputLines:       []string{words[0] + "_suffix"},
			ExpectedExitCode: 0,
			ColorMode:        utils.ColorNever,
		},
		{
			Pattern:          fmt.Sprintf("^%s", words[1]),
			InputLines:       []string{"prefix_" + words[1]},
			ExpectedExitCode: 1,
			ColorMode:        utils.ColorNever,
		},
		{
			Pattern:          "do+g",
			InputLines:       []string{"doooog"},
			ExpectedExitCode: 0,
			ColorMode:        utils.ColorNever,
		},
		{
			Pattern: fmt.Sprintf("(%s|%s)", animals[0], animals[1]),
			InputLines: []string{
				fmt.Sprintf("It's raining %s and %s", animals[0], animals[1]),
				fmt.Sprintf("It's raining %s and %s", animals[1], animals[0]),
			},
			ExpectedExitCode: 0,
			ColorMode:        utils.ColorNever,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
