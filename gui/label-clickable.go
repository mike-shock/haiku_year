package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type ClickableLabel struct {
	widget.Label
	onTapped func()
}

func (c *ClickableLabel) Tapped(ev *fyne.PointEvent) {
	if c.onTapped != nil {
		c.onTapped()
	}
}

func NewClickableLabel(text string, tapped func()) *ClickableLabel {
	l := &ClickableLabel{onTapped: tapped}
	l.ExtendBaseWidget(l)
	l.SetText(text)
	return l
}

// clickableLabel := newClickableLabel("Click me!", func() { println("Clicked!") })
