package internal

import (
	"fmt"

	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testRecursiveFileSearch(stageHarness *test_case_harness.TestCaseHarness) error {
	testFiles := []TestFile{
		{Path: "dir/fruits.txt", Content: "pear\nstrawberry"},
		{Path: "dir/subdir/vegetables.txt", Content: "celery\ncarrot"},
		{Path: "dir/vegetables.txt", Content: "cucumber\ncorn"},
	}
	if err := CreateTestFiles(testFiles, stageHarness.Logger, stageHarness); err != nil {
		return fmt.Errorf("Failed to create test files: %v", err)
	}

	testCases := []FileSearchTestCase{
		{
			Pattern:          ".*er",
			FilePaths:        []string{"dir/"},
			ExpectedExitCode: 0,
			ExpectedOutput:   []string{"dir/fruits.txt:strawberry", "dir/subdir/vegetables.txt:celery", "dir/vegetables.txt:cucumber"},
			Recursive:        true,
		},
		{
			Pattern:          ".*ar",
			FilePaths:        []string{"dir/"},
			ExpectedExitCode: 0,
			ExpectedOutput:   []string{"dir/fruits.txt:pear", "dir/subdir/vegetables.txt:carrot"},
			Recursive:        true,
		},
		{
			Pattern:          "missing_fruit",
			FilePaths:        []string{"dir/"},
			ExpectedExitCode: 1,
			ExpectedOutput:   []string{},
			Recursive:        true,
		},
	}

	return RunFileSearchTestCases(testCases, stageHarness)
}
