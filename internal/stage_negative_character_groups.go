package internal

import (
	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testNegativeCharacterGroups(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	testCaseCollection := test_cases.StdinTestCaseCollection{
		{
			Pattern: "[^xyz]",
			Input:   "apple",
		},
		{
			Pattern: "[^abc]",
			Input:   "apple",
		},
		{
			Pattern: "[^anb]",
			Input:   "banana",
		},
		{
			Pattern: "[^opq]",
			Input:   "orange",
		},
	}

	return testCaseCollection.Run(stageHarness)
}
