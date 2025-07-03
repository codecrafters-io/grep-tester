package internal

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/random"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

var SMALL_WORDS = []string{"foo", "bar", "baz", "qux", "quz"}

type TestFile struct {
	Path    string
	Content string
}

func CreateTestFiles(testFiles []TestFile, logger *logger.Logger, stageHarness *test_case_harness.TestCaseHarness) error {
	randomDir, err := GetShortRandomDirectory(stageHarness)
	if err != nil {
		return fmt.Errorf("Failed to create random directory: %v", err)
	}

	for i, file := range testFiles {
		testFiles[i].Path = path.Join(randomDir, file.Path)
	}

	if err := writeFiles(testFiles, logger); err != nil {
		return fmt.Errorf("Failed to write files: %v", err)
	}

	return nil
}

// GetShortRandomDirectory creates a random directory in /tmp,
// creates the directories and returns the full path
// directory is of the form `/tmp/<random-word>`
// Cleanup is performed automatically, and as the total possible directories
// is very small, this should not be used without cleanup
func GetShortRandomDirectory(stageHarness *test_case_harness.TestCaseHarness) (string, error) {
	randomDir := path.Join("/tmp", random.RandomElementFromArray(SMALL_WORDS))
	if err := os.MkdirAll(randomDir, 0755); err != nil {
		return "", fmt.Errorf("Failed to create directory %s: %v", randomDir, err)
	}
	// Automatically cleanup the directory when the test is completed
	stageHarness.RegisterTeardownFunc(func() {
		cleanupDirectories([]string{randomDir})
	})

	return randomDir, nil
}

// writeFile writes a file to the given path with the given content
func writeFile(filePath string, content string) error {
	return os.WriteFile(filePath, []byte(content), 0644)
}

func writeFiles(testFiles []TestFile, logger *logger.Logger) error {
	for _, testFile := range testFiles {
		logger.UpdateSecondaryPrefix("setup")
		lines := strings.Split(strings.TrimRight(testFile.Content, "\n"), "\n")
		logger.Infof("echo -n %q > %q", lines[0], testFile.Path)
		if len(lines) > 1 {
			for _, line := range lines[1:] {
				logger.Infof("echo -n %q >> %q", line, testFile.Path)
			}
		}
		logger.ResetSecondaryPrefix()

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
