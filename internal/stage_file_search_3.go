package internal

import (
	"fmt"

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

	testCases := []FileSearchTestCase{
		{
			Pattern:          "b.*$",
			FilePaths:        []string{testFiles[0].Path, testFiles[1].Path},
			ExpectedExitCode: 0,
			ExpectedOutput:   []string{"fruits.txt:banana", "fruits.txt:blueberry", "vegetables.txt:broccoli"},
		},
		{
			Pattern:          "missing_fruit",
			FilePaths:        []string{testFiles[0].Path, testFiles[1].Path},
			ExpectedExitCode: 1,
			ExpectedOutput:   []string{},
		},
		{
			Pattern:          "carrot",
			FilePaths:        []string{testFiles[0].Path, testFiles[1].Path},
			ExpectedExitCode: 0,
			ExpectedOutput:   []string{"vegetables.txt:carrot"},
		},
	}

	return RunFileSearchTestCases(testCases, stageHarness)
}
