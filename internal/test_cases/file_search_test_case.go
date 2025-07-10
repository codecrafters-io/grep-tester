package test_cases

import (
	"fmt"
	"path"
	"slices"
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

		actualOutputLines := strings.Split(actualOutput, "\n")
		expectedOutputLines := strings.Split(expectedOutput, "\n")

		foundLines := []string{}
		missingLines := []string{}
		extraLines := []string{}

		for _, expectedLine := range expectedOutputLines {
			if slices.Contains(actualOutputLines, expectedLine) {
				foundLines = append(foundLines, expectedLine)
			} else {
				missingLines = append(missingLines, expectedLine)
			}
		}

		for _, actualLine := range actualOutputLines {
			if !slices.Contains(expectedOutputLines, actualLine) {
				extraLines = append(extraLines, actualLine)
			}
		}

		if len(missingLines) == 0 && len(extraLines) == 0 && len(foundLines) == len(expectedOutputLines) {
			logger.Successf("✓ Stdout contains %d expected line(s)", len(expectedOutputLines))
		} else {
			for _, line := range foundLines {
				logger.Successf("✓ Found line '%s'", line)
			}

			if len(missingLines) > 0 {
				logger.Infof("Expected %d line(s) in output, only found %d. Missing line(s):", len(expectedOutputLines), len(foundLines))
				errorMessage := ""
				for _, line := range missingLines {
					errorMessage += fmt.Sprintf("⨯ Line not found: \"%s\"\n", line)
				}
				return fmt.Errorf("%s", errorMessage)
			}

			if len(extraLines) > 0 {
				logger.Infof("Expected %d line(s) in output, found %d. Unexpected line(s):", len(expectedOutputLines), len(actualOutputLines))
				errorMessage := ""
				for _, line := range extraLines {
					errorMessage += fmt.Sprintf("⨯ Extra line found: \"%s\"\n", line)
				}
				return fmt.Errorf("%s", errorMessage)
			}
		}
	}

	return nil
}
