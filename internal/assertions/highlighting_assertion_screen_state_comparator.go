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

	// We've already asserted the text content previously, so if the contents do not match, panic
	if expected.Content != actual.Content {
		panic(
			fmt.Sprintf(
				"Codecrafters Internal Error - compareCells expected cell content %q got %q",
				expected.Content,
				actual.Content,
			),
		)
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

	return firstError
}

func (c *screenStateComparator) checkFgColor(expectedCell, actualCell *uv.Cell) error {
	// Skip foreground color check on white-space cells
	if isWhiteSpaceCell(actualCell) {
		return nil
	}

	if expectedCell.Style.Fg != actualCell.Style.Fg {
		return fmt.Errorf("Expected color to be %s, got %s", getFgColorName(expectedCell.Style.Fg), getFgColorName(actualCell.Style.Fg))
	}

	c.addPartialSuccessLog(fmt.Sprintf("✓ Color is %s", getFgColorName(expectedCell.Style.Fg)))
	return nil
}

func (c *screenStateComparator) checkBoldAttr(expectedCell, actualCell *uv.Cell) error {
	// Skip bold attribute check on white-space cells
	if isWhiteSpaceCell(actualCell) {
		return nil
	}

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
	expectedVisibleAttrs := getVisibleAttributesExceptBold(expectedCell)
	actualVisibleAttrs := getVisibleAttributesExceptBold(actualCell)

	if expectedVisibleAttrs != actualVisibleAttrs {
		attributeNames := attributesToNames(actualVisibleAttrs)
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

// getVisibleAttributesExceptBold returns the attributes that should be checked for a given cell.
// For whitespace cells, only AttrReverse (which inverts bg and fg color) and AttrStrikethrough (crossed out text) are visible.
// For non-whitespace cells, all attributes are visible
func getVisibleAttributesExceptBold(cell *uv.Cell) uint8 {
	attributesWithoutBold := cell.Style.Attrs &^ uv.AttrBold

	if isWhiteSpaceCell(cell) {
		return attributesWithoutBold & (uv.AttrReverse | uv.AttrStrikethrough)
	}

	return attributesWithoutBold
}

// isWhiteSpaceCell returns true if the cell's content is a whitespace character
func isWhiteSpaceCell(cell *uv.Cell) bool {
	return strings.TrimSpace(cell.Content) == ""
}
