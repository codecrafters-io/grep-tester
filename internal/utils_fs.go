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

// GetShortRandomDirectory creates a random directory in /tmp,
// creates the directories and returns the full path
// directory is of the form `/tmp/<random-word>`
// Cleanup is performed automatically, and as the total possible directories
// is very small, this should not be used without cleanup
func GetShortRandomDirectory(stageHarness *test_case_harness.TestCaseHarness) (string, error) {
	randomDir := path.Join("/tmp", random.RandomElementFromArray(SMALL_WORDS))
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

// TODO: Not ideal, use a `test_file` struct with path and content
// Logs need to handle the situation with longer files better
// A\nB\n\C > file is not ideal.
// writeFiles writes a list of files to the given paths with the given contents
func writeFiles(testFiles []TestFile, logger *logger.Logger) error {
	for _, testFile := range testFiles {
		logger.UpdateSecondaryPrefix("setup")
		logger.Infof("echo -n %q > %q", strings.TrimRight(testFile.Content, "\n"), testFile.Path)
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
