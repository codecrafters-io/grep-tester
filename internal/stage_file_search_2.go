package internal

import (
	"fmt"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testMultiLineFileSearch(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	testFiles := []TestFile{
		{Path: "fruits.txt", Content: "banana\ngrape\nblackberry\nblueberry"},
	}
	if err := CreateTestFiles(testFiles, stageHarness); err != nil {
		return fmt.Errorf("Failed to create test files: %v", err)
	}

	testCases := test_cases.FileSearchTestCaseCollection{
		{
			Pattern:   ".*berry",
			FilePaths: []string{"fruits.txt"},
		},
		{
			Pattern:   "carrot",
			FilePaths: []string{"fruits.txt"},
		},
		{
			Pattern:   "grape",
			FilePaths: []string{"fruits.txt"},
		},
	}

	return testCases.Run(stageHarness)
}
