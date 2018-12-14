package tui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

type Severity int

const (
	Unknown Severity = iota
	Refresh
	Healthy
	Warning
	Alert
)

type Status struct {
	Severity Severity
	Message  string
}

type statusbar struct {
	*tview.Table

	app *tview.Application
}

func (s *statusbar) init(app *tview.Application) {
	s.app = app
	s.Table = tview.NewTable().
		SetBorders(false)
	s.SetBorderPadding(1, 1, 1, 1)
}

func (s *statusbar) add(label string, updates <-chan *Status) *statusbar {
	n := s.GetColumnCount()

	c := tview.NewTableCell(label)
	s.SetCell(0, n, c)

	c = tview.NewTableCell("").SetExpansion(2)
	s.SetCell(0, n+1, c)

	go func() {
		for status := range updates {
			s.app.QueueUpdateDraw(func() {
				updateStatusCell(c, status)
			})
		}
	}()

	return s
}

func updateStatusCell(c *tview.TableCell, status *Status) {
	msg := status.Message
	switch status.Severity {
	case Refresh:
		c.SetText(msg).SetTextColor(tcell.ColorBlue)
	case Healthy:
		c.SetText(msg).SetTextColor(tcell.ColorGreen)
	case Warning:
		c.SetText(msg).SetTextColor(tcell.ColorYellow)
	case Alert:
		c.SetText(msg).SetTextColor(tcell.ColorRed)
	default:
		c.SetText(msg).SetTextColor(tcell.ColorDefault)
	}
}
