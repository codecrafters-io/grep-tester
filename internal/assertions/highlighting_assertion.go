package assertions

import (
	"errors"
	"fmt"
	"strings"

	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/x/ansi"
	"github.com/codecrafters-io/grep-tester/internal/utils"
	"github.com/codecrafters-io/grep-tester/virtual_terminal"
	"github.com/codecrafters-io/tester-utils/executable"
	"github.com/codecrafters-io/tester-utils/logger"
)

// HighlightingAssertion is a single line assertion that asserts the output's content
// and highlighting color (bold-red: default color) against the expected sequence
// using a virtual terminal
type HighlightingAssertion struct {
	ExpectedOutput string

	// This is computed based on expected screen state
	matchesShouldbeHighlighted bool
	originalResult             executable.ExecutableResult
	logger                     *logger.Logger
	vtCellComparisonContext    *vtCellComparisonContext
}

func (a *HighlightingAssertion) Run(result executable.ExecutableResult, logger *logger.Logger) error {
	defer func() {
		// Reset the computed value(s) after the execution is over
		a.matchesShouldbeHighlighted = false
		a.originalResult = executable.ExecutableResult{}
		a.logger = nil
		a.vtCellComparisonContext = nil
	}()

	a.originalResult = result
	a.logger = logger.Clone()

	// The dimensions for the VT will be the value that is maximum among the expected and actual output's width and height
	// 1 more than the max length because even in the case of empty input, we need to spawn a vt
	maxTerminalWidth := 1 + max(
		len(a.ExpectedOutput),
		len(result.Stdout),
	)

	expectedLines := utils.ProgramOutputToLines(a.ExpectedOutput)
	actualLines := utils.ProgramOutputToLines(string(result.Stdout))

	// we still need a virtual terminal with a single row
	maxTerminalHeight := 1 + max(
		len(expectedLines),
		len(actualLines),
	)

	virtualTerminal1 := virtual_terminal.NewCustomVT(maxTerminalHeight, maxTerminalWidth)
	defer virtualTerminal1.Close()

	virtualTerminal2 := virtual_terminal.NewCustomVT(maxTerminalHeight, maxTerminalWidth)
	defer virtualTerminal2.Close()

	if _, err := virtualTerminal1.WriteWithCRLFTranslation([]byte(a.ExpectedOutput)); err != nil {
		return err
	}

	if _, err := virtualTerminal2.WriteWithCRLFTranslation(result.Stdout); err != nil {
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

func (a *HighlightingAssertion) getExpectedOutputAtLineIdx(lineIdx int) string {
	lines := utils.ProgramOutputToLines(a.ExpectedOutput)
	return lines[lineIdx]
}

func (a *HighlightingAssertion) getActualOutputAtLineIdx(lineIdx int) string {
	lines := utils.ProgramOutputToLines(string(a.originalResult.Stdout))
	return lines[lineIdx]
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
	a.vtCellComparisonContext = newEmptVtCellComparisonContext()

	expectedLines := expectedScreenState.GetLinesOfTextUptoCursor()
	actualLines := actualScreenState.GetLinesOfTextUptoCursor()

	rowsCount := len(expectedLines)

	for rowIdx, expectedRow := range expectedScreenState.GetAllRows() {
		// Update context: row index
		a.vtCellComparisonContext.updateRowIndex(rowIdx)

		// Only assert up to the row in which cursor is present
		if rowIdx >= rowsCount {
			break
		}

		actualRow := actualScreenState.GetRowAtIndex(rowIdx)
		if err := a.compareRows(expectedRow, actualRow); err != nil {
			return err
		}

		if a.matchesShouldbeHighlighted {
			a.logger.Successf("✓ All matches in the line %q are highlighted", actualLines[rowIdx])
		} else {
			a.logger.Successf("✓ Line %q is not highlighted", actualLines[rowIdx])
		}
	}

	return nil
}

func (a *HighlightingAssertion) compareRows(expected *virtual_terminal.Row, actual *virtual_terminal.Row) error {
	expectedCells := expected.GetCellsArray()
	actualCells := actual.GetCellsArray()

	for cellIdx, expectedCell := range expectedCells {
		actualCell := actualCells[cellIdx]

		// Update context, column idx, expected, and actual cells
		a.vtCellComparisonContext.updateColumnIndex(cellIdx)
		a.vtCellComparisonContext.updateExpectedCell(expectedCell)
		a.vtCellComparisonContext.updateActualCell(actualCell)

		if err := a.compareCells(); err != nil {
			return err
		}
	}

	return nil
}

func (a *HighlightingAssertion) compareCells() error {
	ctx := a.vtCellComparisonContext
	// Reset for each comparison
	ctx.resetSuccessLogs()

	// If a single cell is found which should be highlighted, which means highlighting is turned on for this run
	if ctx.expectedCell.Style.Fg == ansi.Red && ctx.expectedCell.Style.Attrs == uv.AttrBold {
		a.matchesShouldbeHighlighted = true
	}

	var firstError error

	// 1. Check foreground color
	if err := ctx.checkCellForegroundColor(); err != nil {
		firstError = err
	}

	// 2. Check bold attribute
	if err := ctx.checkCellBoldAttribute(); err != nil && firstError == nil {
		firstError = err
	}

	// 3. Reject extra attributes
	if err := ctx.checkInvalidStylingAndAttributes(); err != nil && firstError == nil {
		firstError = err
	}

	// Only build the error from the first error
	// This is because we still want to display success logs from later checks
	if firstError != nil {
		return a.buildError(firstError)
	}

	return nil
}

func (a *HighlightingAssertion) buildError(errorMessage error) error {
	ctx := a.vtCellComparisonContext
	// Print comparison with cursor
	a.logger.Plainln(
		buildComparisonErrorMessageWithCursor(
			a.getExpectedOutputAtLineIdx(ctx.cellRowIdx),
			a.getActualOutputAtLineIdx(ctx.cellRowIdx),
			ctx.cellColumnIdx,
		),
	)

	// Print success messages
	for _, successLog := range ctx.successLogs {
		a.logger.Successln(successLog)
	}

	// Print error message
	a.logger.Errorln(fmt.Sprintf("⨯ %s", errorMessage))

	// Print ANSI sequence comparison
	a.logger.Plainln(
		buildAnsiCodeMismatchComplaint(
			ctx.expectedCell.Style.String(),
			ctx.actualCell.Style.String(),
		),
	)

	return errors.New("Wrong highlighting pattern")
}
