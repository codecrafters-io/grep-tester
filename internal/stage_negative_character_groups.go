package internal

import (
	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testNegativeCharacterGroups(stageHarness *test_case_harness.TestCaseHarness) error {
	MoveGrepToTemp(stageHarness, stageHarness.Logger)

	testCases := []test_cases.StdinTestCase{
		{
			Pattern:          "[^xyz]",
			Input:            "apple",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "[^abc]",
			Input:            "apple",
			ExpectedExitCode: 0,
		},
		{
			Pattern:          "[^anb]",
			Input:            "banana",
			ExpectedExitCode: 1,
		},
		{
			Pattern:          "[^opq]",
			Input:            "orange",
			ExpectedExitCode: 0,
		},
	}

	return test_cases.RunStdinTestCases(testCases, stageHarness)
}
