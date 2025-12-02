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
			Stdin:            fruit,
			ExpectedExitCode: 0,
			HighlightingMode: utils.ColorAlways,
		},
		{
			Pattern:          animals[0],
			Stdin:            animals[0] + " in the wild",
			ExpectedExitCode: 0,
			HighlightingMode: utils.ColorAlways,
			RunInsideTty:     true,
		},

		// Two never option tests
		{
			Pattern:          `\d+`,
			Stdin:            "123" + words[1],
			ExpectedExitCode: 0,
			HighlightingMode: utils.ColorNever,
		},
		{
			Pattern:          fruit,
			Stdin:            "I love " + fruit,
			ExpectedExitCode: 0,
			HighlightingMode: utils.ColorNever,
			RunInsideTty:     true,
		},

		// Two auto option tests (matching)
		{
			Pattern:          fmt.Sprintf("%s$", words[2]),
			Stdin:            "prefix_" + words[2],
			ExpectedExitCode: 0,
			HighlightingMode: utils.ColorAuto,
		},
		{
			Pattern:          animals[1],
			Stdin:            "The " + animals[1] + " runs fast",
			ExpectedExitCode: 0,
			HighlightingMode: utils.ColorAuto,
			RunInsideTty:     true,
		},
		// One non-matching test case as well
		{
			Pattern:          `\d{5}`,
			Stdin:            "no numbers here",
			ExpectedExitCode: 1,
			HighlightingMode: utils.ColorAuto,
			RunInsideTty:     true,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
