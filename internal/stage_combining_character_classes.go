package internal

import (
	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testCombiningCharacterClasses(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	testCaseCollection := test_cases.StdinTestCaseCollection{
		{
			Pattern: `\d apple`,
			Input:   "sally has 3 apples",
		},
		{
			Pattern: `\d apple`,
			Input:   "sally has 1 orange",
		},
		{
			Pattern: `\d\d\d apples`,
			Input:   "sally has 124 apples",
		},
		{
			Pattern: `\d\\d\\d apples`,
			Input:   "sally has 12 apples",
		},
		{
			Pattern: `\d \w\w\ws`,
			Input:   "sally has 3 dogs",
		},
		{
			Pattern: `\d \w\w\ws`,
			Input:   "sally has 4 dogs",
		},
		{
			Pattern: `\d \w\w\ws`,
			Input:   "sally has 1 dog",
		},
	}

	return testCaseCollection.Run(stageHarness)
}
