package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	//	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"haiku_year/calendar"
	"haiku_year/haiku"
)

var (
	windowWidth, windowHeight             float32 = 280, 320
	todayHaiku                            []haiku.Haiku
	currentYear, currentMonth, currentDay string
)

func main() {
	currentYear, currentMonth, currentDay = "2026", "03", "" // calendar.CurrentDate()
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
	todayDate, finalText := "", ""
	if len(todayHaiku) > 0 {
		finalText = todayHaiku[0].Verse()
		todayDate = todayHaiku[0].Date()
	}
	dateLabel := widget.NewLabel(todayDate)
	verseText := widget.NewRichTextWithText(finalText)
	content := container.NewVBox(dateLabel, verseText)
	return content
}

func tabCalendar() fyne.CanvasObject {
	c := calendar.NewCalendar(todayHaiku[0].Date())
	days := c.Days()

	grid := layout.NewGridLayout(7)
	gridContainer := container.New(grid)

	for _, wd := range calendar.WeekDays("RU") {
		gridContainer.Add(widget.NewLabel(wd))
	}

	for row := 0; row < 6; row++ {
		for col := 0; col < 7; col++ {
			d := days[row][col]
			if d == "  " {
				gridContainer.Add(widget.NewLabel(d))
			} else {
				b := widget.NewButton(d, func() {
					log.Println(d)
				})
				date := fmt.Sprintf("%04s-%02s-%02s", currentYear, currentMonth, d)
				if haiku.IsHaiku(date) {
					b.Importance = widget.HighImportance
				}
				gridContainer.Add(b)
			}
		}
	}

	for _, wd := range calendar.WeekDays("JP") {
		gridContainer.Add(widget.NewLabel(wd))
	}

	content := container.NewVBox(gridContainer)
	return content
}

/*
func tabCalendarGrid() fyne.CanvasObject {
	c := calendar.NewCalendar(todayHaiku[0].Date())
	days := c.Days()

	grid := layout.NewGridLayout(7)
	gridContainer := container.New(grid)

	for _, wd := range calendar.WeekDays("RU") {
		gridContainer.Add(widget.NewLabel(wd))
	}

	for row := 0; row < len(days); row++ {
		for col := 0; col < len(days[row]); col++ {
			gridContainer.Add(widget.NewLabel(days[row][col]))
		}
	}

	for _, wd := range calendar.WeekDays("JP") {
		gridContainer.Add(widget.NewLabel(wd))
	}

	content := container.NewVBox(gridContainer)
	return content
}
*/
/*
func tabCalendarText() fyne.CanvasObject {
	calendar := calendar.NewCalendar(todayHaiku[0].Date())
	info := calendar.String()
	days := widget.NewRichTextWithText(info)
	content := container.NewVBox(days)
	return content
}
*/
