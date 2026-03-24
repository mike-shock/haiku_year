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
	currentDate, selectedDate             string
)

func main() {
	currentYear, currentMonth, currentDay = calendar.CurrentDate()
	currentDate = calendar.Today("RU")
	todayHaiku = haiku.Today()
	a := app.New()
	w := a.NewWindow("Год хайку | 俳句の年")

	tabs := container.NewAppTabs()
	w.SetContent(tabs)
	tabHaiku := container.NewTabItem("Сегодня | 今日", tabToday())
	tabMonth := container.NewTabItem("Календарь | 暦", tabCalendar())
	tabs.Append(tabHaiku)
	tabs.Append(tabMonth)

	w.Resize(fyne.NewSize(windowWidth, windowHeight))
	w.CenterOnScreen()
	w.ShowAndRun()
}

func tabToday() fyne.CanvasObject {
	todaySeason := calendar.Season(currentDate, "RU") + " | " + calendar.Season(currentDate, "JP")
	todayDate := calendar.ThisDay(currentDate, "RU") + " | " + calendar.ThisDay(currentDate, "JP")
	finalText, haikuDate, haikuComment, haikuAuthor := "", "", "", ""
	if len(todayHaiku) > 0 {
		finalText = todayHaiku[0].Verse()
		haikuDate = todayHaiku[0].Date()
		haikuComment = todayHaiku[0].Comment()
		haikuAuthor = todayHaiku[0].Author()
	}

	headerText := fmt.Sprintf("%s\n%s\n", todaySeason, todayDate)
	headerLabel := widget.NewLabel(headerText)

	verseText := widget.NewLabelWithStyle(finalText, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	footerText := fmt.Sprintf("%s\n%s\n%s", haikuDate, haikuAuthor, haikuComment)
	footerLabel := widget.NewLabelWithStyle(footerText, fyne.TextAlignTrailing, fyne.TextStyle{Italic: true})

	content := container.NewVBox(headerLabel, verseText, footerLabel)
	return content
}

func tabCalendar() fyne.CanvasObject {
	c := calendar.NewCalendar(currentDate)
	days := c.Days()

	grid := layout.NewGridLayout(calendar.Cols)
	gridContainer := container.New(grid)

	for _, wd := range calendar.WeekDays("RU") {
		gridContainer.Add(widget.NewLabel(wd))
	}

	for row := 0; row < calendar.Rows; row++ {
		for col := 0; col < calendar.Cols; col++ {
			d := days[row][col]
			if d == "  " {
				gridContainer.Add(widget.NewLabel(d))
			} else {
				b := widget.NewButton(d, func() {
					selectedDate = fmt.Sprintf("%04s-%02s-%02s", currentYear, currentMonth, d)
					log.Println(selectedDate)
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

	monthText := calendar.Month(currentDate, "RU") + " | " + calendar.Month(currentDate, "JP")
	monthLabel := widget.NewLabelWithStyle(monthText, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	content := container.NewVBox(monthLabel, gridContainer)
	return content
}
