package virtual_terminal

import (
	"strings"

	uv "github.com/charmbracelet/ultraviolet"
)

type Row struct {
	cellsArray      []*uv.Cell
	cursorCellIndex int // cursorCellIndex will be -1 if the cursor is not on the line
}

func (r *Row) hasCursor() bool {
	return r.cursorCellIndex != -1
}

func (r *Row) IsEmpty() bool {
	return r.String() == ""
}

func (r *Row) GetCellsArray() []*uv.Cell {
	return r.cellsArray
}

func (r *Row) String() string {
	rawCellContents := ""

	for _, cell := range r.cellsArray {
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
