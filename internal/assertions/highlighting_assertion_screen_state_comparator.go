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

func (c *screenStateComparator) resetSuccessLogs() {
	c.successLogs = []string{}
}

func (c *screenStateComparator) addSuccessLog(successLog string) {
	c.successLogs = append(c.successLogs, successLog)
}

func (c *screenStateComparator) CompareHighlighting(expected, actual *virtual_terminal.ScreenState) *ComparisonError {
	cursorPosition := expected.GetCursorPosition()

	// Compare upto the row in which the cursor is present
	for rowIdx := range cursorPosition.RowIndex + 1 {
		c.cellRowIdx = rowIdx

		// Compare upto the column before in which the cursor is present
		columnsCount := expected.GetColumnsCount()
		if cursorPosition.RowIndex == rowIdx {
			columnsCount = cursorPosition.ColumnIndex
		}

		// Compare cells
		for columnIdx := range columnsCount {
			c.cellColumnIdx = columnIdx

			expectedCell := expected.MustGetCellAtPosition(rowIdx, columnIdx)
			actualCell := actual.MustGetCellAtPosition(rowIdx, columnIdx)

			if err := c.compareCells(expectedCell, actualCell); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *screenStateComparator) compareCells(expected, actual *uv.Cell) *ComparisonError {
	// Reset for each comparison
	c.resetSuccessLogs()

	// If a single cell is found which should be highlighted, which means highlighting is turned on for this run
	if expected.Style.Fg == ansi.Red && expected.Style.Attrs == uv.AttrBold {
		c.matchesShouldbeHighlighted = true
	}

	var firstError error

	// 1. Check foreground color
	if err := c.checkFgColor(expected, actual); err != nil {
		firstError = err
	}

	// 2. Check bold attribute
	if err := c.checkBoldAttr(expected, actual); err != nil && firstError == nil {
		firstError = err
	}

	// 3. Reject extra attributes
	if err := c.checkInvalidStyleAndAttrs(expected, actual); err != nil && firstError == nil {
		firstError = err
	}

	// Return comparison error if there was an issue
	if firstError != nil {
		return &ComparisonError{
			RowIdx:       c.cellRowIdx,
			ColumnIdx:    c.cellColumnIdx,
			ExpectedCell: expected,
			ActualCell:   actual,
			Message:      firstError,
			SuccessLogs:  c.successLogs,
		}
	}

	return nil
}

func (c *screenStateComparator) checkFgColor(expectedCell, actualCell *uv.Cell) error {
	if expectedCell.Style.Fg != actualCell.Style.Fg {
		return fmt.Errorf("Expected %s, got %s", getFgColorName(expectedCell.Style.Fg), getFgColorName(actualCell.Style.Fg))
	}

	c.addSuccessLog(fmt.Sprintf("✓ Color is %s", getFgColorName(expectedCell.Style.Fg)))
	return nil
}

func (c *screenStateComparator) checkBoldAttr(expectedCell, actualCell *uv.Cell) error {
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
		c.addSuccessLog("✓ Bold attribute is present")
	} else {
		c.addSuccessLog("✓ Bold attribute is not present")
	}

	return nil
}

func (c *screenStateComparator) checkInvalidStyleAndAttrs(expectedCell, actualCell *uv.Cell) error {
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
