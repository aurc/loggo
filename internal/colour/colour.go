package colour

import "github.com/gdamore/tcell/v2"

const (
	//ColourBackgroundField    = tcell.Color236
	ColourBackgroundField    = tcell.ColorBlack
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
	ClField   = "[#ffaf00::b]"
	ClWhite   = "[#ffffff::-]"
	ClNumeric = "[#00afff]"
	ClString  = "[#6A9F59]"
)
