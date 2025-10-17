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

	sampleLogs := []string{"token_created", "device_registered", "session_validated"}

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
			Pattern: `^\d{4}-\d{1,2}-\d{1,2} \d{1,2}:\d{1,2} LOG (INFO|DEBUG) \w+$`,
			Input: fmt.Sprintf("%d-%d-%d %d:%d LOG %s %s",
				// Date string
				random.RandomInt(2020, 2030),
				random.RandomInt(1, 13),
				random.RandomInt(1, 31),
				random.RandomInt(0, 24),
				random.RandomInt(0, 60),
				random.RandomElementFromArray([]string{"DEBUG", "INFO"}),
				random.RandomElementFromArray(sampleLogs)),
			ExpectedExitCode: 0,
		},
		{
			Pattern: `^\d{1,2}:\d{1,2}:\d{1,2}:\d{2,3} LOG (INFO|DEBUG) \w+$`,
			Input: fmt.Sprintf("%d:%d:%d:%d WARN %s %s",
				// Date
				random.RandomInt(0, 24),
				random.RandomInt(0, 60),
				random.RandomInt(0, 60),
				// Single digit: Regex demands 2 or 3
				random.RandomInt(0, 9),
				random.RandomElementFromArray([]string{"INFO", "DEBUG"}),
				random.RandomElementFromArray(sampleLogs)),
			ExpectedExitCode: 1,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
