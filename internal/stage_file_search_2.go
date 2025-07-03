package internal

import (
	"fmt"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testMultiLineFileSearch(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	MoveGrepToTemp(stageHarness, logger)

	testFiles := []TestFile{
		{Path: "fruits.txt", Content: "banana\ngrape\nblackberry\nblueberry"},
	}
	if err := CreateTestFiles(testFiles, logger, stageHarness); err != nil {
		return fmt.Errorf("Failed to create test files: %v", err)
	}

	testCases := []test_cases.FileSearchTestCase{
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

	return test_cases.RunFileSearchTestCases(testCases, stageHarness)
}
