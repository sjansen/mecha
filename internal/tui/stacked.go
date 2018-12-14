package tui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type StackedTextViews struct {
	statusbar

	app   *tview.Application
	rows  *tview.Flex
	views []*stdview
}

func NewStackedTextViews() *StackedTextViews {
	app := tview.NewApplication()

	// components
	rows := tview.NewFlex().
		SetDirection(tview.FlexRow)
	screen := &StackedTextViews{
		app:   app,
		rows:  rows,
		views: make([]*stdview, 0),
	}
	screen.statusbar.init(app)

	// event handlers
	app.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		key := e.Key()
		if key == tcell.KeyTab {
			screen.focusNext()
			screen.app.Draw()
			return nil
		} else if key == tcell.KeyBacktab {
			screen.focusPrev()
			screen.app.Draw()
			return nil
		} else if key == tcell.KeyRune && e.Rune() == 'q' {
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

	return screen
}

func (s *StackedTextViews) AddStatusItem(label string, updates <-chan *Status) *StackedTextViews {
	s.statusbar.add(label, updates)
	return s
}

func (s *StackedTextViews) AddStdView(label string, stdout, stderr <-chan string) *StackedTextViews {
	view := &stdview{}
	view.init(s.app, label, nil, nil)
	s.rows.AddItem(view, 0, 1, false)
	s.views = append(s.views, view)
	return s
}

func (s *StackedTextViews) Run() error {
	return s.app.Run()
}

func (s *StackedTextViews) Stop() {
	s.app.Stop()
}
func (s *StackedTextViews) focusNext() {
	if s.statusbar.table.HasFocus() {
		if len(s.views) > 0 {
			box := s.views[0]
			box.SetBorderColor(tcell.ColorBlue)
			s.app.SetFocus(box)
		}
		return
	}
	for i, view := range s.views {
		if view.HasFocus() {
			view.SetBorderColor(tcell.ColorDefault)
			if (i + 1) < len(s.views) {
				box := s.views[i+1]
				box.SetBorderColor(tcell.ColorBlue)
				s.app.SetFocus(box)
			} else {
				s.app.SetFocus(s.statusbar.table)
			}
			return
		}
	}
}

func (s *StackedTextViews) focusPrev() {
	if s.statusbar.table.HasFocus() {
		if len(s.views) > 0 {
			box := s.views[len(s.views)-1]
			box.SetBorderColor(tcell.ColorBlue)
			s.app.SetFocus(box)
		}
		return
	}
	for i, view := range s.views {
		if view.HasFocus() {
			view.SetBorderColor(tcell.ColorDefault)
			if i == 0 {
				s.app.SetFocus(s.statusbar.table)
			} else {
				box := s.views[i-1]
				box.SetBorderColor(tcell.ColorBlue)
				s.app.SetFocus(box)
			}
			return
		}
	}
}
