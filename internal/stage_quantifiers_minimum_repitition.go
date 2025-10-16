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

	logLevels := []string{"INFO", "WARN", "ERROR", "DEBUG"}
	sampleLogs := []string{"disk_full", "device_unreachable", "invalid_token"}

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
			Pattern: `^\d{2}:\d{2} LOG \w{3,} \w+$`,
			Input: fmt.Sprintf("%02d:%02d LOG %s %s",
				random.RandomInt(0, 25),
				random.RandomInt(0, 60),
				random.RandomElementFromArray(logLevels),
				random.RandomElementFromArray(sampleLogs)),
			ExpectedExitCode: 0,
		},
		{
			Pattern: `^\d{2}:\d{2} LOG \w{3,} \w+$`,
			Input: fmt.Sprintf("%02d:%02d LOG OK %s",
				random.RandomInt(0, 25),
				random.RandomInt(0, 60),
				random.RandomElementFromArray(sampleLogs)),
			ExpectedExitCode: 1,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
