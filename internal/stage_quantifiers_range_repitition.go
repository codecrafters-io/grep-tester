package internal

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testQuantifierRangeRepetition(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	fruits := random.RandomElementsFromArray(FRUITS, 2)
	fruit1 := fruits[0]
	fruit2 := fruits[1]
	vegetable := random.RandomElementFromArray(VEGETABLES)
	animals := random.RandomElementsFromArray(ANIMALS, 3)
	animal1 := animals[0]
	animal2 := animals[1]
	animal3 := animals[2]

	testCaseCollection := test_cases.StdinTestCaseCollection{
		{
			Pattern:          vegetable + string(vegetable[len(vegetable)-1]) + "{2,5}",
			Input:            vegetable + strings.Repeat(string(vegetable[len(vegetable)-1]), random.RandomInt(2, 5)),
			ExpectedExitCode: 0,
		},
		{
			Pattern:          animal1 + `\d{2,4}`,
			Input:            fmt.Sprintf("%s%d", animal1, random.RandomInt(100, 999)),
			ExpectedExitCode: 0,
		},
		{
			Pattern:          fruit1 + `,\w{1,3}`,
			Input:            fmt.Sprintf("%s,_ab", fruit1),
			ExpectedExitCode: 0,
		},
		{
			Pattern:          fmt.Sprintf("(%s|%s|%s){2,3}", animal1, animal2, animal3),
			Input:            animal3 + animal1 + animal2,
			ExpectedExitCode: 0,
		},
		{
			Pattern:          fruit2 + `[aeiou]{2,4}`,
			Input:            fruit2 + "aei",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          fmt.Sprintf(`^[a-z]{2,5} v\d{1,2}_\d{1,2}_\d{1,2} (stable|beta) %s (installed|available)$`, fruit1),
			Input:            fmt.Sprintf("%s v%d_%d_%d stable %s installed", vegetable[:random.RandomInt(2, 5)], random.RandomInt(1, 9), random.RandomInt(0, 20), random.RandomInt(0, 50), fruit1),
			ExpectedExitCode: 0,
		},
		{
			Pattern:          fmt.Sprintf(`^[a-z]{2,5} v\d{1,2}_\d{1,2}_\d{1,2} (stable|beta) %s (installed|available)$`, fruit1),
			Input:            fmt.Sprintf("%s v%d_%d_ beta %s available", vegetable[:random.RandomInt(2, 5)], random.RandomInt(1, 9), random.RandomInt(0, 20), fruit1),
			ExpectedExitCode: 1,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
