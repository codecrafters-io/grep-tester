package assertions

import (
	"fmt"
	"strings"

	uv "github.com/charmbracelet/ultraviolet"
)

// vtCellComparisonContext holds cell information for comparison
type vtCellComparisonContext struct {
	cellRowIdx    int
	cellColumnIdx int
	expectedCell  *uv.Cell
	actualCell    *uv.Cell
	successLogs   []string
}

func newEmptVtCellComparisonContext() *vtCellComparisonContext {
	return &vtCellComparisonContext{}
}

func (ctx *vtCellComparisonContext) updateRowIndex(newIdx int) {
	ctx.cellRowIdx = newIdx
}

func (ctx *vtCellComparisonContext) updateColumnIndex(newIdx int) {
	ctx.cellColumnIdx = newIdx
}

func (ctx *vtCellComparisonContext) updateExpectedCell(newExpectedCell *uv.Cell) {
	ctx.expectedCell = newExpectedCell
}

func (ctx *vtCellComparisonContext) updateActualCell(newActualCell *uv.Cell) {
	ctx.actualCell = newActualCell
}

func (ctx *vtCellComparisonContext) resetSuccessLogs() {
	ctx.successLogs = []string{}
}

func (ctx *vtCellComparisonContext) addSuccessLog(successLog string) {
	ctx.successLogs = append(ctx.successLogs, successLog)
}

func (ctx *vtCellComparisonContext) checkCellForegroundColor() error {
	expectedCell := ctx.expectedCell
	actualCell := ctx.actualCell

	if expectedCell.Style.Fg != actualCell.Style.Fg {

		// No color expected, use a different wording
		if expectedCell.Style.Fg == nil {
			return fmt.Errorf("Expected no color, got %s", getFgColorName(actualCell.Style.Fg))
		}

		// No color actually, use a different wording
		if actualCell.Style.Fg == nil {
			return fmt.Errorf("Expected color to be %s, got no color", getFgColorName(expectedCell.Style.Fg))
		}

		return fmt.Errorf("Expected color to be %s, got %s", getFgColorName(expectedCell.Style.Fg), getFgColorName(actualCell.Style.Fg))
	}

	ctx.addSuccessLog(fmt.Sprintf("✓ Color is %s", getFgColorName(expectedCell.Style.Fg)))
	return nil
}

func (ctx *vtCellComparisonContext) checkCellBoldAttribute() error {
	expectedCell := ctx.expectedCell
	actualCell := ctx.actualCell

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

func (ctx *vtCellComparisonContext) checkInvalidStylingAndAttributes() error {
	// Intentionally don't add success messages for extra attributes

	if ctx.actualCell.Style.Attrs != ctx.expectedCell.Style.Attrs {
		attributesWithoutBold := ctx.actualCell.Style.Attrs &^ uv.AttrBold
		attributeNames := attributesToNames(attributesWithoutBold)
		return fmt.Errorf(
			"Expected no extra attributes, got attributes: %s",
			strings.Join(attributeNames, ", "),
		)
	}

	if ctx.actualCell.Style.Bg != ctx.expectedCell.Style.Bg {
		return fmt.Errorf("Expected no background color, found %s", getBgColorName(ctx.actualCell.Style.Bg))
	}

	if ctx.actualCell.Style.Underline != ctx.expectedCell.Style.Underline {
		return fmt.Errorf("Expected no underline, found underline")
	}

	return nil
}
