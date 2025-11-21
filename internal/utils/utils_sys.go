package utils

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"

	"github.com/codecrafters-io/tester-utils/test_case_harness"
)

// RelocateSystemGrep moves the system grep binary to a temporary directory
// And registers a teardown function to restore the original system grep binary
func RelocateSystemGrep(harness *test_case_harness.TestCaseHarness) {
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
		// On CI we don't run as a root user so need to use sudo
		sudoCommand := fmt.Sprintf("sudo %s", command)
		sudoMoveCmd := exec.Command("sh", "-c", sudoCommand)
		sudoMoveCmd.Stdout = io.Discard
		sudoMoveCmd.Stderr = io.Discard
		if sudoErr := sudoMoveCmd.Run(); sudoErr != nil {
			os.RemoveAll(tmpGrepDir)
			panic(fmt.Sprintf("CodeCrafters Internal Error: mv grep to tmp directory failed: %v", err))
		}
	}

	// Register teardown function to automatically restore system grep
	harness.RegisterTeardownFunc(func() { restoreSystemGrep(tmpGrepPath, oldGrepPath) })
}

// RestoreSystemGrep moves the system grep binary back to its original location and cleans up
func restoreSystemGrep(newPath string, originalPath string) error {
	command := fmt.Sprintf("mv %s %s", newPath, originalPath)
	moveCmd := exec.Command("sh", "-c", command)
	moveCmd.Stdout = io.Discard
	moveCmd.Stderr = io.Discard
	if err := moveCmd.Run(); err != nil {
		// On CI we don't run as a root user so need to use sudo
		sudoCommand := fmt.Sprintf("sudo %s", command)
		sudoMoveCmd := exec.Command("sh", "-c", sudoCommand)
		sudoMoveCmd.Stdout = io.Discard
		sudoMoveCmd.Stderr = io.Discard
		if sudoErr := sudoMoveCmd.Run(); sudoErr != nil {
			panic(fmt.Sprintf("CodeCrafters Internal Error: mv restore for grep failed: %v", err))
		}
	}

	if err := os.RemoveAll(path.Dir(newPath)); err != nil {
		panic(fmt.Sprintf("CodeCrafters Internal Error: delete tmp grep directory failed: %s", path.Dir(newPath)))
	}

	return nil
}
