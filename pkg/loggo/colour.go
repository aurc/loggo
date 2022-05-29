package loggo

import "github.com/gdamore/tcell/v2"

const (
	ColourBackgroundField    = tcell.Color236
	ColourForegroundField    = tcell.ColorWhite
	ColourSelectedBackground = tcell.Color69
	ColourSelectedForeground = tcell.ColorWhite
	ColourSecondaryBorder    = tcell.Color240
)

var (
	FieldStyle = tcell.StyleDefault.
			Background(ColourBackgroundField).
			Foreground(ColourForegroundField)
	SelectStyle = tcell.StyleDefault.
			Background(ColourSelectedBackground).
			Foreground(ColourSelectedForeground)
)

const (
	clField   = "[#ffaf00::b]"
	clWhite   = "[#ffffff::-]"
	clNumeric = "[#00afff]"
	clString  = "[#6A9F59]"
)
