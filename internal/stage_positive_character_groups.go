package internal

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testPositiveCharacterGroups(stageHarness *test_case_harness.TestCaseHarness) error {
	words := random.RandomWords(3)

	letterInsideWord0 := random.RandomElementFromArray(strings.Split(words[0], ""))
	letterOutsideWord1 := pickLetterOutsideWord(words[1], len(words[1]))

	testCases := []TestCase{
		{
			Pattern:          fmt.Sprintf("[%s]", words[0]),
			Input:            letterInsideWord0,
			ExpectedExitCode: 0,
		},
		{
			Pattern:          fmt.Sprintf("[%s]", words[1]),
			Input:            letterOutsideWord1,
			ExpectedExitCode: 1,
		},
		{
			Pattern:          fmt.Sprintf("[%s]", words[2]),
			Input:            "[]",
			ExpectedExitCode: 1,
		},
	}

	return RunTestCases(testCases, stageHarness)
}

func pickLetterOutsideWord(word string, n int) string {
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
