package internal

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

type FileSearchTestCase struct {
	Pattern          string
	FilePaths        []string
	ExpectedExitCode int
	ExpectedOutput   []string
	Recursive        bool
	TestFiles        []TestFile
}

type TestFile struct {
	Path    string
	Content string
}

func RunFileSearchTestCases(testCases []FileSearchTestCase, stageHarness *test_case_harness.TestCaseHarness) error {
	logger := stageHarness.Logger
	executable := stageHarness.Executable

	// Create test files once at the beginning of the stage
	if len(testCases) > 0 {
		logger.UpdateSecondaryPrefix("setup")
		if err := createTestFiles(testCases[0].TestFiles, logger); err != nil {
			return fmt.Errorf("failed to create test files: %v", err)
		}
		logger.ResetSecondaryPrefix()
		// Ensure cleanup happens even if test fails
		defer cleanupTestFiles(testCases[0].TestFiles)
	}

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
		// Verify stdout output
		actualOutput := strings.TrimSpace(string(result.Stdout))

		if len(testCase.ExpectedOutput) == 0 {
			// Expect empty output
			if actualOutput != "" {
				return fmt.Errorf("expected empty output, got: %q", actualOutput)
			}
		} else {
			// Verify expected output lines are present
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

		logger.Successf("✓ Received exit code %d and correct output.", testCase.ExpectedExitCode)
	}

	return nil
}

func createTestFiles(files []TestFile, logger *logger.Logger) error {
	for _, file := range files {
		// Create directory if it doesn't exist
		dir := filepath.Dir(file.Path)
		if dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", dir, err)
			}
			logger.Infof("$ mkdir -p %s", dir)
		}

		// Write file content - log the shell commands used to create it
		lines := strings.Split(file.Content, "\n")
		if len(lines) == 1 {
			// Single line file
			logger.Infof("$ echo \"%s\" > %s", file.Content, file.Path)
		} else {
			// Multi-line file
			logger.Infof("$ echo \"%s\" > %s", lines[0], file.Path)
			for _, line := range lines[1:] {
				logger.Infof("$ echo \"%s\" >> %s", line, file.Path)
			}
		}

		// Actually write the file
		if err := os.WriteFile(file.Path, []byte(file.Content), 0644); err != nil {
			return fmt.Errorf("failed to create file %s: %v", file.Path, err)
		}
	}
	return nil
}

func cleanupTestFiles(files []TestFile) {
	for _, file := range files {
		os.Remove(file.Path)
	}
}
