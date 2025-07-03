package internal

import (
	"fmt"

	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testSingleLineFileSearch(stageHarness *test_case_harness.TestCaseHarness) error {
	testFiles := []TestFile{
		{Path: "fruits.txt", Content: "apple"},
	}
	if err := CreateTestFiles(testFiles, stageHarness.Logger, stageHarness); err != nil {
		return fmt.Errorf("Failed to create test files: %v", err)
	}

	testCases := []FileSearchTestCase{
		{
			Pattern:          "appl.*",
			FilePaths:        []string{"fruits.txt"},
			ExpectedExitCode: 0,
			ExpectedOutput:   []string{"apple"},
		},
		{
			Pattern:          "carrot",
			FilePaths:        []string{"fruits.txt"},
			ExpectedExitCode: 1,
			ExpectedOutput:   []string{},
		},
		{
			Pattern:          ".*ple",
			FilePaths:        []string{"fruits.txt"},
			ExpectedExitCode: 0,
			ExpectedOutput:   []string{"apple"},
		},
	}

	return RunFileSearchTestCases(testCases, stageHarness)
}
