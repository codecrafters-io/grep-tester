package internal

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testQuantifierExactRepitition(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	allLetters := strings.Split("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", "")
	allNumbers := strings.Split("1234567890", "")
	allAlphaNumeric := append(append(allLetters, allNumbers...), "_")

	lettersInPattern := random.RandomElementsFromArray(allLetters, 5)
	startLetter := lettersInPattern[0]
	endLetter := lettersInPattern[1]
	repeatedLetter := lettersInPattern[2]
	groupLetter1 := lettersInPattern[3]
	groupLetter2 := lettersInPattern[4]

	chosenFruits := random.RandomElementsFromArray(FRUITS, 3)

	testCaseCollection := test_cases.StdinTestCaseCollection{
		{
			Pattern:          startLetter + fmt.Sprintf("%s{2}", repeatedLetter) + endLetter,
			Input:            startLetter + strings.Repeat(repeatedLetter, 2) + endLetter,
			ExpectedExitCode: 0,
		},
		{
			Pattern:          startLetter + `\d{3}` + endLetter,
			Input:            startLetter + strings.Join(random.RandomElementsFromArray(allNumbers, 3), "") + endLetter,
			ExpectedExitCode: 0,
		},
		{
			Pattern:          startLetter + `\d{3}` + endLetter,
			Input:            startLetter + strings.Join(random.RandomElementsFromArray(allNumbers, 2), "") + endLetter,
			ExpectedExitCode: 1,
		},
		{
			Pattern:          startLetter + `[0-9]{2}` + endLetter,
			Input:            startLetter + strings.Join(random.RandomElementsFromArray(allNumbers, 3), "") + endLetter,
			ExpectedExitCode: 1,
		},
		{
			Pattern:          startLetter + `\w{4}` + endLetter,
			Input:            startLetter + strings.Join(random.RandomElementsFromArray(allAlphaNumeric, 4), "") + endLetter,
			ExpectedExitCode: 0,
		},
		{
			Pattern: startLetter +
				"(" + strings.Join(chosenFruits, "|") + ")" +
				fmt.Sprintf("{%d}", len(chosenFruits)) +
				endLetter,
			Input:            startLetter + strings.Join(random.ShuffleArray(chosenFruits), "") + endLetter,
			ExpectedExitCode: 0,
		},
		{
			Pattern:          startLetter + fmt.Sprintf("[%s]{3}", groupLetter1+groupLetter2) + endLetter,
			Input:            startLetter + groupLetter1 + groupLetter2 + endLetter,
			ExpectedExitCode: 1,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
