package internal

import (
	"fmt"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/grep-tester/internal/utils"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testHighlightingAlwaysMultipleMatches(stageHarness *test_case_harness.TestCaseHarness) error {
	utils.RelocateSystemGrep(stageHarness)

	fruit := random.RandomElementFromArray(utils.FRUITS)
	animals := random.RandomElementsFromArray(utils.ANIMALS, 3)
	vegetables := random.RandomElementsFromArray(utils.VEGETABLES, 2)

	testCaseCollection := test_cases.HighlightingTestCaseCollection{
		{
			Pattern:          `\d`,
			Stdin:            "a1b2c3",
			ExpectedExitCode: 0,
			ColorMode:        utils.ColorAlways,
		},
		{
			Pattern:          fmt.Sprintf(`(%s|%s|%s)`, animals[0], animals[1], animals[2]),
			Stdin:            fmt.Sprintf("%s_%s_%s", animals[2], animals[0], animals[1]),
			ExpectedExitCode: 0,
			ColorMode:        utils.ColorAlways,
		},
		{
			Pattern:          `\d`,
			Stdin:            fruit,
			ExpectedExitCode: 1,
			ColorMode:        utils.ColorAlways,
		},
		{
			Pattern:          `\w\w`,
			Stdin:            "xxyyzz",
			ExpectedExitCode: 0,
			ColorMode:        utils.ColorAlways,
		},
		{
			Pattern:          `\\w`,
			Stdin:            "$$##@@",
			ExpectedExitCode: 1,
			ColorMode:        utils.ColorAlways,
		},
		{
			Pattern:          fmt.Sprintf(`I see \d+ (%s|%s)s?`, vegetables[0], vegetables[1]),
			Stdin:            fmt.Sprintf("I see 3 %ss. Also, I see 4 %ss.", vegetables[1], vegetables[0]),
			ExpectedExitCode: 0,
			ColorMode:        utils.ColorAlways,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
