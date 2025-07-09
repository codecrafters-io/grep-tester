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
	ExpectedExitCode          int
	ShouldEnableRecursiveFlag bool
}

type FileSearchTestCaseCollection []FileSearchTestCase

func (c FileSearchTestCaseCollection) Run(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	executable := stageHarness.Executable

	for _, testCase := range c {
		args := []string{}
		if testCase.ShouldEnableRecursiveFlag {
			args = append(args, "-r")
		}
		args = append(args, "-E", testCase.Pattern)
		args = append(args, testCase.FilePaths...)

		logger.Infof("$ ./%s %s", path.Base(executable.Path), strings.Join(args, " "))

		grepArgs := []string{}
		if testCase.ShouldEnableRecursiveFlag {
			grepArgs = append(grepArgs, "-r")
		}
		grepArgs = append(grepArgs, "-E", testCase.Pattern)
		grepArgs = append(grepArgs, testCase.FilePaths...)

		expectedResult := grep.EmulateGrep(grepArgs, []byte{})
		actualResult, err := executable.Run(args...)
		if err != nil {
			return err
		}

		if testCase.ExpectedExitCode != expectedResult.ExitCode {
			panic(fmt.Sprintf("CodeCrafters Internal Error: Expected exit code %v, grep returned %v", testCase.ExpectedExitCode, expectedResult.ExitCode))
		}
		if actualResult.ExitCode != testCase.ExpectedExitCode {
			return fmt.Errorf("Expected exit code %v, got %v", testCase.ExpectedExitCode, actualResult.ExitCode)
		}
		logger.Successf("✓ Received exit code %d.", actualResult.ExitCode)

		actualOutput := strings.TrimSpace(string(actualResult.Stdout))
		expectedOutput := strings.TrimSpace(string(expectedResult.Stdout))

		if expectedOutput == "" {
			if actualOutput != "" {
				return fmt.Errorf("Expected no output, got: %v", actualOutput)
			}
		} else {
			actualOutputLines := strings.Split(actualOutput, "\n")
			expectedOutputLines := strings.Split(expectedOutput, "\n")

			if len(actualOutputLines) != len(expectedOutputLines) {
				return fmt.Errorf("Expected %d output lines, got %d\nExpected: [%s]\nActual: [%s]",
					len(expectedOutputLines), len(actualOutputLines), strings.Join(expectedOutputLines, ", "), strings.Join(actualOutputLines, ", "))
			}

			// We have no expectations on the order of the output lines
			sort.Strings(actualOutputLines)
			sort.Strings(expectedOutputLines)

			for i, expectedLine := range expectedOutputLines {
				if actualOutputLines[i] != expectedLine {
					return fmt.Errorf("Expected: [%s] (in any order), got [%s]", strings.Join(expectedOutputLines, ", "), strings.Join(actualOutputLines, ", "))
				}
			}
		}

		logger.Successf("✓ Received correct output.")
	}

	return nil
}
