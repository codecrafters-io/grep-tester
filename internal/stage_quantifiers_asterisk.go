package internal

import (
	"fmt"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testQuantifierAsterisk(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	fruit := random.RandomElementFromArray(FRUITS)
	vegetable := random.RandomElementFromArray(VEGETABLES)
	animals := random.RandomElementsFromArray(ANIMALS, 2)
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
			Pattern:          fmt.Sprintf(`^(test|run)*-[A-Z]*\d* %s (passed|failed)$`, fruit),
			Input:            fmt.Sprintf("testrun-BUILD%d %s passed", random.RandomInt(1, 9), fruit),
			ExpectedExitCode: 0,
		},
		{
			Pattern:          fmt.Sprintf(`^(test|run)*-[A-Z]*\d* %s (passed|failed)$`, fruit),
			Input:            fmt.Sprintf("runtest-BUILD %d %s failed", random.RandomInt(1, 9), fruit),
			ExpectedExitCode: 1,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
