package tui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type StackedTextViews struct {
	app       *tview.Application
	rows      *tview.Flex
	statusbar *tview.Table
	statuses  map[string]*tview.TableCell
}

func NewStackedTextViews() *StackedTextViews {
	app := tview.NewApplication()

	// components
	rows := tview.NewFlex().
		SetDirection(tview.FlexRow)
	statusbar := tview.NewTable().
		SetBorders(false)
	statusbar.SetBorderPadding(1, 1, 1, 1)
	screen := &StackedTextViews{
		app:       app,
		rows:      rows,
		statusbar: statusbar,
		statuses:  make(map[string]*tview.TableCell, 0),
	}

	// event handlers
	app.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		key := e.Key()
		if key == tcell.KeyRune && e.Rune() == 'q' {
			app.Stop()
		}
		return e
	})

	// layout
	grid := tview.NewGrid().
		SetRows(3, 0).
		AddItem(statusbar, 0, 0, 1, 1, 0, 0, true).
		AddItem(rows, 1, 0, 1, 1, 0, 0, false)
	app.SetRoot(grid, true)

	// placeholders
	placeholder := tview.NewTextView().
		SetChangedFunc(func() {
			app.Draw()
		})
	placeholder.SetBorder(true).
		SetTitle("TODO").
		SetTitleAlign(tview.AlignLeft)
	placeholder.SetChangedFunc(func() {
		screen.app.Draw()
	})
	rows.AddItem(placeholder, 0, 1, false)

	return screen
}

func (s *StackedTextViews) AddStatusItem(id, label string) *StackedTextViews {
	t := s.statusbar
	n := t.GetColumnCount()

	c := tview.NewTableCell(label)
	t.SetCell(0, n, c)

	c = tview.NewTableCell("TODO").SetExpansion(2)
	t.SetCell(0, n+1, c)
	s.statuses[id] = c

	return s
}

func (s *StackedTextViews) Run() error {
	return s.app.Run()
}

func (s *StackedTextViews) Stop() {
	s.app.Stop()
}
