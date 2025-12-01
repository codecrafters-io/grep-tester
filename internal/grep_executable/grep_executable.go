package grep_executable

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/codecrafters-io/grep-tester/internal/condition_reader"
	"github.com/codecrafters-io/tester-utils/executable"
	"github.com/codecrafters-io/tester-utils/linewriter"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
	ptylib "github.com/creack/pty"
	"golang.org/x/sys/unix"
)

type loggerWriter struct {
	loggerFunc func(string)
}

func newLoggerWriter(loggerFunc func(string)) *loggerWriter {
	return &loggerWriter{
		loggerFunc: loggerFunc,
	}
}

func (w *loggerWriter) Write(bytes []byte) (n int, err error) {
	w.loggerFunc(string(bytes[:len(bytes)-1]))
	return len(bytes), nil
}

type GrepExecutable struct {
	executable *executable.Executable

	// Used only for TTY
	stdoutBuffer     *bytes.Buffer
	cmd              *exec.Cmd
	ptyMaster        *os.File
	ptySlave         *os.File
	stdoutLineWriter *linewriter.LineWriter
	stderrPipe       io.ReadCloser
	stderrLineWriter *linewriter.LineWriter
}

func NewGrepExecutable(stageHarness *test_case_harness.TestCaseHarness) *GrepExecutable {
	return &GrepExecutable{
		executable: stageHarness.NewExecutable(),
	}
}

func (e *GrepExecutable) Path() string {
	return e.executable.Path
}

func (e *GrepExecutable) Run(args ...string) (executable.ExecutableResult, error) {
	return e.executable.Run(args...)
}

func (e *GrepExecutable) RunWithStdin(stdin []byte, args ...string) (executable.ExecutableResult, error) {
	return e.executable.RunWithStdin(stdin, args...)
}

func (e *GrepExecutable) WriteLineToTTY(input []byte) (n int, err error) {
	input = fmt.Appendf(input, "\n")
	return e.ptyMaster.Write(input)
}

func (e *GrepExecutable) WriteVeofCharacter() (n int, err error) {
	return e.ptyMaster.Write([]byte{4})
}

func (e *GrepExecutable) RunWithStdinInTty(stdin []byte, args ...string) (executable.ExecutableResult, error) {
	if err := e.startInTty(args...); err != nil {
		return executable.ExecutableResult{}, err
	}
	defer func() {
		e.ptyMaster.Close()
		e.ptySlave.Close()
	}()

	waitErrorChan := make(chan error)
	go func() {
		// Why 1 exit status still returns properly???
		waitErrorChan <- e.cmd.Wait()
	}()

	if _, err := e.WriteLineToTTY(stdin); err != nil {
		return executable.ExecutableResult{}, err
	}

	if _, err := e.WriteVeofCharacter(); err != nil {
		return executable.ExecutableResult{}, err
	}

	waitError := <-waitErrorChan
	if waitError != nil {
		return executable.ExecutableResult{}, waitError
	}

	return executable.ExecutableResult{
		ExitCode: e.cmd.ProcessState.ExitCode(),
		Stdout:   e.stdoutBuffer.Bytes(),
	}, nil
}

func (e *GrepExecutable) startInTty(args ...string) error {
	cmd := exec.Command(e.executable.Path, args...)
	e.cmd = cmd

	master, slave, err := ptylib.Open()
	if err != nil {
		return fmt.Errorf("Failed to open PTY: %s", err)
	}
	disableEcho(slave)

	e.ptyMaster = master
	e.ptySlave = slave

	// Wire up command's stdin, stdout, and stderr with pseudoterminal
	cmd.Stdin = slave
	cmd.Stdout = slave
	e.stderrPipe, err = cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("Failed to set up stderr pipe: %s", err)
	}

	// Start the program
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Failed to start program: ", err)
	}

	// Set up i/o relay
	e.stdoutLineWriter = linewriter.New(newLoggerWriter(e.executable.GetLoggerFunction()), 500*time.Millisecond)
	e.stderrLineWriter = linewriter.New(newLoggerWriter(e.executable.GetLoggerFunction()), 500*time.Millisecond)

	e.stdoutBuffer = bytes.NewBuffer([]byte{})
	stdoutReader := condition_reader.NewConditionReader(io.TeeReader(master, io.MultiWriter(e.stdoutLineWriter, e.stdoutBuffer)))
	stderrReader := condition_reader.NewConditionReader(io.TeeReader(e.stderrPipe, e.stderrLineWriter))

	go func() {
		stdoutReader.ReadUntilConditionOrTimeout(func() bool { return false }, time.Duration(e.executable.TimeoutInMilliseconds))
	}()

	go func() {
		stderrReader.ReadUntilConditionOrTimeout(func() bool { return false }, time.Duration(e.executable.TimeoutInMilliseconds))
	}()

	return nil
}

func disableEcho(file *os.File) error {
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
