// Demo code for the TextView primitive.
package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var titles = []string{
	"foo", "bar", "baz", "qux", "quux", "corge",
	"grault", "garply", "waldo", "fred", "plugh",
	"xyzzy", "thud",
}

type screen struct {
	app   *tview.Application
	grid  *tview.Grid
	rows  *tview.Flex
	views []*tview.TextView
}

func (s *screen) addView() {
	idx := len(s.views)
	view := tview.NewTextView().
		SetDynamicColors(true)
	view.SetBorder(true).
		SetTitle(titles[idx]).
		SetTitleAlign(tview.AlignLeft)
	view.SetChangedFunc(func() {
		s.app.Draw()
	})
	go func() {
		for i := 1; i <= 15; i++ {
			if i%10 == 0 {
				fmt.Fprintf(view, "[red::r]line #%d[-:-:-]\n", i)
			} else {
				fmt.Fprintf(view, "[green]line #%d[-:-:-]\n", i)
			}
			n := rand.Int()%750 + 250
			time.Sleep(time.Duration(n) * time.Millisecond)
		}
		fmt.Fprintf(view, `[yellow]Press "q" to quit.`)
	}()
	s.rows.AddItem(view, 0, 1, false)
	s.views = append(s.views, view)
}

func main() {
	app := tview.NewApplication()
	grid := tview.NewGrid().
		SetColumns(15, 0).
		SetRows(0)
	rows := tview.NewFlex().
		SetDirection(tview.FlexRow)
	screen := screen{
		app:   app,
		grid:  grid,
		rows:  rows,
		views: make([]*tview.TextView, 0, 3),
	}

	menu := tview.NewList().
		AddItem("Add Row", "", 0, nil).
		AddItem("Quit", "", 0, nil).
		ShowSecondaryText(false)
	menu.SetBorderPadding(1, 1, 1, 1)

	grid.AddItem(menu, 0, 0, 1, 1, 0, 0, true)
	grid.AddItem(rows, 0, 1, 1, 1, 0, 0, false)

	app.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		key := e.Key()
		if key == tcell.KeyRune && e.Rune() == 'q' {
			app.Stop()
		}
		return e
	})

	menu.SetSelectedFunc(func(_ int, action string, _ string, _ rune) {
		switch action {
		case "Add Row":
			screen.addView()
		case "Quit":
			app.Stop()
		}
	})

	for i := 0; i < cap(screen.views); i++ {
		screen.addView()
	}

	if err := app.SetRoot(grid, true).Run(); err != nil {
		panic(err)
	}
}
