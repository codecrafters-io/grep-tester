package internal

import (
	"fmt"
	"path"
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

func RunFileSearchTestCases(testCases []FileSearchTestCase, stageHarness *test_case_harness.TestCaseHarness) error {
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
			return fmt.Errorf("expected exit code %v, got %v", testCase.ExpectedExitCode, result.ExitCode)
		}
		logger.Successf("✓ Received exit code %d.", testCase.ExpectedExitCode)

		actualOutput := strings.TrimSpace(string(result.Stdout))
		if len(testCase.ExpectedOutput) == 0 {
			if actualOutput != "" {
				return fmt.Errorf("expected empty output, got: %q", actualOutput)
			}
		} else {
			actualLines := strings.Split(actualOutput, "\n")
			if len(actualLines) != len(testCase.ExpectedOutput) {
				return fmt.Errorf("expected %d output lines, got %d\nExpected: %v\nActual: %v",
					len(testCase.ExpectedOutput), len(actualLines), testCase.ExpectedOutput, actualLines)
			}

			for i, expectedLine := range testCase.ExpectedOutput {
				if actualLines[i] != expectedLine {
					return fmt.Errorf("expected line %d to be %q, got %q", i+1, expectedLine, actualLines[i])
				}
			}
		}

		logger.Successf("✓ Received correct output.")
	}

	return nil
}
