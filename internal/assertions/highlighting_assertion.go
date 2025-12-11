package assertions

import (
	"errors"
	"fmt"
	"strings"

	"github.com/codecrafters-io/grep-tester/internal/utils"
	"github.com/codecrafters-io/grep-tester/virtual_terminal"
	"github.com/codecrafters-io/tester-utils/executable"
	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/fatih/color"
)

// HighlightingAssertion asserts the output's content
// and highlighting color (bold-red: default color) against the expected sequence
// using a virtual terminal
type HighlightingAssertion struct {
	ExpectedOutput string

	// These are temporary values calculated and reset in each run
	highlightingIsTurnedOn bool
	actualResult           executable.ExecutableResult
	logger                 *logger.Logger
}

func (a *HighlightingAssertion) getExpectedOutputOnRowIdx(lineIdx int) string {
	lines := utils.ProgramOutputToLines(a.ExpectedOutput)
	return lines[lineIdx]
}

func (a *HighlightingAssertion) getActualOutputOnLineIdx(lineIdx int) string {
	lines := utils.ProgramOutputToLines(string(a.actualResult.Stdout))
	return lines[lineIdx]
}

func (a *HighlightingAssertion) Run(result executable.ExecutableResult, logger *logger.Logger) error {
	defer func() {
		// Reset the computed value(s) after the execution is over
		a.highlightingIsTurnedOn = false
		a.actualResult = executable.ExecutableResult{}
		a.logger = nil
	}()

	a.actualResult = result
	a.logger = logger.Clone()

	expectedLines := utils.ProgramOutputToLines(a.ExpectedOutput)
	actualLines := utils.ProgramOutputToLines(string(result.Stdout))

	// The dimensions for the VT will be the maximum of the expected and actual output's width and height.
	// Add 1 to the maximum length because, even in the case of empty input, we need to spawn a VT.
	maxTerminalWidth := 1 + max(
		len(a.ExpectedOutput),
		len(result.Stdout),
	)

	maxTerminalHeight := 1 + max(
		len(expectedLines),
		len(actualLines),
	)

	virtualTerminal1 := virtual_terminal.NewCustomVT(maxTerminalHeight, maxTerminalWidth)
	defer virtualTerminal1.Close()

	virtualTerminal2 := virtual_terminal.NewCustomVT(maxTerminalHeight, maxTerminalWidth)
	defer virtualTerminal2.Close()

	if _, err := virtualTerminal1.WriteWithCrlfTranslation([]byte(a.ExpectedOutput)); err != nil {
		return err
	}

	if _, err := virtualTerminal2.WriteWithCrlfTranslation(result.Stdout); err != nil {
		return err
	}

	expectedScreenState := virtualTerminal1.GetScreenState()
	a.panicIfExpectedScreenStateisFlawed(expectedScreenState)
	actualScreenState := virtualTerminal2.GetScreenState()

	// Assert text contents first
	if err := a.assertTextContents(expectedScreenState, actualScreenState); err != nil {
		return err
	}

	// Assert highlighting
	return a.assertHighlighting(expectedScreenState, actualScreenState)
}

func (a *HighlightingAssertion) assertTextContents(expectedScreenState, actualScreenState *virtual_terminal.ScreenState) error {
	expectedLines := expectedScreenState.GetLinesOfTextUptoCursor()
	actualLines := actualScreenState.GetLinesOfTextUptoCursor()

	orderedLinesAssertion := OrderedLinesAssertion{
		ExpectedOutputLines: expectedLines,
	}

	// Simulate the result as if it would have been produced without highlighting it
	actualOutput := strings.Join(actualLines, "\n")

	actualResult := executable.ExecutableResult{
		Stdout: []byte(actualOutput),
	}

	return orderedLinesAssertion.Run(actualResult, a.logger)
}

func (a *HighlightingAssertion) assertHighlighting(expectedScreenState, actualScreenState *virtual_terminal.ScreenState) error {
	// Initialize comparison context
	outputLinesCount := len(expectedScreenState.GetLinesOfTextUptoCursor())

	// Hand off to comparator for comparison
	screenStateComparator := newScreenStateComparator()
	comparisonError := screenStateComparator.CompareHighlighting(expectedScreenState, actualScreenState)

	lastSuccessFulRowIndex := outputLinesCount - 1
	if comparisonError != nil {
		lastSuccessFulRowIndex = comparisonError.RowIdx - 1
	}

	// Log success messages for each row
	actualLines := actualScreenState.GetLinesOfTextUptoCursor()
	for rowIdx := 0; rowIdx <= lastSuccessFulRowIndex; rowIdx++ {
		if screenStateComparator.matchesShouldbeHighlighted {
			a.logger.Successf("✓ All matches in the line %q are highlighted", actualLines[rowIdx])
		} else {
			a.logger.Successf("✓ Line %q is not highlighted", actualLines[rowIdx])
		}
	}

	if comparisonError != nil {
		// Build error using the information from context
		return a.buildError(comparisonError)
	}

	return nil
}

func (a *HighlightingAssertion) buildError(compErr *ComparisonError) error {
	var b strings.Builder

	// Comparison error message first
	b.WriteString(
		buildComparisonErrorMessageWithCursor(
			a.getExpectedOutputOnRowIdx(compErr.RowIdx),
			a.getActualOutputOnLineIdx(compErr.RowIdx),
			compErr.ColumnIdx,
		),
	)
	b.WriteString("\n")

	// Print success messages
	for _, successLog := range compErr.SuccessLogs {
		b.WriteString(colorizeString(color.FgHiGreen, successLog))
		b.WriteString("\n")
	}

	// Print error message
	b.WriteString(colorizeString(color.FgHiRed, fmt.Sprintf("⨯ %s\n", compErr.Message)))

	// Print ANSI sequence comparison
	b.WriteString(
		buildAnsiCodeMismatchComplaint(
			compErr.ExpectedCell.Style.String(),
			compErr.ActualCell.Style.String(),
		),
	)

	return errors.New(b.String())
}
