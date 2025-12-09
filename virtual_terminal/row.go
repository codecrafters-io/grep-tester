package virtual_terminal

import (
	"strings"

	uv "github.com/charmbracelet/ultraviolet"
)

type Row struct {
	cells           []*uv.Cell
	cursorCellIndex int // cursorCellIndex will be -1 if the cursor is not on the line
}

func (r *Row) hasCursor() bool {
	return r.cursorCellIndex != -1
}

// GetCells returns a copy of all the cells in the row
func (r *Row) GetCellsArray() []*uv.Cell {
	cells := make([]*uv.Cell, len(r.cells))

	for i, cell := range r.cells {
		cells[i] = cell.Clone()
	}

	return cells
}

func (r *Row) Clone() *Row {
	return &Row{
		cursorCellIndex: r.cursorCellIndex,
		cells:           r.GetCellsArray(),
	}
}

func (r *Row) GetCellsCount() int {
	return len(r.cells)
}

// GetContents returns the contents present in given row
// If cursor is not present in the row, the spaces are trimmed from the right
// If the row has a cursor, the spaces before the cursor is preserved,
// and spaces are trimmed from the right
func (r *Row) GetContents() string {
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
