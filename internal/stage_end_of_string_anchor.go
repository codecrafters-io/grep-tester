package internal

import (
	"fmt"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testEndOfStringAnchor(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	words := random.RandomWords(3)

	testCaseCollection := test_cases.StdinTestCaseCollection{
		{
			Pattern:          fmt.Sprintf("%s$", words[1]),
			Input:            words[0] + words[1],
			ExpectedExitCode: 0,
		},
		{
			Pattern:          fmt.Sprintf("%s$", words[1]),
			Input:            words[1] + words[0],
			ExpectedExitCode: 1,
		},
		{
			Pattern:          fmt.Sprintf("^%s$", words[2]),
			Input:            words[2],
			ExpectedExitCode: 0,
		},
		{
			Pattern:          fmt.Sprintf("^%s$", words[2]),
			Input:            words[1] + words[2] + words[0],
			ExpectedExitCode: 1,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
