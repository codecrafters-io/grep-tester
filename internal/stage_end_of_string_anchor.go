package internal

import (
	"fmt"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/grep-tester/internal/utils"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testEndOfStringAnchor(stageHarness *test_case_harness.TestCaseHarness) error {
	utils.RelocateSystemGrep(stageHarness)

	words := utils.RandomWordsWithoutSubstrings(3)

	testCaseCollection := test_cases.StdinTestCaseCollection{
		{
			Pattern:          fmt.Sprintf("%s$", words[1]),
			Input:            words[0] + "_" + words[1],
			ExpectedExitCode: 0,
		},
		{
			Pattern:          fmt.Sprintf("%s$", words[1]),
			Input:            words[1] + "_" + words[0],
			ExpectedExitCode: 1,
		},
		{
			Pattern:          fmt.Sprintf("^%s$", words[2]),
			Input:            words[2],
			ExpectedExitCode: 0,
		},
		{
			Pattern:          fmt.Sprintf("^%s$", words[2]),
			Input:            words[2] + "_" + words[2],
			ExpectedExitCode: 1,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
