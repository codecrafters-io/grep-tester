package internal

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/grep-tester/internal/test_cases"
	"github.com/codecrafters-io/grep-tester/internal/utils"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

func testMultiLineFileSearch(stageHarness *test_case_harness.TestCaseHarness) error {
	utils.RelocateSystemGrep(stageHarness)

	file_name := "fruits-" + utils.RandomFilePrefix() + ".txt"
	fruits := append(random.RandomElementsFromArray(utils.FRUITS, 2), "blueberry", "strawberry")
	vegetable_1 := random.RandomElementFromArray(utils.VEGETABLES)
	testFiles := []utils.TestFile{
		{Path: file_name, Content: strings.Join(fruits, "\n")},
	}
	if err := utils.CreateTestFiles(testFiles, stageHarness); err != nil {
		panic(fmt.Sprintf("CodeCrafters Internal Error: Failed to create test files: %v", err))
	}

	testCaseCollection := test_cases.FileSearchTestCaseCollection{
		{
			Pattern:          ".+berry",
			FilePaths:        []string{file_name},
			ExpectedExitCode: 0,
		},
		{
			Pattern:          vegetable_1,
			FilePaths:        []string{file_name},
			ExpectedExitCode: 1,
		},
		{
			Pattern:          fruits[0],
			FilePaths:        []string{file_name},
			ExpectedExitCode: 0,
		},
	}

	return testCaseCollection.Run(stageHarness)
}
