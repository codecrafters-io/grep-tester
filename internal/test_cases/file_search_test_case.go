package test_cases

import (
	"fmt"
	"path"
	"sort"
	"strings"

	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

type FileSearchTestCase struct {
	Pattern          string
	FilePaths        []string
	ExpectedExitCode int
	ExpectedOutput   []string
	Recursive        bool
}

type FileSearchTestCaseCollection []FileSearchTestCase

func (testCases FileSearchTestCaseCollection) Run(stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	executable := stageHarness.Executable

	for _, testCase := range testCases {
		args := []string{}
		if testCase.Recursive {
			args = append(args, "-r")
		}
		args = append(args, "-E", testCase.Pattern)
		args = append(args, testCase.FilePaths...)

		logger.Infof("$ ./%s %s", path.Base(executable.Path), strings.Join(args, " "))

		result, err := executable.Run(args...)
		if err != nil {
			return err
		}

		if result.ExitCode != testCase.ExpectedExitCode {
			return fmt.Errorf("Expected exit code %v, got %v", testCase.ExpectedExitCode, result.ExitCode)
		}
		logger.Successf("✓ Received exit code %d.", testCase.ExpectedExitCode)

		actualOutput := strings.TrimSpace(string(result.Stdout))
		if len(testCase.ExpectedOutput) == 0 {
			if actualOutput != "" {
				return fmt.Errorf("Expected no output, got: %q", actualOutput)
			}
		} else {
			actualOutputLines := strings.Split(actualOutput, "\n")
			if len(actualOutputLines) != len(testCase.ExpectedOutput) {
				return fmt.Errorf("Expected %d output lines, got %d\nExpected: %v\nActual: %v",
					len(testCase.ExpectedOutput), len(actualOutputLines), testCase.ExpectedOutput, actualOutputLines)
			}

			// We have no expectations on the order of the output lines
			sort.Strings(actualOutputLines)
			sort.Strings(testCase.ExpectedOutput)

			for i, expectedLine := range testCase.ExpectedOutput {
				if actualOutputLines[i] != expectedLine {
					return fmt.Errorf("Expected line %d to be %q, got %q", i+1, expectedLine, actualOutputLines[i])
				}
			}
		}

		logger.Successf("✓ Received correct output.")
	}

	return nil
}
