package main

import (
	//	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	//	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"haiku_year/calendar"
	"haiku_year/haiku"
)

var (
	windowWidth, windowHeight float32 = 280, 320
	todayHaiku                []haiku.Haiku
)

func main() {
	todayHaiku = haiku.Today()
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
	if len(todayHaiku) > 0 {
		final = todayHaiku[0].Verse()
	}
	verse := widget.NewRichTextWithText(final)
	content := container.NewVBox(verse)
	return content
}

func tabCalendar() fyne.CanvasObject {
	calendar := calendar.NewCalendar(todayHaiku[0].Date())
	info := calendar.String()
	days := widget.NewRichTextWithText(info)
	content := container.NewVBox(days)
	/*
		days := calendar.Days()
		list := widget.NewTable(
			func() (int, int) {
				return len(days), len(days[0])
			},
			func() fyne.CanvasObject {
				return widget.NewLabel("Calendar")
			},
			func(i widget.TableCellID, o fyne.CanvasObject) {
				o.(*widget.Label).SetText(days[i.Row][i.Col])
			})
		content := container.NewVBox(list)
	*/
	return content
}
