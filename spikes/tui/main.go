// Demo code for the TextView primitive.
package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

var BOXES = []string{"foo", "bar", "baz", "qux", "quux", "corge"}

func addView() *tview.TextView {
	textView := tview.NewTextView().
		SetDynamicColors(true)
	go func() {
		for i := 1; i <= 25; i += 1 {
			if i%10 == 0 {
				fmt.Fprintf(textView, "[red::r]line #%d[-:-:-]\n", i)
			} else {
				fmt.Fprintf(textView, "[green]line #%d[-:-:-]\n", i)
			}
			n := rand.Int()%900 + 100
			time.Sleep(time.Duration(n) * time.Millisecond)
		}
		fmt.Fprintf(textView, "[yellow]Press Enter to quit.")
	}()
	return textView
}

func main() {
	app := tview.NewApplication()
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		key := event.Key()
		if key == tcell.KeyEnter {
			app.Stop()
		} else if key == tcell.KeyRune && event.Rune() == 'q' {
			app.Stop()
		}
		return event
	})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow)
	for _, title := range BOXES {
		textView := addView().
			SetChangedFunc(func() {
				app.Draw()
			})
		textView.
			SetBorder(true).
			SetTitle(title).
			SetTitleAlign(tview.AlignLeft)
		flex.AddItem(textView, 0, 1, false)
	}

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
