package internal

import (
	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testCombiningCharacterClasses(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	testCases := test_cases.StdinTestCases{
		{
			Pattern:          `\d apple`,
			Input:            "sally has 3 apples",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          `\d apple`,
			Input:            "sally has 1 orange",
			ExpectedExitCode: 1,
		},
		{
			Pattern:          `\d\d\d apples`,
			Input:            "sally has 124 apples",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          `\d\\d\\d apples`,
			Input:            "sally has 12 apples",
			ExpectedExitCode: 1,
		},
		{
			Pattern:          `\d \w\w\ws`,
			Input:            "sally has 3 dogs",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          `\d \w\w\ws`,
			Input:            "sally has 4 dogs",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          `\d \w\w\ws`,
			Input:            "sally has 1 dog",
			ExpectedExitCode: 1,
		},
	}

	return testCases.Run(stageHarness)
}
