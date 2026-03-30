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
	windowWidth, windowHeight             float32       = 280, 320
	todayHaiku                            []haiku.Haiku // 今日の俳句
	todayHaikuIndex                       int           = 0
	currentYear, currentMonth, currentDay string        // 現在の年、現在の月、現在の日
	currentDate, selectedDate             string        // 現在の日付
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

	headerLabel := widget.NewLabel(fmt.Sprintf("%s\n%s\n", todaySeason, todayDate))
	moreButton := widget.NewButtonWithIcon("", theme.MoreHorizontalIcon(), nextVerse)
	header := container.NewHBox(headerLabel)
	if len(todayHaiku) > 1 {
		header.Add(moreButton)
	}

	if len(todayHaiku) > 0 {
		finalText = todayHaiku[todayHaikuIndex].Verse()
		haikuDate = todayHaiku[todayHaikuIndex].Date()
		haikuComment = todayHaiku[todayHaikuIndex].Comment()
		haikuAuthor = todayHaiku[todayHaikuIndex].Author()
	}
	verseText := widget.NewLabelWithStyle(finalText, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	infoText := fmt.Sprintf("%s\n%s\n%s", haikuDate, haikuAuthor, haikuComment)
	infoLabel := widget.NewLabelWithStyle(infoText, fyne.TextAlignTrailing, fyne.TextStyle{Italic: true})
	currentVerse := container.NewVBox(verseText, infoLabel)

	box := container.NewVBox(header, currentVerse)
	return box
}

func tabCalendar() *fyne.Container {
	content := setCalendar()
	return content
}

func setCalendar() *fyne.Container {
	currentDate = currentYear + "-" + currentMonth + "-" + currentDay
	fmt.Println("setCalendar", currentDate, currentYear, currentMonth, currentDay)
	monthText := calendar.Month(currentDate, "RU") + " | " + calendar.Month(currentDate, "JP")
	monthLabel := widget.NewLabelWithStyle(monthText, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	grid := layout.NewGridLayout(calendar.Cols)
	gridContainer := container.New(grid)

	backButton := widget.NewButtonWithIcon("", theme.MediaFastRewindIcon(), backMonth)
	todayButton := widget.NewButtonWithIcon("", theme.MediaRecordIcon(), nowDay)
	nextButton := widget.NewButtonWithIcon("", theme.MediaFastForwardIcon(), nextMonth)
	buttons := container.NewHBox(layout.NewSpacer(), backButton, todayButton, nextButton, layout.NewSpacer())

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
					currentDay = d
					thisDay()
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
	cD := currentDate
	currentDate := calendar.PreviousMonth(currentDate)
	currentYear, currentMonth, currentDay = calendar.YyyyMmDd(currentDate)
	log.Println("Previous:", cD, "-->", currentDate)
	tabMonth.Content = setCalendar()
	tabs.Select(tabMonth)
}

func thisDay() {
	currentDate = fmt.Sprintf("%04s-%02s-%02s", currentYear, currentMonth, currentDay)
	log.Println(currentDate)
	todayHaikuIndex = 0
	tabHaiku.Content = setHaiku()
	tabs.Select(tabHaiku)
}

func nowDay() {
	currentDate = calendar.Today("RU")
	log.Println(currentDate)
	todayHaikuIndex = 0
	tabHaiku.Content = setHaiku()
	tabs.Select(tabHaiku)
}

func nextVerse() {
	currentIndex := todayHaikuIndex
	if len(todayHaiku) > 1 {
		todayHaikuIndex++
		if todayHaikuIndex == len(todayHaiku) {
			todayHaikuIndex = 0
		}
		if todayHaikuIndex != currentIndex {
			tabHaiku.Content = setHaiku()
			//tabs.Select(tabHaiku)
		}
	}

}
