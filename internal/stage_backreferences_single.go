package internal

import (
	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testBackreferencesSingle(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	testCaseCollection := test_cases.StdinTestCaseCollection{
		// Base case
		{
			Pattern:          `(cat) and \1`,
			Input:            "cat and cat",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          `(cat) and \1`,
			Input:            "cat and dog",
			ExpectedExitCode: 1,
		},
		{
			Pattern:          `(\w+) and \1`,
			Input:            "cat and cat",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          `(\w+) and \1`,
			Input:            "cat and dog",
			ExpectedExitCode: 1,
		},
		{
			Pattern:          `^([act]+) is \1, not [^xyz]+$`,
			Input:            "cat is cat, not dog",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          `^([act]+) is \1, not [^xyz]+$`,
			Input:            "cat is c@t, not d0g",
			ExpectedExitCode: 1,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
