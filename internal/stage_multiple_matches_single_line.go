package internal

import (
	"fmt"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/grep-tester/internal/utils"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testMultipleMatchesSingleLine(stageHarness *test_case_harness.TestCaseHarness) error {
	utils.RelocateSystemGrep(stageHarness)

	animals := random.RandomElementsFromArray(utils.ANIMALS, 2)
	fruits := random.RandomElementsFromArray(utils.FRUITS, 3)

	animal1 := animals[0]
	animal2 := animals[1]
	fruit1 := fruits[0]
	fruit2 := fruits[1]
	fruit3 := fruits[2]

	testCaseCollection := test_cases.PrintMatchesOnlyTestCaseCollection{
		{
			Pattern:          `\d`,
			InputLines:       []string{"a1b2c3"},
			ExpectedExitCode: 0,
		},
		{
			Pattern:          fmt.Sprintf(`(%s|%s|%s)`, fruit1, fruit2, fruit3),
			InputLines:       []string{fmt.Sprintf("%s_%s_%s", fruit3, fruit2, fruit1)},
			ExpectedExitCode: 0,
		},
		{
			Pattern:          `\d`,
			InputLines:       []string{"cherry"},
			ExpectedExitCode: 1,
		},
		{
			Pattern:          `\w\w`,
			InputLines:       []string{"xx, yy, zz"},
			ExpectedExitCode: 0,
		},
		{
			Pattern:          `\w`,
			InputLines:       []string{"##$$%"},
			ExpectedExitCode: 1,
		},
		{
			Pattern:          fmt.Sprintf(`I see \d+ (%s|%s)s?`, animal1, animal2),
			InputLines:       []string{fmt.Sprintf("I see 3 %ss. Also, I see 4 %ss.", animal1, animal2)},
			ExpectedExitCode: 0,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
