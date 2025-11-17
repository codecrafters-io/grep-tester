package internal

import (
	"fmt"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testPrintingMatchesSingleLine(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	words := randomWordsWithoutSubstrings(3)
	fruit := random.RandomElementFromArray(FRUITS)

	testCaseCollection := test_cases.PrintMatchingLinesTestCaseCollection{
		{
			Pattern:             `\d`,
			InputLines:          []string{"banana123"},
			ExpectedExitCode:    0,
			ExpectedOutputLines: []string{"banana123"},
		},
		{
			Pattern:          `\d`,
			InputLines:       []string{"cherry"},
			ExpectedExitCode: 1,
		},
		{
			Pattern:             fmt.Sprintf("^%s", words[0]),
			InputLines:          []string{words[0] + "_suffix"},
			ExpectedExitCode:    0,
			ExpectedOutputLines: []string{words[0] + "_suffix"},
		},
		{
			Pattern:          fmt.Sprintf("^%s", words[0]),
			InputLines:       []string{"prefix_" + words[0]},
			ExpectedExitCode: 1,
		},
		{
			Pattern:             "ca+t",
			InputLines:          []string{"cat"},
			ExpectedExitCode:    0,
			ExpectedOutputLines: []string{"cat"},
		},
		{
			Pattern:             fmt.Sprintf("[%s]", fruit),
			InputLines:          []string{string(fruit[0]) + "test"},
			ExpectedExitCode:    0,
			ExpectedOutputLines: []string{string(fruit[0]) + "test"},
		},
	}

	return testCaseCollection.Run(stageHarness)
}
