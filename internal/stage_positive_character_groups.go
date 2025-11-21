package internal

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/grep-tester/internal/utils"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testPositiveCharacterGroups(stageHarness *test_case_harness.TestCaseHarness) error {
	utils.RelocateSystemGrep(stageHarness)

	words := utils.RandomWordsWithoutSubstrings(3)
	letterInsideWord0 := random.RandomElementFromArray(strings.Split(words[0], ""))
	lettersOutsideWord0 := pickLettersOutsideWord(words[0], 2)
	lettersOutsideWord1 := pickLettersOutsideWord(words[1], len(words[1]))

	testCaseCollection := test_cases.StdinTestCaseCollection{
		{
			Pattern:          fmt.Sprintf("[%s]", words[0]),
			Input:            letterInsideWord0,
			ExpectedExitCode: 0,
		},
		{
			Pattern:          fmt.Sprintf("[%s]", words[0]),
			Input:            letterInsideWord0 + lettersOutsideWord0,
			ExpectedExitCode: 0,
		},
		{
			Pattern:          fmt.Sprintf("[%s]", lettersOutsideWord1),
			Input:            words[1],
			ExpectedExitCode: 1,
		},
		{
			Pattern:          fmt.Sprintf("[%s]", words[2]),
			Input:            "[]",
			ExpectedExitCode: 1,
		},
	}

	return testCaseCollection.Run(stageHarness)
}

func pickLettersOutsideWord(word string, n int) string {
	letters := ""

	for x := 'a'; x <= 'z'; x++ {
		if !strings.Contains(word, string(x)) {
			letters += string(x)
		}
		if len(letters) >= n {
			break
		}
	}

	return letters
}
