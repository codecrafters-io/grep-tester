package grep_executable

import (
	"os"

	"golang.org/x/sys/unix"
)

// disableInputEchoing disables input echoing on the slave terminal
func disableInputEchoing(file *os.File) error {
	fd := int(file.Fd())

	// Get current terminal attributes
	termios, err := unix.IoctlGetTermios(fd, unix.TIOCGETA)
	if err != nil {
		return err
	}

	// Disable echo flag
	termios.Lflag &^= unix.ECHO

	// Apply new attributes
	if err := unix.IoctlSetTermios(fd, unix.TIOCSETA, termios); err != nil {
		return err
	}

	return nil
}
