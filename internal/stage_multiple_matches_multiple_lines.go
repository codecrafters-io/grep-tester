package internal

import (
	"fmt"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/grep-tester/internal/utils"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testMultipleMatchesMultipleLines(stageHarness *test_case_harness.TestCaseHarness) error {
	utils.RelocateSystemGrep(stageHarness)

	animals := random.RandomElementsFromArray(utils.ANIMALS, 2)
	fruits := random.RandomElementsFromArray(utils.FRUITS, 2)

	animal1 := animals[0]
	animal2 := animals[1]
	fruit1 := fruits[0]
	fruit2 := fruits[1]

	testCaseCollection := test_cases.PrintMatchesOnlyTestCaseCollection{
		{
			Pattern: `\d`,
			InputLines: []string{
				"a1b",
				"no digits here",
				"2c3d",
			},
			ExpectedExitCode: 0,
		},
		{
			Pattern: fmt.Sprintf(`(%s|%s)`, fruit1, fruit2),
			InputLines: []string{
				fmt.Sprintf("I like %s", fruit1),
				"nothing here",
				fmt.Sprintf("%s and %s are tasty", fruit2, fruit1),
			},
			ExpectedExitCode: 0,
		},
		{
			Pattern: `XYZ123`,
			InputLines: []string{
				"abc",
				"def",
				"ghi",
			},
			ExpectedExitCode: 1,
		},
		{
			Pattern: `cat`,
			InputLines: []string{
				"dogdogdog",
				"elephant",
				"cat-cat-dog",
			},
			ExpectedExitCode: 0,
		},
		{
			Pattern: fmt.Sprintf(`I saw \d+ (%s|%s)s?`, animal1, animal2),
			InputLines: []string{
				fmt.Sprintf("Yesterday I saw 3 %s and I saw 45 %s.", animal1, animal2),
				"Nothing interesting today.",
				fmt.Sprintf("Last week I saw 12 %ss.", animal2),
			},
			ExpectedExitCode: 0,
		},
		{
			Pattern: `cats and dogs`,
			InputLines: []string{
				"today is sunny",
				"no rains here",
				"tomorrow maybe rainy",
			},
			ExpectedExitCode: 1,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
