package main

import (
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	//	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"haiku_year/calendar"
	"haiku_year/haiku"
)

var (
	windowWidth, windowHeight             float32 = 280, 320
	todayHaiku                            []haiku.Haiku
	currentYear, currentMonth, currentDay string
	currentDate, selectedDate             string
	tabs                                  *container.AppTabs
	tabHaiku, tabMonth                    *container.TabItem
)

func main() {
	currentYear, currentMonth, currentDay = calendar.CurrentDate()
	currentDate = calendar.Today("RU")
	todayHaiku = haiku.Today()
	a := app.New()
	w := a.NewWindow("Год хайку | 俳句の年")

	tabs = container.NewAppTabs()
	w.SetContent(tabs)
	tabHaiku = container.NewTabItem("Сегодня | 今日", tabToday())
	tabMonth = container.NewTabItem("Календарь | 暦", tabCalendar())
	tabs.Append(tabHaiku)
	tabs.Append(tabMonth)

	w.Resize(fyne.NewSize(windowWidth, windowHeight))
	w.CenterOnScreen()
	w.ShowAndRun()
}

func tabToday() fyne.CanvasObject {
	content := setHaiku()
	return content
}

func setHaiku() *fyne.Container {
	todaySeason := calendar.Season(currentDate, "RU") + " | " + calendar.Season(currentDate, "JP")
	todayDate := calendar.ThisDay(currentDate, "RU") + " | " + calendar.ThisDay(currentDate, "JP")
	finalText, haikuDate, haikuComment, haikuAuthor := "", "", "", ""
	todayHaiku, _ = haiku.ThisDay(currentDate)
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

	box := container.NewVBox(headerLabel, verseText, footerLabel)
	return box
}

func tabCalendar() *fyne.Container {
	content := setCalendar()
	return content
}

func setCalendar() *fyne.Container {
	monthText := calendar.Month(currentDate, "RU") + " | " + calendar.Month(currentDate, "JP")
	monthLabel := widget.NewLabelWithStyle(monthText, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	grid := layout.NewGridLayout(calendar.Cols)
	gridContainer := container.New(grid)

	backButton := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), backMonth) // MediaSkipPreviousIcon
	nextButton := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), nextMonth) // MediaSkipNextIcon
	buttons := container.NewHBox(backButton, nextButton)

	c := calendar.NewCalendar(currentDate)
	days := c.Days()

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
					currentDate = fmt.Sprintf("%04s-%02s-%02s", currentYear, currentMonth, d)
					log.Println(currentDate)
					tabHaiku.Content = setHaiku()
					tabs.Select(tabHaiku)
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

	box := container.NewVBox(monthLabel, gridContainer, buttons)
	return box
}

func nextMonth() {
	currentDate = calendar.NextMonth(currentDate)
	currentYear, currentMonth, currentDay = calendar.YyyyMmDd(currentDate)
	log.Println("Next:", currentDate)
	tabMonth.Content = setCalendar()
	tabs.Select(tabMonth)
}

func backMonth() {
	currentDate := calendar.PreviousMonth(currentDate)
	currentYear, currentMonth, currentDay = calendar.YyyyMmDd(currentDate)
	log.Println("Previous:", currentDate)
	tabMonth.Content = setCalendar()
	tabs.Select(tabMonth)
}
