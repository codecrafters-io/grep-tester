package internal

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testQuantifierMinimumRepetition(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	fruit := random.RandomElementFromArray(FRUITS)
	vegetables := random.RandomElementsFromArray(VEGETABLES, 3)
	vegetable1 := vegetables[0]
	vegetable2 := vegetables[1]
	vegetable3 := vegetables[2]
	animal := random.RandomElementFromArray(ANIMALS)

	testCaseCollection := test_cases.StdinTestCaseCollection{
		{
			Pattern:          fruit + string(fruit[len(fruit)-1]) + "{2,}",
			Input:            fruit + strings.Repeat(string(fruit[len(fruit)-1]), 2),
			ExpectedExitCode: 0,
		},
		{
			Pattern:          fruit + string(fruit[len(fruit)-1]) + "{2,}",
			Input:            fruit + strings.Repeat(string(fruit[len(fruit)-1]), random.RandomInt(3, 5)),
			ExpectedExitCode: 0,
		},
		{
			Pattern:          animal + `\d{3,}`,
			Input:            fmt.Sprintf("%s%d", animal, random.RandomInt(100, 9999)),
			ExpectedExitCode: 0,
		},
		{
			Pattern:          vegetable1 + `,\w{2,}`,
			Input:            fmt.Sprintf("%s,_tag", vegetable1),
			ExpectedExitCode: 0,
		},
		{
			Pattern:          fmt.Sprintf("(%s|%s|%s){3,}", vegetable1, vegetable2, vegetable3),
			Input:            vegetable3 + vegetable1 + vegetable2 + vegetable3,
			ExpectedExitCode: 0,
		},
		{
			Pattern:          fmt.Sprintf(`^(success|failure) [a-z]{3,}_\d{3,} [A-Z]{2,} %s (deployed|rolled back)$`, animal),
			Input:            fmt.Sprintf("success %s_%d %s %s deployed", vegetable1, random.RandomInt(100, 999), "PROD", animal),
			ExpectedExitCode: 0,
		},
		{
			Pattern:          fmt.Sprintf(`^(success|failure) [a-z]{3,}_\d+ [A-Z]{2,} %s (deployed|rolled back)$`, animal),
			Input:            fmt.Sprintf("failure ab_%d %s %s rolled back", random.RandomInt(100, 999), "PROD", animal),
			ExpectedExitCode: 1,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
