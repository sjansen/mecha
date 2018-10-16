package main

import (
	"fmt"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	view := tview.NewTextView().
		SetDynamicColors(true)
	view.SetBorder(true).
		SetBorderColor(tcell.ColorBlue).
		SetBorderPadding(0, 0, 1, 1)
	view.SetChangedFunc(func() {
		app.Draw()
	})
	fmt.Fprint(view, "[::r] Press Any Key ('q' to quit) [-:-:-]\n")

	app.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		key := e.Key()
		if key == tcell.KeyRune {
			if e.Rune() == 'q' {
				app.Stop()
			}
			s := fmt.Sprintf("%#v\n", e)
			view.Write([]byte(s))
		} else {
			s := fmt.Sprintf("[green]%#v[-:-:-]\n", e)
			view.Write([]byte(s))
		}
		return e
	})

	app.SetRoot(view, true)
	app.Run()
}
