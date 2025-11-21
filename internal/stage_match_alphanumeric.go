package internal

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/grep-tester/internal/utils"

	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testMatchAlphanumeric(stageHarness *test_case_harness.TestCaseHarness) error {
	utils.RelocateSystemGrep(stageHarness)

	words := utils.RandomWordsWithoutSubstrings(2)
	specialCharacters := []string{"+", "-", "÷", "×", "=", "#", "%"}

	nonWord1 := strings.Join(random.RandomElementsFromArray(specialCharacters, 3), "")
	nonWord2 := strings.Join(random.RandomElementsFromArray(specialCharacters, 3), "")
	nonWord3 := strings.Join(random.RandomElementsFromArray(specialCharacters, 6), "")

	testCaseCollection := test_cases.StdinTestCaseCollection{
		{
			Pattern:          `\w`,
			Input:            words[0],
			ExpectedExitCode: 0,
		},
		{
			Pattern:          `\w`,
			Input:            strings.ToUpper(words[1]),
			ExpectedExitCode: 0,
		},

		{
			Pattern:          `\w`,
			Input:            fmt.Sprintf("%d", random.RandomInt(100, 1000)),
			ExpectedExitCode: 0,
		},
		{
			Pattern:          `\w`,
			Input:            nonWord1 + "_" + nonWord2,
			ExpectedExitCode: 0,
		},
		{
			Pattern:          `\w`,
			Input:            nonWord3,
			ExpectedExitCode: 1,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
