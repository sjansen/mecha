package tui

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Screen struct {
	app       *tview.Application
	menu      *tview.List
	menuItems map[string]MenuAction
	rows      *tview.Flex
	views     []*tview.TextView
}

type MenuAction func()

func NewScreen() *Screen {
	// components
	menu := tview.NewList().
		SetSelectedBackgroundColor(tcell.ColorBlue).
		ShowSecondaryText(false)
	menu.SetBorderPadding(1, 1, 1, 1)
	rows := tview.NewFlex().
		SetDirection(tview.FlexRow)
	// layout
	grid := tview.NewGrid().
		SetColumns(15, 0).
		SetRows(0).
		AddItem(menu, 0, 0, 1, 1, 0, 0, true).
		AddItem(rows, 0, 1, 1, 1, 0, 0, false)
	// event handlers
	app := tview.NewApplication().
		SetRoot(grid, true)
	screen := &Screen{
		app:       app,
		menu:      menu,
		menuItems: make(map[string]MenuAction, 0),
		rows:      rows,
		views:     make([]*tview.TextView, 0),
	}
	app.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		key := e.Key()
		if key == tcell.KeyTab {
			screen.changeFocus()
			screen.app.Draw()
			return nil
		} else if key == tcell.KeyRune && e.Rune() == 'q' {
			app.Stop()
		}
		return e
	})
	menu.SetSelectedFunc(func(_ int, name string, _ string, _ rune) {
		if action, ok := screen.menuItems[name]; ok {
			action()
		}
	})
	return screen
}

func (s *Screen) AddMenuItem(name string, action MenuAction) *Screen {
	s.menu.AddItem(name, "", 0, nil)
	s.menuItems[name] = action
	return s
}

func (s *Screen) AddStreamPair(title string, stdout, stderr <-chan string) {
	view := tview.NewTextView().
		SetDynamicColors(true)
	view.SetBorder(true).
		SetTitle(title).
		SetTitleAlign(tview.AlignLeft)
	view.SetChangedFunc(func() {
		s.app.Draw()
	})
	s.rows.AddItem(view, 0, 1, false)
	s.views = append(s.views, view)
	go func() {
	loop:
		for {
			select {
			case line, ok := <-stdout:
				if ok {
					fmt.Fprint(view, "[green]")
					fmt.Fprint(view, line)
					fmt.Fprint(view, "[-:-:-]\n")
				} else {
					stdout = nil
					if stderr == nil {
						break loop
					}
				}
			case line, ok := <-stderr:
				if ok {
					fmt.Fprint(view, "[red]")
					fmt.Fprint(view, line)
					fmt.Fprint(view, "[-:-:-]\n")
				} else {
					stderr = nil
					if stdout == nil {
						break loop
					}
				}
			}
		}
		fmt.Fprintf(view, `[::r]Press "q" to quit.`)
	}()
}

func (s *Screen) Run() error {
	return s.app.Run()
}

func (s *Screen) Stop() {
	s.app.Stop()
}

func (s *Screen) changeFocus() {
	if s.menu.HasFocus() {
		s.menu.SetSelectedBackgroundColor(tcell.ColorWhite)
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
				s.menu.SetSelectedBackgroundColor(tcell.ColorBlue)
				s.app.SetFocus(s.menu)
			}
			return
		}
	}
}
