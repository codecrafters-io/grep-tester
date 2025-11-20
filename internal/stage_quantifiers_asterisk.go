package internal

import (
	"fmt"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/grep-tester/internal/utils"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testQuantifierAsterisk(stageHarness *test_case_harness.TestCaseHarness) error {
	utils.RelocateSystemGrep(stageHarness)

	fruit := random.RandomElementFromArray(utils.FRUITS)
	vegetable := random.RandomElementFromArray(utils.VEGETABLES)
	animals := random.RandomElementsFromArray(utils.ANIMALS, 2)
	animal1 := animals[0]
	animal2 := animals[1]

	testCaseCollection := test_cases.StdinTestCaseCollection{
		{
			Pattern:          fruit[:len(fruit)-1] + string(fruit[len(fruit)-1]) + "*",
			Input:            fruit,
			ExpectedExitCode: 0,
		},
		{
			Pattern:          fruit[:len(fruit)-1] + string(fruit[len(fruit)-1]) + "*",
			Input:            fruit[:len(fruit)-1],
			ExpectedExitCode: 0,
		},
		{
			Pattern:          animal1 + `\d*` + animal2,
			Input:            animal1 + fmt.Sprintf("%d", random.RandomInt(1, 100)) + animal2,
			ExpectedExitCode: 0,
		},
		{
			Pattern:          vegetable + `,\w*`,
			Input:            fmt.Sprintf("%s,_fresh", vegetable),
			ExpectedExitCode: 0,
		},
		{
			Pattern:          fmt.Sprintf("(%s|%s)*", animal1, animal2),
			Input:            animal2 + animal1 + animal1 + animal2,
			ExpectedExitCode: 0,
		},
		{
			Pattern:          fmt.Sprintf(`^LOG [FION]* \d+ (%s|%s)$`, fruit, vegetable),
			Input:            fmt.Sprintf("LOG INFO %d %s", random.RandomInt(10, 99), fruit),
			ExpectedExitCode: 0,
		},
		{
			Pattern:          fmt.Sprintf(`^LOG [FION]* \d+ (%s|%s)$`, fruit, vegetable),
			Input:            fmt.Sprintf("LOG info %d %s", random.RandomInt(10, 99), fruit),
			ExpectedExitCode: 1,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
