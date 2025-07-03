package internal

import (
	"fmt"

	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testMultiLineFileSearch(stageHarness *test_case_harness.TestCaseHarness) error {
	testFiles := []TestFile{
		{Path: "fruits.txt", Content: "banana\ngrape\nblackberry\nblueberry"},
	}
	if err := CreateTestFiles(testFiles, stageHarness.Logger, stageHarness); err != nil {
		return fmt.Errorf("Failed to create test files: %v", err)
	}

	testCases := []FileSearchTestCase{
		{
			Pattern:          ".*berry",
			FilePaths:        []string{"fruits.txt"},
			ExpectedExitCode: 0,
			ExpectedOutput:   []string{"blackberry", "blueberry"},
		},
		{
			Pattern:          "carrot",
			FilePaths:        []string{"fruits.txt"},
			ExpectedExitCode: 1,
			ExpectedOutput:   []string{},
		},
		{
			Pattern:          "grape",
			FilePaths:        []string{"fruits.txt"},
			ExpectedExitCode: 0,
			ExpectedOutput:   []string{"grape"},
		},
	}

	return RunFileSearchTestCases(testCases, stageHarness)
}
