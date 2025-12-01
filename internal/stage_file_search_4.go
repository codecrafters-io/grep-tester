package internal

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/grep-tester/internal/utils"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testRecursiveFileSearch(stageHarness *test_case_harness.TestCaseHarness) error {
	utils.RelocateSystemGrep(stageHarness)

	file_path_1 := "dir/fruits-" + utils.RandomFilePrefix() + ".txt"
	file_path_2 := "dir/subdir/vegetables-" + utils.RandomFilePrefix() + ".txt"
	file_path_3 := "dir/vegetables-" + utils.RandomFilePrefix() + ".txt"
	fruits_1 := append(random.RandomElementsFromArray(utils.FRUITS, 1), "pear")
	vegetables_1 := append(random.RandomElementsFromArray(utils.VEGETABLES, 1), "celery", "cauliflower")
	vegetables_2 := append(random.RandomElementsFromArray(utils.VEGETABLES, 2), "cucumber")

	testFiles := []utils.TestFile{
		{Path: file_path_1, Content: strings.Join(fruits_1, "\n")},
		{Path: file_path_2, Content: strings.Join(vegetables_1, "\n")},
		{Path: file_path_3, Content: strings.Join(vegetables_2, "\n")},
	}
	if err := utils.CreateTestFiles(testFiles, stageHarness); err != nil {
		panic(fmt.Sprintf("CodeCrafters Internal Error: Failed to create test files: %v", err))
	}

	testCaseCollection := test_cases.FileSearchTestCaseCollection{
		{
			Pattern:                   ".+er",
			FilePaths:                 []string{"dir/"},
			ShouldEnableRecursiveFlag: true,
			ExpectedExitCode:          0,
		},
		{
			Pattern:                   ".+ar",
			FilePaths:                 []string{"dir/"},
			ShouldEnableRecursiveFlag: true,
			ExpectedExitCode:          0,
		},
		{
			Pattern:                   fmt.Sprintf("missing_fruit_%s", fruits_1),
			FilePaths:                 []string{"dir/"},
			ShouldEnableRecursiveFlag: true,
			ExpectedExitCode:          1,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
