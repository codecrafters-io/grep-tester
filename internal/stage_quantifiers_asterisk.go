package internal

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testQuantifierAsterisk(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	allLetters := strings.Split("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", "")
	allNumbers := strings.Split("1234567890", "")
	allAlphaNumerics := append(allLetters, allNumbers...)

	lettersInPattern := random.RandomElementsFromArray(allLetters, 4)
	startLetter := lettersInPattern[0]
	endLetter := lettersInPattern[1]
	repeatedLetter := lettersInPattern[2]
	letterNotInPattern := lettersInPattern[3]

	testCaseCollection := test_cases.StdinTestCaseCollection{
		{
			Pattern:          startLetter + fmt.Sprintf("%s*", repeatedLetter) + endLetter,
			Input:            startLetter + endLetter,
			ExpectedExitCode: 0,
		},
		{
			Pattern:          startLetter + fmt.Sprintf("%s*", repeatedLetter) + endLetter,
			Input:            startLetter + strings.Repeat(repeatedLetter, random.RandomInt(1, 4)) + endLetter,
			ExpectedExitCode: 0,
		},
		{
			Pattern:          startLetter + fmt.Sprintf("%s*", repeatedLetter) + endLetter,
			Input:            startLetter + strings.Repeat(repeatedLetter, random.RandomInt(1, 4)) + letterNotInPattern,
			ExpectedExitCode: 1,
		},
		{
			Pattern:          startLetter + `\d*` + endLetter,
			Input:            startLetter + strings.Join(random.RandomElementsFromArray(allNumbers, 3), "") + endLetter,
			ExpectedExitCode: 0,
		},
		{
			Pattern:          startLetter + `\d*` + endLetter,
			Input:            startLetter + fmt.Sprintf("12%s34", letterNotInPattern) + endLetter,
			ExpectedExitCode: 1,
		},
		{
			Pattern:          `,\w*,`,
			Input:            "," + strings.Join(random.RandomElementsFromArray(allAlphaNumerics, 3), "") + ",",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          startLetter + `[0-9]*` + endLetter,
			Input:            startLetter + strings.Join(random.RandomElementsFromArray(allNumbers, 3), "") + endLetter,
			ExpectedExitCode: 0,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
