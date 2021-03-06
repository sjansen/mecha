package tui

import (
	"fmt"

	"github.com/rivo/tview"
)

type stdview struct {
	*tview.TextView

	app *tview.Application
}

func (v *stdview) init(app *tview.Application, label string, stdout, stderr <-chan string, updates <-chan *Status) {
	view := tview.NewTextView()
	v.TextView = view
	v.app = app

	view.SetDynamicColors(true).
		SetBorder(true).
		SetTitle(label).
		SetTitleAlign(tview.AlignLeft)
	view.SetChangedFunc(func() {
		app.Draw()
	})

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
					fmt.Fprint(view, "[blue::b]stdout closed[-:-:-]\n")
					stdout = nil
					if stderr == nil && updates == nil {
						break loop
					}
				}
			case line, ok := <-stderr:
				if ok {
					fmt.Fprint(view, "[red]")
					fmt.Fprint(view, line)
					fmt.Fprint(view, "[-:-:-]\n")
				} else {
					fmt.Fprint(view, "[blue::b]stderr closed[-:-:-]\n")
					stderr = nil
					if stdout == nil && updates == nil {
						break loop
					}
				}
			case update, ok := <-updates:
				if ok {
					app.QueueUpdateDraw(func() {
						title := label + "(" + update.Message + ")"
						view.SetTitle(title)
					})
				} else {
					updates = nil
					if stderr == nil && stdout == nil {
						break loop
					}
				}
			}
		}
	}()
}
