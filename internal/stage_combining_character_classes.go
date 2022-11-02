package internal

import tester_utils "github.com/codecrafters-io/tester-utils"

func testCombiningCharacterClasses(stageHarness *tester_utils.StageHarness) error {
	testCases := []TestCase{
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

	return RunTestCases(testCases, stageHarness)
}
