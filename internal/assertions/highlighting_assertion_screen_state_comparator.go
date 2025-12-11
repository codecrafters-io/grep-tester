package assertions

import (
	"fmt"
	"strings"

	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/x/ansi"
	"github.com/codecrafters-io/grep-tester/virtual_terminal"
	"github.com/dustin/go-humanize/english"
)

// screenStateComparator holds cell information for comparison
type screenStateComparator struct {
	cellRowIdx                 int
	cellColumnIdx              int
	successLogs                []string
	matchesShouldbeHighlighted bool
}

// ComparisonError represents an error found during screen state comparison
type ComparisonError struct {
	RowIdx       int
	ColumnIdx    int
	ExpectedCell *uv.Cell
	ActualCell   *uv.Cell
	Message      error
	SuccessLogs  []string
}

func newScreenStateComparator() *screenStateComparator {
	return &screenStateComparator{}
}

func (ctx *screenStateComparator) resetSuccessLogs() {
	ctx.successLogs = []string{}
}

func (ctx *screenStateComparator) addSuccessLog(successLog string) {
	ctx.successLogs = append(ctx.successLogs, successLog)
}

func (ctx *screenStateComparator) CompareHighlightingForNRows(expected, actual *virtual_terminal.ScreenState, rowsCount int) *ComparisonError {
	if rowsCount >= actual.GetRowsCount() {
		panic(fmt.Sprintf("Codecrafters Internal Error - Cannot compare up to %d rows in a screen with %d rows", rowsCount, actual.GetRowsCount()))
	}

	for rowIdx := range rowsCount {
		ctx.cellRowIdx = rowIdx

		for columnIdx := range expected.GetColumnsCount() {
			ctx.cellColumnIdx = columnIdx

			expectedCell := expected.MustGetCellAtPosition(rowIdx, columnIdx)
			actualCell := actual.MustGetCellAtPosition(rowIdx, columnIdx)

			if err := ctx.compareCells(expectedCell, actualCell); err != nil {
				return err
			}
		}
	}

	return nil
}

func (ctx *screenStateComparator) compareCells(expected, actual *uv.Cell) *ComparisonError {
	// Reset for each comparison
	ctx.resetSuccessLogs()

	// If a single cell is found which should be highlighted, which means highlighting is turned on for this run
	if expected.Style.Fg == ansi.Red && expected.Style.Attrs == uv.AttrBold {
		ctx.matchesShouldbeHighlighted = true
	}

	var firstError error

	// 1. Check foreground color
	if err := ctx.checkFgColor(expected, actual); err != nil {
		firstError = err
	}

	// 2. Check bold attribute
	if err := ctx.checkBoldAttr(expected, actual); err != nil && firstError == nil {
		firstError = err
	}

	// 3. Reject extra attributes
	if err := ctx.checkInvalidStyleAndAttrs(expected, actual); err != nil && firstError == nil {
		firstError = err
	}

	// Return comparison error if there was an issue
	if firstError != nil {
		return &ComparisonError{
			RowIdx:       ctx.cellRowIdx,
			ColumnIdx:    ctx.cellColumnIdx,
			ExpectedCell: expected,
			ActualCell:   actual,
			Message:      firstError,
			SuccessLogs:  ctx.successLogs,
		}
	}

	return nil
}

func (ctx *screenStateComparator) checkFgColor(expectedCell, actualCell *uv.Cell) error {
	if expectedCell.Style.Fg != actualCell.Style.Fg {
		// No color expected, use a different wording
		return fmt.Errorf("Expected %s, got %s", getFgColorName(expectedCell.Style.Fg), getFgColorName(actualCell.Style.Fg))
	}

	ctx.addSuccessLog(fmt.Sprintf("✓ Color is %s", getFgColorName(expectedCell.Style.Fg)))
	return nil
}

func (ctx *screenStateComparator) checkBoldAttr(expectedCell, actualCell *uv.Cell) error {
	expectedBold := expectedCell.Style.Attrs&uv.AttrBold != 0
	actualBold := actualCell.Style.Attrs&uv.AttrBold != 0

	if expectedBold != actualBold {
		if expectedBold {
			return fmt.Errorf("Expected character to be bold (ANSI code 01), was not bold")
		} else {
			return fmt.Errorf("Expected character to not be bold, was bold")
		}
	}

	if expectedBold {
		ctx.addSuccessLog("✓ Bold attribute is present")
	} else {
		ctx.addSuccessLog("✓ Bold attribute is not present")
	}

	return nil
}

func (ctx *screenStateComparator) checkInvalidStyleAndAttrs(expectedCell, actualCell *uv.Cell) error {
	// Intentionally don't add success messages for extra attributes

	if actualCell.Style.Attrs != expectedCell.Style.Attrs {
		attributesWithoutBold := actualCell.Style.Attrs &^ uv.AttrBold
		attributeNames := attributesToNames(attributesWithoutBold)
		return fmt.Errorf(
			"Expected no extra attributes, got %s: %s",
			english.PluralWord(len(attributeNames), "attribute", "attributes"),
			strings.Join(attributeNames, ", "),
		)
	}

	if actualCell.Style.Bg != expectedCell.Style.Bg {
		return fmt.Errorf(
			"Expected %s in background, found %s",
			getBgColorName(expectedCell.Style.Bg),
			getBgColorName(actualCell.Style.Bg),
		)
	}

	if actualCell.Style.Underline != expectedCell.Style.Underline {
		return fmt.Errorf("Expected no underline, found underline")
	}

	return nil
}
