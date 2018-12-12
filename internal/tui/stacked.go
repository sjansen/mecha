package tui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type StackedTextViews struct {
	statusbar

	app  *tview.Application
	rows *tview.Flex
}

func NewStackedTextViews() *StackedTextViews {
	app := tview.NewApplication()

	// components
	rows := tview.NewFlex().
		SetDirection(tview.FlexRow)
	screen := &StackedTextViews{
		app:  app,
		rows: rows,
	}
	screen.statusbar.init(app)

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
		AddItem(screen.statusbar.table, 0, 0, 1, 1, 0, 0, true).
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

func (s *StackedTextViews) AddStatusItem(id, label string, updates <-chan *Status) *StackedTextViews {
	s.statusbar.add(id, label, updates)
	return s
}

func (s *StackedTextViews) Run() error {
	return s.app.Run()
}

func (s *StackedTextViews) Stop() {
	s.app.Stop()
}
