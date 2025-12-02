package grep_executable

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/codecrafters-io/grep-tester/internal/condition_reader"
	"github.com/codecrafters-io/tester-utils/executable"
	"github.com/codecrafters-io/tester-utils/linewriter"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
	ptylib "github.com/creack/pty"
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
	stdoutBuffer *bytes.Buffer
	stderrBuffer *bytes.Buffer
	cmd          *exec.Cmd
	ptyMaster    *os.File
	ptySlave     *os.File
	stdoutLogger *linewriter.LineWriter
	stderrLogger *linewriter.LineWriter
	stderrPipe   io.ReadCloser
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

func (e *GrepExecutable) RunWithStdinInTty(stdin []byte, args ...string) (executable.ExecutableResult, error) {
	// Start executable in TTY
	if err := e.startInTty(args...); err != nil {
		return executable.ExecutableResult{}, err
	}

	// Close after executable has exitted
	defer func() {
		e.ptyMaster.Close()
		e.ptySlave.Close()
	}()

	// Write input and VEOF character
	if _, err := e.writeLineToTTY(stdin); err != nil {
		return executable.ExecutableResult{}, err
	}

	if _, err := e.writeVeofCharacter(); err != nil {
		return executable.ExecutableResult{}, err
	}

	return e.wait()
}

func (e *GrepExecutable) startInTty(args ...string) error {
	cmd := exec.Command(e.executable.Path, args...)
	e.cmd = cmd

	master, slave, err := ptylib.Open()
	if err != nil {
		return fmt.Errorf("Failed to open PTY pair: %s", err)
	}

	// Disable input echoing so we only read the output of the program
	disableInputEchoing(slave)

	e.ptyMaster = master
	e.ptySlave = slave

	// Wire up command's stdin, stdout with the pseudoterminal
	cmd.Stdin = slave
	cmd.Stdout = slave

	// Wire up command's stderr to a pipe so that user's debugging lines aren't tested
	e.stderrPipe, err = cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("Failed to set up stderr pipe: %s", err)
	}

	// Start the program
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Failed to start program: ", err)
	}

	// Set up loggers
	e.stdoutLogger = linewriter.New(newLoggerWriter(e.executable.GetLoggerFunction()), 500*time.Millisecond)
	e.stderrLogger = linewriter.New(newLoggerWriter(e.executable.GetLoggerFunction()), 500*time.Millisecond)

	// Setup Logging loop
	e.stdoutBuffer = bytes.NewBuffer([]byte{})
	e.stderrBuffer = bytes.NewBuffer([]byte{})
	e.setupLoggingLoop(master, io.MultiWriter(e.stdoutLogger, e.stdoutBuffer))
	e.setupLoggingLoop(e.stderrPipe, io.MultiWriter(e.stderrLogger, e.stderrBuffer))

	return nil
}

func (e *GrepExecutable) setupLoggingLoop(source io.Reader, destination io.Writer) {
	reader := condition_reader.NewConditionReader(io.TeeReader(source, destination))
	go func() {
		// Loop until either the program exits or timeout is exceeded
		reader.ReadUntilConditionOrTimeout(func() bool { return false }, time.Duration(e.executable.TimeoutInMilliseconds))
	}()
}

func (e *GrepExecutable) writeLineToTTY(input []byte) (n int, err error) {
	input = fmt.Appendf(input, "\n")
	return e.ptyMaster.Write(input)
}

func (e *GrepExecutable) writeVeofCharacter() (n int, err error) {
	return e.ptyMaster.Write([]byte{4})
}

func (e *GrepExecutable) wait() (executable.ExecutableResult, error) {
	if e.cmd == nil {
		panic("CodeCrafters internal error: WaitForTermination called before command was run")
	}

	err := e.cmd.Wait()

	exitCode := e.cmd.ProcessState.ExitCode()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitCode == -1 {
				if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
					// If the process was terminated by a signal, extract the signal number
					if status.Signaled() {
						exitCode = 128 + int(status.Signal())
					}
				}
			}
		} else {
			// Ignore other exit errors, we'd rather send the exit code back
			return executable.ExecutableResult{}, err
		}
	}

	e.stdoutLogger.Flush()
	e.stderrLogger.Flush()

	stdout := e.stdoutBuffer.Bytes()
	stderr := e.stderrBuffer.Bytes()

	result := executable.ExecutableResult{
		Stdout:   stdout,
		Stderr:   stderr,
		ExitCode: exitCode,
	}

	return result, nil
}
