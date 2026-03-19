package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"haiku-year/haiku"
)

var windowWidth, windowHeight float32 = 280, 320

func main() {
	//	haiku.HaikuDir = haikuDir
	a := app.New()
	w := a.NewWindow("Год хайку | 俳句の年")

	tabs := container.NewAppTabs()
	w.SetContent(tabs)
	tabs.Append(container.NewTabItem("Сегодня | 今日", tabToday()))
	tabs.Append(container.NewTabItem("Календарь | 暦", tabCalendar()))

	w.Resize(fyne.NewSize(windowWidth, windowHeight))
	w.CenterOnScreen()
	w.ShowAndRun()
}

func tabToday() fyne.CanvasObject {
	final := ""
	h := haiku.Today()
	if len(h) > 0 {
		final = h[0].Verse()
	}
	verse := widget.NewRichTextWithText(final)
	content := container.NewVBox(verse)
	return content
}

func tabCalendar() fyne.CanvasObject {
	info := "Здесь будет календарь на месяц."

	verse := widget.NewRichTextWithText(info)
	content := container.NewVBox(verse)
	return content
}
