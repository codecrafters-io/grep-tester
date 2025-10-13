package internal

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testQuantifierRangeRepitition(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	allLetters := strings.Split("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", "")

	lettersInPattern := random.RandomElementsFromArray(allLetters, 5)
	startLetter := lettersInPattern[0]
	endLetter := lettersInPattern[1]
	repeatedLetter := lettersInPattern[2]
	groupLetter1 := lettersInPattern[3]
	groupLetter2 := lettersInPattern[4]

	chosenFruits := random.RandomElementsFromArray(FRUITS, 3)

	testCaseCollection := test_cases.StdinTestCaseCollection{
		{
			Pattern:          startLetter + fmt.Sprintf("%s{2,4}", repeatedLetter) + endLetter,
			Input:            startLetter + strings.Repeat(repeatedLetter, 2) + endLetter,
			ExpectedExitCode: 0,
		},
		{
			Pattern:          startLetter + fmt.Sprintf("%s{2,4}", repeatedLetter) + endLetter,
			Input:            startLetter + strings.Repeat(repeatedLetter, 3) + endLetter,
			ExpectedExitCode: 0,
		},
		{
			Pattern:          startLetter + fmt.Sprintf("%s{2,4}", repeatedLetter) + endLetter,
			Input:            startLetter + strings.Repeat(repeatedLetter, 4) + endLetter,
			ExpectedExitCode: 0,
		},
		{
			Pattern:          startLetter + fmt.Sprintf("%s{2,4}", repeatedLetter) + endLetter,
			Input:            startLetter + repeatedLetter + endLetter,
			ExpectedExitCode: 1,
		},
		{
			Pattern:          startLetter + `\d{2,4}` + endLetter,
			Input:            startLetter + "123" + endLetter,
			ExpectedExitCode: 0,
		},
		{
			Pattern:          startLetter + `\w{2,3}` + endLetter,
			Input:            startLetter + "ab" + endLetter,
			ExpectedExitCode: 0,
		},
		{
			Pattern: startLetter +
				"(" + strings.Join(chosenFruits, "|") + ")" +
				"{2,4}" +
				endLetter,
			Input:            startLetter + chosenFruits[2] + chosenFruits[1] + endLetter,
			ExpectedExitCode: 0,
		},
		{
			Pattern:          startLetter + fmt.Sprintf("[%s]{2,4}", groupLetter1+groupLetter2+"e") + endLetter,
			Input:            startLetter + groupLetter1 + endLetter,
			ExpectedExitCode: 1,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
