package internal

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"

	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testMatchAlphanumeric(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	words := random.RandomWords(2)
	specialCharacters := []string{"+", "-", "÷", "×", "$", "€"}

	nonWord1 := strings.Join(random.RandomElementsFromArray(specialCharacters, 3), "")
	nonWord2 := strings.Join(random.RandomElementsFromArray(specialCharacters, 3), "")
	nonWord3 := strings.Join(random.RandomElementsFromArray(specialCharacters, 6), "")

	testCases := test_cases.StdinTestCaseCollection{
		{
			Pattern: `\w`,
			Input:   words[0],
		},
		{
			Pattern: `\w`,
			Input:   strings.ToUpper(words[1]),
		},

		{
			Pattern: `\w`,
			Input:   fmt.Sprintf("%d", random.RandomInt(100, 1000)),
		},
		{
			Pattern: `\w`,
			Input:   nonWord1 + "_" + nonWord2,
		},
		{
			Pattern: `\w`,
			Input:   nonWord3,
		},
	}

	return testCases.Run(stageHarness)
}
