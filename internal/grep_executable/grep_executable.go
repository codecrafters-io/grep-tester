package grep_executable

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/codecrafters-io/tester-utils/executable"
	"github.com/codecrafters-io/tester-utils/linewriter"
	"github.com/codecrafters-io/tester-utils/test_case_harness"
	ptylib "github.com/creack/pty"
)

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
	readDone     chan bool
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
	if err := e.startInTty(args...); err != nil {
		return executable.ExecutableResult{}, err
	}

	// Write input and VEOF character to signal end of input
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
		return fmt.Errorf("failed to open PTY pair: %s", err)
	}

	// Disable input echoing so we only read the output of the program
	if err := disableInputEchoing(slave); err != nil {
		master.Close()
		slave.Close()
		return fmt.Errorf("failed to disable input echoing: %s", err)
	}

	e.ptyMaster = master
	e.ptySlave = slave

	// Wire up command's stdin, stdout with the pseudoterminal
	cmd.Stdin = slave
	cmd.Stdout = slave

	// Wire up command's stderr to a pipe so that user's debugging lines aren't tested
	e.stderrPipe, err = cmd.StderrPipe()
	if err != nil {
		master.Close()
		slave.Close()
		return fmt.Errorf("failed to set up stderr pipe: %s", err)
	}

	// Start the program
	if err := cmd.Start(); err != nil {
		master.Close()
		slave.Close()
		return fmt.Errorf("failed to start program: %s", err)
	}

	slave.Close()

	// Initialize buffers and loggers
	e.stdoutBuffer = bytes.NewBuffer([]byte{})
	e.stderrBuffer = bytes.NewBuffer([]byte{})
	e.stdoutLogger = linewriter.New(newLoggerWriter(e.executable.GetLoggerFunction()), 500*time.Millisecond)
	e.stderrLogger = linewriter.New(newLoggerWriter(e.executable.GetLoggerFunction()), 500*time.Millisecond)

	// Setup I/O relay - parameter order matches executable.go: (source, buffer, logger)
	e.readDone = make(chan bool)
	e.setupIORelay(master, e.stdoutBuffer, e.stdoutLogger)
	e.setupIORelay(e.stderrPipe, e.stderrBuffer, e.stderrLogger)

	return nil
}

// setupIORelay sets up I/O relay similar to executable.SetupIORelay but tracks completion locally
func (e *GrepExecutable) setupIORelay(source io.Reader, buffer *bytes.Buffer, logger *linewriter.LineWriter) {
	go func() {
		combinedDestination := io.MultiWriter(buffer, logger)
		// Limit to 30KB (~250 lines at 120 chars per line)
		bytesWritten, err := io.Copy(combinedDestination, io.LimitReader(source, 30000))
		if err != nil {
			panic(err)
		}

		if bytesWritten == 30000 {
			e.executable.GetLoggerFunction()("Warning: Logs exceeded allowed limit, output might be truncated.\n")
		}

		e.readDone <- true
		io.Copy(io.Discard, source) // Drain the pipe in case any content is leftover
	}()
}

// newLoggerWriter creates a writer that logs to the executable's logger function
func newLoggerWriter(loggerFunc func(string)) io.Writer {
	return &loggerWriter{loggerFunc: loggerFunc}
}

type loggerWriter struct {
	loggerFunc func(string)
}

func (w *loggerWriter) Write(bytes []byte) (n int, err error) {
	w.loggerFunc(string(bytes[:len(bytes)-1]))
	return len(bytes), nil
}

func (e *GrepExecutable) writeLineToTTY(input []byte) (int, error) {
	inputWithNewline := fmt.Appendf(input, "\n")
	return e.ptyMaster.Write(inputWithNewline)
}

func (e *GrepExecutable) writeVeofCharacter() (int, error) {
	return e.ptyMaster.Write([]byte{4}) // VEOF character (Ctrl+D)
}

func (e *GrepExecutable) wait() (executable.ExecutableResult, error) {
	defer e.cleanup()

	if e.cmd == nil {
		panic("CodeCrafters internal error: wait called before command was started")
	}

	// Wait for I/O relay to complete (both stdout and stderr)
	<-e.readDone
	<-e.readDone

	err := e.cmd.Wait()
	exitCode := e.extractExitCode(err)

	// Flush loggers to ensure all output is captured
	e.stdoutLogger.Flush()
	e.stderrLogger.Flush()

	return executable.ExecutableResult{
		Stdout:   e.stdoutBuffer.Bytes(),
		Stderr:   e.stderrBuffer.Bytes(),
		ExitCode: exitCode,
	}, nil
}

// extractExitCode extracts the exit code from cmd.Wait() error, handling signals
func (e *GrepExecutable) extractExitCode(err error) int {
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
		}
	}

	return exitCode
}

// cleanup cleans up resources used during TTY execution
func (e *GrepExecutable) cleanup() {
	if e.ptyMaster != nil {
		e.ptyMaster.Close()
		e.ptyMaster = nil
	}
	if e.ptySlave != nil {
		e.ptySlave.Close()
		e.ptySlave = nil
	}
	if e.stderrPipe != nil {
		e.stderrPipe.Close()
		e.stderrPipe = nil
	}
	e.cmd = nil
	e.stdoutBuffer = nil
	e.stderrBuffer = nil
	e.stdoutLogger = nil
	e.stderrLogger = nil
	e.readDone = nil
}
