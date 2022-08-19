package internal

import tester_utils "github.com/codecrafters-io/tester-utils"

func testCombiningCharacterClasses(stageHarness *tester_utils.StageHarness) error {
	testCases := []TestCase{
		{
			Pattern:          "\\d apple",
			Input:            "3 apples",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "\\d apple",
			Input:            "1 orange",
			ExpectedExitCode: 1,
		},
		{
			Pattern:          "\\d\\d\\d apples",
			Input:            "124 apples",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "\\d\\d\\d apples",
			Input:            "12 apples",
			ExpectedExitCode: 1,
		},
		{
			Pattern:          "\\d \\w\\w\\ws",
			Input:            "3 dogs",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "\\d \\w\\w\\ws",
			Input:            "4 dogs",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "\\d \\w\\w\\ws",
			Input:            "1 dog",
			ExpectedExitCode: 1,
		},
	}

	return RunTestCases(testCases, stageHarness)
}
