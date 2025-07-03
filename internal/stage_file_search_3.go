package internal

import (
	"fmt"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testMultiFileSearch(stageHarness *test_case_harness.TestCaseHarness) error {
	testFiles := []TestFile{
		{Path: "fruits.txt", Content: "banana\nblueberry"},
		{Path: "vegetables.txt", Content: "broccoli\ncarrot"},
	}
	if err := CreateTestFiles(testFiles, stageHarness.Logger, stageHarness); err != nil {
		return fmt.Errorf("Failed to create test files: %v", err)
	}

	testCases := []test_cases.FileSearchTestCase{
		{
			Pattern:          "b.*$",
			FilePaths:        []string{"fruits.txt", "vegetables.txt"},
			ExpectedExitCode: 0,
			ExpectedOutput:   []string{"fruits.txt:banana", "fruits.txt:blueberry", "vegetables.txt:broccoli"},
		},
		{
			Pattern:          "missing_fruit",
			FilePaths:        []string{"fruits.txt", "vegetables.txt"},
			ExpectedExitCode: 1,
			ExpectedOutput:   []string{},
		},
		{
			Pattern:          "carrot",
			FilePaths:        []string{"fruits.txt", "vegetables.txt"},
			ExpectedExitCode: 0,
			ExpectedOutput:   []string{"vegetables.txt:carrot"},
		},
	}

	return test_cases.RunFileSearchTestCases(testCases, stageHarness)
}
