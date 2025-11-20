package utils

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

var SMALL_WORDS = []string{"foo", "bar", "baz", "qux", "quz"}

type TestFile struct {
	Path    string
	Content string
}

// CreateTestFiles creates the given test files and registers a
// teardown function to cleanup the files and directories
func CreateTestFiles(testFiles []TestFile, stageHarness *test_case_harness.TestCaseHarness) error {
	for _, testFile := range testFiles {
		dir := path.Dir(testFile.Path)
		if dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("Failed to create directory %s: %v", dir, err)
			}
			stageHarness.RegisterTeardownFunc(func() {
				cleanupDirectories([]string{dir})
			})
		} else {
			stageHarness.RegisterTeardownFunc(func() {
				cleanupFiles([]string{testFile.Path})
			})
		}
	}

	if err := writeFiles(testFiles, stageHarness.Logger); err != nil {
		return fmt.Errorf("Failed to write files: %v", err)
	}

	return nil
}

// writeFile writes a file to the given path with the given content
func writeFile(filePath string, content string) error {
	return os.WriteFile(filePath, []byte(content), 0644)
}

// writeFiles writes the given test files to the filesystem
// and logs the file creation process to the logger
func writeFiles(testFiles []TestFile, logger *logger.Logger) error {
	for _, testFile := range testFiles {
		logger.UpdateLastSecondaryPrefix("setup")
		lines := strings.Split(strings.TrimRight(testFile.Content, "\n"), "\n")
		logger.Infof("echo %q > %q", lines[0], testFile.Path)
		if len(lines) > 1 {
			for _, line := range lines[1:] {
				logger.Infof("echo %q >> %q", line, testFile.Path)
			}
		}
		logger.ResetSecondaryPrefixes()

		if err := writeFile(testFile.Path, testFile.Content); err != nil {
			logger.Errorf("Error writing file %s: %v", testFile.Path, err)
			return err
		}
	}
	return nil
}

func cleanupDirectories(dirs []string) {
	for _, dir := range dirs {
		if err := os.RemoveAll(dir); err != nil {
			panic(fmt.Sprintf("CodeCrafters internal error: Failed to cleanup directories: %s", err))
		}
	}
}

func cleanupFiles(files []string) {
	for _, file := range files {
		if err := os.Remove(file); err != nil {
			panic(fmt.Sprintf("CodeCrafters internal error: Failed to cleanup files: %s", err))
		}
	}
}
