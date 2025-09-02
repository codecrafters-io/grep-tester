package internal

import (
	"fmt"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testStartOfStringAnchor(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	words := randomWordsNoSubstrings(2)

	testCaseCollection := test_cases.StdinTestCaseCollection{
		{
			Pattern:          fmt.Sprintf("^%s", words[0]),
			Input:            words[0] + "_" + words[1],
			ExpectedExitCode: 0,
		},
		{
			Pattern:          fmt.Sprintf("^%s", words[0]),
			Input:            words[1] + "_" + words[0],
			ExpectedExitCode: 1,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
