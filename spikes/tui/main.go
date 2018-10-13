package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type menuAction func()

type screen struct {
	app       *tview.Application
	menu      *tview.List
	menuItems map[string]menuAction
	rows      *tview.Flex
	views     []*tview.TextView
}

func newScreen() *screen {
	// components
	menu := tview.NewList().
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
	screen := &screen{
		app:       app,
		menu:      menu,
		menuItems: make(map[string]menuAction, 0),
		rows:      rows,
		views:     make([]*tview.TextView, 0),
	}
	app.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		key := e.Key()
		if key == tcell.KeyRune && e.Rune() == 'q' {
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

func (s *screen) addMenuItem(name string, action menuAction) *screen {
	s.menu.AddItem(name, "", 0, nil)
	s.menuItems[name] = action
	return s
}

func (s *screen) addStreamPair(title string, stdout, stderr <-chan string) {
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

func (s *screen) run() error {
	return s.app.Run()
}

func (s *screen) stop() {
	s.app.Stop()
}

func main() {
	screen := newScreen()
	addStreamPair := func() {
		stdout := make(chan string)
		stderr := make(chan string)
		screen.addStreamPair("TODO", stdout, stderr)
		go func() {
			for i := 1; i <= 15; i++ {
				if i%10 == 0 {
					stderr <- fmt.Sprintf("line #%d", i)
				} else {
					stdout <- fmt.Sprintf("line #%d", i)
				}
				n := rand.Int()%750 + 250
				time.Sleep(time.Duration(n) * time.Millisecond)
			}
			close(stdout)
			close(stderr)
		}()
	}
	screen.
		addMenuItem("Add Row", addStreamPair).
		addMenuItem("Quit Row", screen.stop)
	for i := 0; i < 3; i++ {
		addStreamPair()
	}
	if err := screen.run(); err != nil {
		panic(err)
	}
}
