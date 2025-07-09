package internal

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testRecursiveFileSearch(stageHarness *test_case_harness.TestCaseHarness) error {
	RelocateSystemGrep(stageHarness)

	file_path_1 := "dir/fruits-" + randomFilePrefix() + ".txt"
	file_path_2 := "dir/subdir/vegetables-" + randomFilePrefix() + ".txt"
	file_path_3 := "dir/vegetables-" + randomFilePrefix() + ".txt"
	fruits_1 := append(random.RandomElementsFromArray(FRUITS, 1), "pear")
	vegetables_1 := append(random.RandomElementsFromArray(VEGETABLES, 1), "celery")
	vegetables_2 := append(random.RandomElementsFromArray(VEGETABLES, 1), "cucumber")

	testFiles := []TestFile{
		{Path: file_path_1, Content: strings.Join(fruits_1, "\n")},
		{Path: file_path_2, Content: strings.Join(vegetables_1, "\n")},
		{Path: file_path_3, Content: strings.Join(vegetables_2, "\n")},
	}
	if err := CreateTestFiles(testFiles, stageHarness); err != nil {
		return fmt.Errorf("Failed to create test files: %v", err)
	}

	testCaseCollection := test_cases.FileSearchTestCaseCollection{
		{
			Pattern:                   ".*er",
			FilePaths:                 []string{"dir/"},
			ShouldEnableRecursiveFlag: true,
		},
		{
			Pattern:                   ".*ar",
			FilePaths:                 []string{"dir/"},
			ShouldEnableRecursiveFlag: true,
		},
		{
			Pattern:                   "missing_fruit",
			FilePaths:                 []string{"dir/"},
			ShouldEnableRecursiveFlag: true,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
