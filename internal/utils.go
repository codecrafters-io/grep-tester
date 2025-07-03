package internal

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"

	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

// MoveGrepToTemp moves the system grep binary to a temporary directory
func MoveGrepToTemp(harness *test_case_harness.TestCaseHarness, logger *logger.Logger) {
	oldGrepPath, err := exec.LookPath("grep")
	if err != nil {
		panic(fmt.Sprintf("CodeCrafters Internal Error: grep executable not found: %v", err))
	}

	tmpGrepDir, err := os.MkdirTemp("/tmp", "grep-*")
	if err != nil {
		panic(fmt.Sprintf("CodeCrafters Internal Error: create tmp grep directory failed: %v", err))
	}
	tmpGrepPath := path.Join(tmpGrepDir, "grep")

	command := fmt.Sprintf("mv %s %s", oldGrepPath, tmpGrepPath)
	moveCmd := exec.Command("sh", "-c", command)
	moveCmd.Stdout = io.Discard
	moveCmd.Stderr = io.Discard
	if err := moveCmd.Run(); err != nil {
		os.RemoveAll(tmpGrepDir)
		panic(fmt.Sprintf("CodeCrafters Internal Error: mv grep to tmp directory failed: %v", err))
	}

	// Register teardown function to automatically restore grep
	harness.RegisterTeardownFunc(func() { restoreGrep(tmpGrepPath, oldGrepPath) })
}

// RestoreGrep moves the grep binary back to its original location and cleans up
func restoreGrep(newPath string, originalPath string) error {
	command := fmt.Sprintf("mv %s %s", newPath, originalPath)
	moveCmd := exec.Command("sh", "-c", command)
	moveCmd.Stdout = io.Discard
	moveCmd.Stderr = io.Discard
	if err := moveCmd.Run(); err != nil {
		panic(fmt.Sprintf("CodeCrafters Internal Error: mv restore for grep failed: %v", err))
	}

	if err := os.RemoveAll(path.Dir(newPath)); err != nil {
		panic(fmt.Sprintf("CodeCrafters Internal Error: delete tmp grep directory failed: %s", path.Dir(newPath)))
	}

	return nil
}
