package virtual_terminal

import (
	"strings"

	"github.com/charmbracelet/x/vt"
)

type VirtualTerminal struct {
	emulator     *vt.Emulator
	rowsCount    int
	columnsCount int
}

func NewStandardVT() *VirtualTerminal {
	return NewCustomVT(80, 124)
}

func NewCustomVT(rowsCount, columnsCount int) *VirtualTerminal {
	return &VirtualTerminal{
		emulator:     vt.NewEmulator(columnsCount, rowsCount),
		rowsCount:    rowsCount,
		columnsCount: columnsCount,
	}
}

func (vt *VirtualTerminal) Write(p []byte) (int, error) {
	return vt.emulator.Write(p)
}

func (vt *VirtualTerminal) RenderText() string {
	lines := make([]string, vt.emulator.CursorPosition().Y)

	for row := range vt.emulator.CursorPosition().Y {
		rowStringBuilder := strings.Builder{}

		for column := range vt.columnsCount {
			cellContent := vt.emulator.CellAt(column, row).Content
			rowStringBuilder.WriteString(cellContent)
		}

		lines[row] = strings.TrimSpace(rowStringBuilder.String())
	}

	result := strings.Join(lines, "\n")
	return result
}
