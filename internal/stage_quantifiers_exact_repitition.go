package internal

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testQuantifierExactRepetition(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	fruit := random.RandomElementFromArray(FRUITS)
	vegetable := random.RandomElementFromArray(VEGETABLES)
	animals := random.RandomElementsFromArray(ANIMALS, 2)
	animal1 := animals[0]
	animal2 := animals[1]

	testCaseCollection := test_cases.StdinTestCaseCollection{
		{
			Pattern:          fruit + string(fruit[len(fruit)-1]) + "{2}",
			Input:            fruit + strings.Repeat(string(fruit[len(fruit)-1]), 2),
			ExpectedExitCode: 0,
		},
		{
			Pattern:          animal1 + `\d{3}`,
			Input:            fmt.Sprintf("%s%d", animal1, random.RandomInt(100, 999)),
			ExpectedExitCode: 0,
		},
		{
			Pattern:          vegetable + `,\w{4}`,
			Input:            fmt.Sprintf("%s,_hi", vegetable),
			ExpectedExitCode: 1,
		},
		{
			Pattern:          fmt.Sprintf("(%s|%s){2}", animal1, animal2),
			Input:            animal2 + animal1,
			ExpectedExitCode: 0,
		},
		{
			Pattern:          fruit + `[aeiou]{3}`,
			Input:            fruit + "uoi",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          fmt.Sprintf(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2} (%s|%s)$`, fruit, vegetable),
			Input:            fmt.Sprintf("%d-%02d-%02d %02d:%02d %s", random.RandomInt(2020, 2025), random.RandomInt(1, 12), random.RandomInt(1, 28), random.RandomInt(0, 23), random.RandomInt(0, 59), fruit),
			ExpectedExitCode: 0,
		},
		{
			Pattern:          fmt.Sprintf(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2} (%s|%s)$`, fruit, vegetable),
			Input:            fmt.Sprintf("%d-%d-%d %d:%d %s", random.RandomInt(2020, 2025), random.RandomInt(1, 12), random.RandomInt(1, 9), random.RandomInt(0, 23), random.RandomInt(0, 59), fruit),
			ExpectedExitCode: 1,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
