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
	actualResult               executable.ExecutableResult
	matchesShouldbeHighlighted bool
	logger                     *logger.Logger
}

func (a *HighlightingAssertion) Run(result executable.ExecutableResult, logger *logger.Logger) error {
	defer func() {
		// Reset the computed value(s) after the execution is over
		a.actualResult = executable.ExecutableResult{}
		a.logger = nil
	}()

	// Presence of ANSI coloring sequence implies that matches should be highlighted
	a.matchesShouldbeHighlighted = strings.Contains(a.ExpectedOutput, "\033[01;31m")
	a.logger = logger.Clone()
	a.actualResult = result

	expectedLines := utils.ProgramOutputToLines(a.ExpectedOutput)
	actualLines := utils.ProgramOutputToLines(string(result.Stdout))

	// The dimensions for the VT will be the maximum of the expected and actual output's width and height.
	// Add 1 to the maximum length because, even in the case of empty input, we need to spawn a VT.
	maxTerminalWidth := 1 + max(maxLineLength(expectedLines), maxLineLength(actualLines))

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
	a.validateExpectedScreenState(expectedScreenState)
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

	lastSuccessfulRowIndex := outputLinesCount - 1
	if comparisonError != nil {
		lastSuccessfulRowIndex = comparisonError.RowIdx - 1
	}

	// Log success messages for each row
	actualLines := actualScreenState.GetLinesOfTextUptoCursor()
	for rowIdx := 0; rowIdx <= lastSuccessfulRowIndex; rowIdx++ {
		if a.matchesShouldbeHighlighted {
			a.logger.Successf("✓ All matches in the line %q are highlighted", actualLines[rowIdx])
		} else {
			a.logger.Successf("✓ Line %q is not highlighted", actualLines[rowIdx])
		}
	}

	if comparisonError != nil {
		// Build error using the information from context
		return a.buildError(comparisonError, expectedScreenState, actualScreenState)
	}

	return nil
}

func (a *HighlightingAssertion) buildError(
	comparisonError *comparisonMismatchError,
	expectedScreenState,
	actualScreenState *virtual_terminal.ScreenState,
) error {
	var b strings.Builder

	// Print the actual raw output, not the text content from the expected and actual screen states
	expectedLine := utils.ProgramOutputToLines(a.ExpectedOutput)[comparisonError.RowIdx]
	actualLine := utils.ProgramOutputToLines(string(a.actualResult.Stdout))[comparisonError.RowIdx]

	// Comparison error message first
	b.WriteString(
		buildComparisonErrorMessageWithCursor(
			expectedLine,
			actualLine,
			comparisonError.ColumnIdx,
		),
	)
	b.WriteString("\n")

	// Print success messages
	for _, successLog := range comparisonError.PartialSuccessLogs {
		b.WriteString(colorizeString(color.FgHiGreen, successLog))
		b.WriteString("\n")
	}

	// Print error message
	b.WriteString(colorizeString(color.FgHiRed, fmt.Sprintf("⨯ %s\n", comparisonError.Error())))
	b.WriteString("\n")

	// Print ANSI code comparison
	expectedCell := expectedScreenState.MustGetCellAtPosition(comparisonError.RowIdx, comparisonError.ColumnIdx)
	actualCell := actualScreenState.MustGetCellAtPosition(comparisonError.RowIdx, comparisonError.ColumnIdx)

	b.WriteString(
		buildAnsiCodeMismatchErrorMessage(
			expectedCell.Style.String(),
			actualCell.Style.String(),
		),
	)

	return errors.New(b.String())
}
