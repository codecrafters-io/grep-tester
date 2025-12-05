package virtual_terminal

import (
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/x/vt"
)

type VirtualTerminal struct {
	vt   *vt.Emulator
	rows int
	cols int
}

func NewCustomVT(rows, cols int) *VirtualTerminal {
	vtInstance := &VirtualTerminal{
		vt:   vt.NewEmulator(cols, rows),
		rows: rows,
		cols: cols,
	}

	return vtInstance
}

func (vt *VirtualTerminal) Close() {
	vt.vt.Close()
}

func (vt *VirtualTerminal) Write(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	return vt.vt.Write(p)
}

func (vt *VirtualTerminal) GetScreenState() *ScreenState {
	cellMatrix := make([]*Row, vt.rows)

	for i := 0; i < vt.rows; i++ {
		cellMatrix[i] = &Row{
			cellsArray: make([]*uv.Cell, vt.cols),
		}

		for j := 0; j < vt.cols; j++ {
			cellMatrix[i].cellsArray[j] = vt.vt.CellAt(j, i)
		}
	}

	return NewScreenState(cellMatrix, CursorPosition{
		RowIndex:    vt.vt.CursorPosition().Y,
		ColumnIndex: vt.vt.CursorPosition().X,
	})
}
