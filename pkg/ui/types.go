package ui

import "github.com/rivo/tview"

type FocusDelegator interface {
	SetFocus(p tview.Primitive) *tview.Application
	QueueUpdateDraw(f func()) *tview.Application
}
