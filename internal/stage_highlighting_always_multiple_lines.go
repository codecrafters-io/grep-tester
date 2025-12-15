package internal

import (
	"fmt"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/grep-tester/internal/utils"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testHighlightingAlwaysMultipleLines(stageHarness *test_case_harness.TestCaseHarness) error {
	utils.RelocateSystemGrep(stageHarness)

	fruit := random.RandomElementFromArray(utils.FRUITS)
	animals := random.RandomElementsFromArray(utils.ANIMALS, 3)
	vegetables := random.RandomElementsFromArray(utils.VEGETABLES, 2)
	words := utils.RandomWordsWithoutSubstrings(2)

	testCaseCollection := test_cases.HighlightingTestCaseCollection{
		{
			Pattern: `\d`,
			InputLines: []string{
				"a1b2c3",
				"no digits here",
				"456def",
			},
			ExpectedExitCode: 0,
			ColorMode:        utils.ColorAlways,
		},
		{
			Pattern: fmt.Sprintf("(%s|%s|%s)", animals[0], animals[1], animals[2]),
			InputLines: []string{
				animals[2] + "_" + animals[0] + "_" + animals[1],
				"no_animal_here",
				animals[0] + " and " + animals[1],
			},
			ExpectedExitCode: 0,
			ColorMode:        utils.ColorAlways,
		},
		{
			Pattern: `\d`,
			InputLines: []string{
				fruit,
				animals[0],
				vegetables[0],
			},
			ExpectedExitCode: 1,
			ColorMode:        utils.ColorAlways,
		},
		{
			Pattern: `\w\w`,
			InputLines: []string{
				"xxyyzz",
				"ab cd ef",
				words[0],
			},
			ExpectedExitCode: 0,
			ColorMode:        utils.ColorAlways,
		},
		{
			Pattern: `\w`,
			InputLines: []string{
				"$$##@@",
				"!@#$%^",
				"+++---",
			},
			ExpectedExitCode: 1,
			ColorMode:        utils.ColorAlways,
		},
		{
			Pattern: fmt.Sprintf(`I see \d+ (%s|%s)s?`, vegetables[0], vegetables[1]),
			InputLines: []string{
				fmt.Sprintf("I see 3 %ss. Also, I see 4 %ss.", vegetables[1], vegetables[0]),
				fmt.Sprintf("I ate 10 %s today", vegetables[0]),
			},
			ExpectedExitCode: 0,
			ColorMode:        utils.ColorAlways,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
