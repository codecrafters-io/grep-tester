package internal

import (
	"fmt"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testMultiFileSearch(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	testFiles := []TestFile{
		{Path: "fruits.txt", Content: "banana\nblueberry"},
		{Path: "vegetables.txt", Content: "broccoli\ncarrot"},
	}
	if err := CreateTestFiles(testFiles, stageHarness); err != nil {
		return fmt.Errorf("Failed to create test files: %v", err)
	}

	testCases := test_cases.FileSearchTestCaseCollection{
		{
			Pattern:   "b.*$",
			FilePaths: []string{"fruits.txt", "vegetables.txt"},
		},
		{
			Pattern:   "missing_fruit",
			FilePaths: []string{"fruits.txt", "vegetables.txt"},
		},
		{
			Pattern:   "carrot",
			FilePaths: []string{"fruits.txt", "vegetables.txt"},
		},
	}

	return testCases.Run(stageHarness)
}
