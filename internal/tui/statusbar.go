package tui

import "github.com/rivo/tview"

type statusbar struct {
	table *tview.Table
	cells map[string]*tview.TableCell
}

func (s *statusbar) init() {
	s.table = tview.NewTable().
		SetBorders(false)
	s.table.SetBorderPadding(1, 1, 1, 1)

	s.cells = make(map[string]*tview.TableCell, 0)
}

func (s *statusbar) add(id, label string) *statusbar {
	t := s.table
	n := t.GetColumnCount()

	c := tview.NewTableCell(label)
	t.SetCell(0, n, c)

	c = tview.NewTableCell("TODO").SetExpansion(2)
	t.SetCell(0, n+1, c)
	s.cells[id] = c

	return s
}
