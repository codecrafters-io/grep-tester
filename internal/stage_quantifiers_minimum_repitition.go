package internal

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testQuantifierMinimumRepitition(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	allLetters := strings.Split("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", "")
	allNumbers := strings.Split("1234567890", "")
	allAlphaNumeric := append(allLetters, allNumbers...)

	lettersInPattern := random.RandomElementsFromArray(allLetters, 5)
	startLetter := lettersInPattern[0]
	endLetter := lettersInPattern[1]
	repeatedLetter := lettersInPattern[2]
	groupLetter1 := lettersInPattern[3]
	groupLetter2 := lettersInPattern[4]

	chosenVegetables := random.RandomElementsFromArray(VEGETABLES, 3)

	testCaseCollection := test_cases.StdinTestCaseCollection{
		{
			Pattern:          startLetter + fmt.Sprintf("%s{2,}", repeatedLetter) + endLetter,
			Input:            startLetter + strings.Repeat(repeatedLetter, 2) + endLetter,
			ExpectedExitCode: 0,
		},
		{
			Pattern:          startLetter + fmt.Sprintf("%s{2,}", repeatedLetter) + endLetter,
			Input:            startLetter + strings.Repeat(repeatedLetter, random.RandomInt(3, 5)) + endLetter,
			ExpectedExitCode: 0,
		},
		{
			Pattern: startLetter + `\d{3,}` + endLetter,
			Input: startLetter +
				strings.Join(random.RandomElementsFromArray(allNumbers, 3), "") +
				endLetter,
			ExpectedExitCode: 0,
		},
		{
			Pattern: startLetter + `\d{3,}` + endLetter,
			Input: startLetter +
				strings.Join(random.RandomElementsFromArray(allNumbers, 2), "") +
				endLetter,
			ExpectedExitCode: 1,
		},
		{
			Pattern: startLetter + `\w{2,}` + endLetter,
			Input: startLetter +
				strings.Join(random.RandomElementsFromArray(allAlphaNumeric, 3), "") +
				endLetter,
			ExpectedExitCode: 0,
		},
		{
			Pattern: startLetter +
				"(" + strings.Join(chosenVegetables, "|") + ")" +
				"{2,}" +
				endLetter,
			Input: startLetter +
				strings.Join(random.ShuffleArray(chosenVegetables), "") +
				endLetter,
			ExpectedExitCode: 0,
		},
		{
			Pattern:          startLetter + fmt.Sprintf("[%s]{3,}", groupLetter1+groupLetter2) + endLetter,
			Input:            startLetter + groupLetter1 + groupLetter2 + endLetter,
			ExpectedExitCode: 1,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
