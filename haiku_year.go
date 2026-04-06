package main

import (
	"embed"
	"fmt"
	"image"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"haiku_year/calendar"
	"haiku_year/haiku"
)

const formatDate = "%04s-%02s-%02s"

var (
	a                                     fyne.App
	w                                     fyne.Window
	windowWidth, windowHeight             float32       = 280, 460 // 320
	todayHaiku                            []haiku.Haiku            // 今日の俳句
	todayHaikuIndex                       int           = 0
	currentYear, currentMonth, currentDay string        // 現在の年、現在の月、現在の日
	currentDate                           string        // 現在の日付
	tabs                                  *container.AppTabs
	tabHaiku, tabMonth                    *container.TabItem
	imageCheckBox                         *widget.Check
	backgroundImage                       *canvas.Image
	darkTheme                             bool = true
)

//go:embed images
var imagesDir embed.FS

func main() {
	currentYear, currentMonth, currentDay = calendar.CurrentDate()
	currentDate = calendar.Today("RU")
	todayHaiku = haiku.Today()
	a = app.NewWithID("com.shokhirev.haiku_year")
	a.Settings().SetTheme(theme.DarkTheme())
	w = a.NewWindow("Год хайку | 俳句の年")
	setDefaults()

	tabs = container.NewAppTabs()
	w.SetContent(tabs)
	tabHaiku = container.NewTabItemWithIcon("今日", theme.MediaRecordIcon(), tabToday())
	tabMonth = container.NewTabItemWithIcon("暦", theme.CalendarIcon(), tabCalendar())
	tabSettings := container.NewTabItemWithIcon("色", theme.ColorPaletteIcon(), tabOptions())
	tabAbout := container.NewTabItemWithIcon("著", theme.InfoIcon(), tabInfo())
	tabQuit := container.NewTabItemWithIcon("", theme.LogoutIcon(), tabExit())
	tabs.Append(tabHaiku)
	tabs.Append(tabMonth)
	tabs.Append(tabSettings)
	tabs.Append(tabAbout)
	tabs.Append(tabQuit)

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
	//moreButton := widget.NewButtonWithIcon("", theme.MoreHorizontalIcon(), nextVerse)
	moreButton := widget.NewButton("..!", nextVerse)
	header := container.NewHBox(headerLabel)
	if len(todayHaiku) > 1 {
		header.Add(moreButton)
	} else {
		header.Add(layout.NewSpacer())
	}
	//quitButton := widget.NewButtonWithIcon("", theme.LogoutIcon(), func() { os.Exit(0) }) // a.Quit() w.Close()
	//header.Add(quitButton)

	if len(todayHaiku) > 0 {
		finalText = todayHaiku[todayHaikuIndex].Verse()
		haikuDate = todayHaiku[todayHaikuIndex].Date()
		haikuComment = todayHaiku[todayHaikuIndex].Comment()
		haikuAuthor = todayHaiku[todayHaikuIndex].Author()
	}

	if imageCheckBox.Checked {
		backgroundImage = embeddedImage()
	} else {
		backgroundImage = colorImage()
	}

	verseText := widget.NewLabelWithStyle(finalText, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	infoText := fmt.Sprintf("%s\n%s\n%s", haikuDate, haikuAuthor, haikuComment)
	infoLabel := widget.NewLabelWithStyle(infoText, fyne.TextAlignTrailing, fyne.TextStyle{Italic: true})
	currentVerse := container.NewVBox(verseText, infoLabel)
	//haiku := container.New(layout.NewStackLayout(), backgroundImage, currentVerse)

	box := container.NewVBox(header, currentVerse)
	content := container.New(layout.NewStackLayout(), backgroundImage, box)
	return content
}

func tabCalendar() *fyne.Container {
	content := setCalendar()
	return content
}

func setCalendar() *fyne.Container {
	currentDate = currentYear + "-" + currentMonth + "-" + currentDay
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
				date := fmt.Sprintf(formatDate, currentYear, currentMonth, d)
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

func tabOptions() fyne.CanvasObject {
	content := setOptions()
	return content
}

func setOptions() *fyne.Container {
	labelTheme := widget.NewLabel("Theme | テーマ")
	themeButtons := container.NewGridWithColumns(2,
		widget.NewButtonWithIcon("Dark | 黒", theme.RadioButtonCheckedIcon(), func() { setTheme(true) }),
		widget.NewButtonWithIcon("Light | 白", theme.RadioButtonIcon(), func() { setTheme(false) }),
	)
	labelImage := widget.NewLabel("Images | 画")
	content := container.NewVBox(labelTheme, themeButtons, labelImage, imageCheckBox)
	return content
}

func setTheme(dark bool) {
	darkTheme = dark
	if darkTheme {
		a.Settings().SetTheme(theme.DarkTheme())
	} else {
		a.Settings().SetTheme(theme.LightTheme())
	}
	thisDay()
}

func setDefaults() {
	imageCheckBox = widget.NewCheck("Visible | 見", func(value bool) {
		setImage(value)
	})
	imageCheckBox.Checked = true
}

func setImage(visible bool) {
	thisDay()
}

func tabInfo() *fyne.Container {
	about := widget.NewLabel("'Haiku Year' -\n a haiku\n for each day\n of the year...")
	authors := widget.NewLabel(" by Mike & Ray Shock.")
	copyleft := widget.NewLabel("Copyleft 🄯 1999-...")
	content := container.NewVBox(about, authors, copyleft)
	return content
}

func tabExit() *fyne.Container {
	top := widget.NewLabel("Close the application.")
	bottom := widget.NewLabel("アプリケーションを閉じます。")
	quitButton := widget.NewButtonWithIcon("", theme.LogoutIcon(), func() { a.Quit() })
	spacer := layout.NewSpacer()
	box := container.NewVBox(top, quitButton, bottom)
	content := container.NewBorder(spacer, spacer, spacer, spacer, box, spacer)
	//                               (top, bottom, left, right, middle, extra)
	return content
}

func nextMonth() {
	currentDate = calendar.NextMonth(currentDate)
	currentYear, currentMonth, currentDay = calendar.YyyyMmDd(currentDate)
	tabMonth.Content = setCalendar()
	tabs.Select(tabMonth)
}

func backMonth() {
	currentDate := calendar.PreviousMonth(currentDate)
	currentYear, currentMonth, currentDay = calendar.YyyyMmDd(currentDate)
	tabMonth.Content = setCalendar()
	tabs.Select(tabMonth)
}

func thisDay() {
	currentDate = fmt.Sprintf(formatDate, currentYear, currentMonth, currentDay)
	todayHaikuIndex = 0
	tabHaiku.Content = setHaiku()
	tabs.Select(tabHaiku)
}

func nowDay() {
	currentDate = calendar.Today("RU")
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
		}
	}
}

func embeddedImage() *canvas.Image {
	fileName := calendar.Season(currentDate, "EN") + ".png"
	file, err := imagesDir.Open("images/" + fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	i := canvas.NewImageFromReader(file, fileName)
	i.FillMode = canvas.ImageFillOriginal
	return i
}

func colorImage() *canvas.Image {
	width, height := int(windowWidth), int(windowHeight)
	p := image.NewRGBA(image.Rect(0, 0, width, height))
	c := color.White
	if darkTheme {
		c = color.Black
	}
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			p.Set(x, y, c)
		}
	}
	i := canvas.NewImageFromImage(p)
	i.FillMode = canvas.ImageFillOriginal
	return i
}
