package internal

import (
	"fmt"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/grep-tester/internal/utils"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testPrintingMatchesMultipleLines(stageHarness *test_case_harness.TestCaseHarness) error {
	utils.RelocateSystemGrep(stageHarness)

	words := utils.RandomWordsWithoutSubstrings(2)
	fruit := random.RandomElementFromArray(utils.FRUITS)
	vegetable := random.RandomElementFromArray(utils.VEGETABLES)
	animals := random.RandomElementsFromArray(utils.ANIMALS, 2)
	animal1 := animals[0]
	animal2 := animals[1]

	testCaseCollection := test_cases.PrintMatchingLinesTestCaseCollection{
		{
			Pattern: `\d`,
			InputLines: []string{
				"apple",
				"banana123",
				"cherry",
				"dog456",
				"elephant",
			},
			ExpectedExitCode: 0,
		},
		{
			Pattern: `\d`,
			InputLines: []string{
				"apple",
				"banana",
				"cherry",
				"dog",
			},
			ExpectedExitCode: 1,
		},
		{
			Pattern: `\w+`,
			InputLines: []string{
				words[0],
				"!@#$",
				words[1],
				"+++",
				"test123",
			},
			ExpectedExitCode: 0,
		},
		{
			Pattern: fmt.Sprintf("(%s|%s)", animal1, animal2),
			InputLines: []string{
				"france",
				animal1,
				"italy",
				animal2,
				"spain",
			},
			ExpectedExitCode: 0,
		},
		{
			Pattern: fmt.Sprintf("(%s|%s)", animal1, animal2),
			InputLines: []string{
				"New York",
				"Washington D.C.",
				"Austin",
				"Los Angeles",
			},
			ExpectedExitCode: 1,
		},
		{
			Pattern: fmt.Sprintf(`^LOG \d+ (%s|%s)$`, fruit, vegetable),
			InputLines: []string{
				"LOG 10 " + fruit,
				"INVALID", "LOG 20 " + vegetable,
				"LOG 30 invalid",
			},
			ExpectedExitCode: 0,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
