package virtual_terminal

import (
	"strings"

	uv "github.com/charmbracelet/ultraviolet"
)

type row struct {
	cells           []*uv.Cell
	cursorCellIndex int // cursorCellIndex will be -1 if the cursor is not on the line
}

func (r *row) hasCursor() bool {
	return r.cursorCellIndex != -1
}

// GetCells returns a copy of all the cells in the row
func (r *row) getCellsArray() []*uv.Cell {
	cells := make([]*uv.Cell, len(r.cells))

	for i, cell := range r.cells {
		cells[i] = cell.Clone()
	}

	return cells
}

func (r *row) GetCellsCount() int {
	return len(r.cells)
}

// getTextContents returns the contents present in given row
// Spaces on the right are not preserved
func (r *row) getTextContents() string {
	rawCellContents := ""

	for _, cell := range r.cells {
		rawCellContents += cell.Content
	}

	if r.hasCursor() {
		// If the cursor is on the line, we need to preserve the spaces before the cursor
		contentsBeforeCursor := rawCellContents[:r.cursorCellIndex]
		contentsAfterCursor := strings.TrimRight(rawCellContents[r.cursorCellIndex:], " ")
		return contentsBeforeCursor + contentsAfterCursor
	} else {
		// If the cursor isn't on the line, we can safely trim spaces from the right
		return strings.TrimRight(rawCellContents, " ")
	}
}
