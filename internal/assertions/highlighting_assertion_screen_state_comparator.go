package assertions

import (
	"fmt"
	"strings"

	uv "github.com/charmbracelet/ultraviolet"
	"github.com/codecrafters-io/grep-tester/virtual_terminal"
	"github.com/dustin/go-humanize/english"
)

// comparisonMismatchError represents an mismatch found during screen state comparison
type comparisonMismatchError struct {
	RowIdx             int
	ColumnIdx          int
	errorMessage       string
	PartialSuccessLogs []string
}

func (e *comparisonMismatchError) Error() string {
	return e.errorMessage
}

type screenStateComparator struct {
	partialSuccessLogs []string
}

func (c *screenStateComparator) CompareHighlighting(expected, actual *virtual_terminal.ScreenState) *comparisonMismatchError {
	cursorPosition := expected.GetCursorPosition()

	// Compare upto the row in which the cursor is present
	for rowIdx := range cursorPosition.RowIndex + 1 {

		// Compare upto the column before in which the cursor is present
		maxColumnsCount := expected.GetColumnsCount()
		if cursorPosition.RowIndex == rowIdx {
			maxColumnsCount = cursorPosition.ColumnIndex
		}

		// Compare cells
		for columnIdx := range maxColumnsCount {

			expectedCell := expected.MustGetCellAtPosition(rowIdx, columnIdx)
			actualCell := actual.MustGetCellAtPosition(rowIdx, columnIdx)

			if err := c.compareCells(expectedCell, actualCell); err != nil {
				return &comparisonMismatchError{
					RowIdx:             rowIdx,
					ColumnIdx:          columnIdx,
					errorMessage:       err.Error(),
					PartialSuccessLogs: c.partialSuccessLogs,
				}
			}
		}
	}

	return nil
}

func newScreenStateComparator() *screenStateComparator {
	return &screenStateComparator{}
}

func (c *screenStateComparator) resetPartialSuccessLogs() {
	c.partialSuccessLogs = []string{}
}

func (c *screenStateComparator) addPartialSuccessLog(successLog string) {
	c.partialSuccessLogs = append(c.partialSuccessLogs, successLog)
}

func (c *screenStateComparator) compareCells(expected, actual *uv.Cell) error {
	// Reset for each cell
	c.resetPartialSuccessLogs()

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

	return firstError
}

func (c *screenStateComparator) checkFgColor(expectedCell, actualCell *uv.Cell) error {
	if expectedCell.Style.Fg != actualCell.Style.Fg {
		return fmt.Errorf("Expected color to be %s, got %s", getFgColorName(expectedCell.Style.Fg), getFgColorName(actualCell.Style.Fg))
	}

	c.addPartialSuccessLog(fmt.Sprintf("✓ Color is %s", getFgColorName(expectedCell.Style.Fg)))
	return nil
}

func (c *screenStateComparator) checkBoldAttr(expectedCell, actualCell *uv.Cell) error {
	expectedBold := expectedCell.Style.Attrs&uv.AttrBold != 0
	actualBold := actualCell.Style.Attrs&uv.AttrBold != 0

	if expectedBold != actualBold {
		if expectedBold {
			return fmt.Errorf("Expected character to be bold (ANSI code 01), was not bold")
		} else {
			return fmt.Errorf("Expected character to not be bold, was bold (ANSI code 01)")
		}
	}

	if expectedBold {
		c.addPartialSuccessLog("✓ Character is bold")
	} else {
		c.addPartialSuccessLog("✓ Character is not bold")
	}

	return nil
}

func (c *screenStateComparator) checkInvalidStyleAndAttrs(expectedCell, actualCell *uv.Cell) error {
	// Intentionally don't add success messages for extra attributes
	expectedCellAttributesWithoutBold := expectedCell.Style.Attrs &^ uv.AttrBold
	actualCellAttributesWithoutBold := actualCell.Style.Attrs &^ uv.AttrBold

	if expectedCellAttributesWithoutBold != actualCellAttributesWithoutBold {
		attributeNames := attributesToNames(actualCellAttributesWithoutBold)
		return fmt.Errorf(
			"Found extra %s: %s",
			english.PluralWord(len(attributeNames), "attribute", "attributes"),
			strings.Join(attributeNames, ", "),
		)
	}

	if actualCell.Style.Bg != expectedCell.Style.Bg {
		return fmt.Errorf(
			"Expected background color to be %s, found %s",
			getBgColorName(expectedCell.Style.Bg),
			getBgColorName(actualCell.Style.Bg),
		)
	}

	if actualCell.Style.Underline != expectedCell.Style.Underline {
		return fmt.Errorf("Expected no underline, found underline")
	}

	return nil
}
