package internal

import (
	"fmt"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/grep-tester/internal/utils"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testHighlightingAutoOption(stageHarness *test_case_harness.TestCaseHarness) error {
	utils.RelocateSystemGrep(stageHarness)

	words := utils.RandomWordsWithoutSubstrings(3)
	animals := random.RandomElementsFromArray(utils.ANIMALS, 2)
	fruit := random.RandomElementFromArray(utils.FRUITS)

	testCaseCollection := test_cases.HighlightingTestCaseCollection{
		// Two always option tests
		{
			Pattern:          `\w+`,
			InputLines:       []string{fruit},
			ExpectedExitCode: 0,
			ColorMode:        utils.ColorAlways,
		},
		{
			Pattern: fmt.Sprintf("(%s|%s)", animals[0], animals[1]),
			InputLines: []string{
				animals[0] + " in the wild",
				animals[1] + "in the air",
			},
			ExpectedExitCode: 0,
			ColorMode:        utils.ColorAlways,
			RunInsideTty:     true,
		},

		// Two never option tests
		{
			Pattern:          `\d+`,
			InputLines:       []string{"123" + words[1]},
			ExpectedExitCode: 0,
			ColorMode:        utils.ColorNever,
		},
		{
			Pattern:          fruit,
			InputLines:       []string{"I love " + fruit},
			ExpectedExitCode: 0,
			ColorMode:        utils.ColorNever,
			RunInsideTty:     true,
		},

		// Two auto option tests (matching)
		{
			Pattern:          fmt.Sprintf("%s$", words[2]),
			InputLines:       []string{"prefix_" + words[2]},
			ExpectedExitCode: 0,
			ColorMode:        utils.ColorAuto,
		},
		{
			Pattern:          animals[1],
			InputLines:       []string{"The " + animals[1] + " runs fast"},
			ExpectedExitCode: 0,
			ColorMode:        utils.ColorAuto,
			RunInsideTty:     true,
		},
		// One non-matching test case as well
		{
			Pattern:          `\d{5}`,
			InputLines:       []string{"no numbers here"},
			ExpectedExitCode: 1,
			ColorMode:        utils.ColorAuto,
			RunInsideTty:     true,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
