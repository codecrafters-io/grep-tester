package test_cases

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/codecrafters-io/grep-tester/internal/grep"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

type FileSearchTestCase struct {
	Pattern                   string
	FilePaths                 []string
	ShouldEnableRecursiveFlag bool
}

type FileSearchTestCaseCollection []FileSearchTestCase

func (testCases FileSearchTestCaseCollection) Run(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	executable := stageHarness.Executable

	for _, testCase := range testCases {
		args := []string{}
		if testCase.ShouldEnableRecursiveFlag {
			args = append(args, "-r")
		}
		args = append(args, "-E", testCase.Pattern)
		args = append(args, testCase.FilePaths...)

		logger.Infof("$ ./%s %s", path.Base(executable.Path), strings.Join(args, " "))

		// Get expected results from internal grep implementation
		expectedResult := grep.SearchFiles(testCase.Pattern, testCase.FilePaths, grep.Options{
			ExtendedRegex: true,
			Recursive:     testCase.ShouldEnableRecursiveFlag,
		})

		// Run the actual executable
		actualResult, err := executable.Run(args...)
		if err != nil {
			return err
		}

		// Compare exit codes
		if actualResult.ExitCode != expectedResult.ExitCode {
			return fmt.Errorf("Expected exit code %v, got %v", expectedResult.ExitCode, actualResult.ExitCode)
		}
		logger.Successf("✓ Received exit code %d.", expectedResult.ExitCode)

		// Compare output
		actualOutput := strings.TrimSpace(string(actualResult.Stdout))

		if len(expectedResult.Stdout) == 0 {
			if actualOutput != "" {
				return fmt.Errorf("Expected no output, got: %q", actualOutput)
			}
		} else {
			actualOutputLines := strings.Split(actualOutput, "\n")
			expectedOutputLines := expectedResult.Stdout

			if len(actualOutputLines) != len(expectedOutputLines) {
				return fmt.Errorf("Expected %d output lines, got %d\nExpected: %v\nActual: %v",
					len(expectedOutputLines), len(actualOutputLines), expectedOutputLines, actualOutputLines)
			}

			// We have no expectations on the order of the output lines
			sort.Strings(actualOutputLines)
			sort.Strings(expectedOutputLines)

			for i, expectedLine := range expectedOutputLines {
				if actualOutputLines[i] != expectedLine {
					return fmt.Errorf("Expected line %d to be %q, got %q", i+1, expectedLine, actualOutputLines[i])
				}
			}
		}

		logger.Successf("✓ Received correct output.")
	}

	return nil
}
