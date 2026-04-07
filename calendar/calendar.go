package calendar

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	dateFormat = "2006-01-02"
	zeroDate   = "0000-00-00"
	Rows       = 6
	Cols       = 7
)

var (
	weekDays = map[string][]string{
		"RU": []string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"},
		"EN": []string{"Mo", "Tu", "We", "Th", "Fr", "Sa", "Su"},
		"JP": []string{"月", "火", "水", "木", "金", "土", "日"}, // 曜日
	}
	months = map[string][]string{
		"RU": []string{"", "Январь", "Февраль", "Март", "Апрель", "Май", "Июнь", "Июль", "Август", "Сентябрь", "Октябрь", "Ноябрь", "Декабрь"},
		"EN": []string{"", "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"},
		"JP": []string{"", "一", "二", "三", "四", "五", "六", "七", "八", "九", "十", "十一", "十二"},
	}
	seasons = map[string][]string{
		"RU": []string{"", "Весна", "Лето", "Осень", "Зима"},
		"EN": []string{"", "Spring", "Summer", "Autumn", "Winter"},
		"JP": []string{"", "春", "夏", "秋", "冬"},
	}
	formats = map[string]string{
		"RU": "%04d-%02d-%02d",
		//		"EN": "%04d-%02d-%02d",
		"JP": "%04s年%02s月%02s日",
	}
)

type Calendar struct { // 暦
	date  string     // 年月日
	year  int        // 年
	month time.Month // 月
	days  [][]string // cols = week days: 0 - Mon, ..., 6 - Sun; rows = days (日)
}

func NewCalendar(date string) Calendar {
	c := Calendar{date: date}
	ymd, err := time.Parse(dateFormat, date)
	if err != nil {
		log.Printf("error in date format: %v\n", err)
		return c
	}
	c.year, c.month, _ = ymd.Date()
	c.emptyDays()
	c.fillDays()
	return c
}

func Today(language string) (date string) {
	today := time.Now().Format(dateFormat) // 今日
	date = ThisDay(today, language)
	return date
}

func ThisDay(someDate, language string) (date string) {
	switch language {
	case "RU", "EN":
		date = someDate // YYYY-MM-DD
	case "JP":
		ymd := strings.Split(someDate, "-")
		date = fmt.Sprintf(formats["JP"], ymd[0], ymd[1], ymd[2])
	}
	return date
}

func Season(date, language string) string {
	ymd := strings.Split(date, "-")
	if len(ymd) == 3 {
		if m, err := strconv.Atoi(ymd[1]); err == nil {
			s := 0
			switch m {
			case 3, 4, 5:
				s = 1
			case 6, 7, 8:
				s = 2
			case 9, 10, 11:
				s = 3
			case 12, 1, 2:
				s = 4
			}
			return seasons[language][s]
		}
	}
	return ""
}

func Month(date, language string) (month string) {
	ymd := strings.Split(date, "-")
	if m, err := strconv.Atoi(ymd[1]); err == nil {
		if m >= 1 && m <= 12 {
			month = months[language][m]
		}
	}
	return month
}

func NextMonth(date string) string {
	t, err := time.Parse(dateFormat, date)
	if err != nil {
		return date
	}
	y, m, d := t.Date()
	y, m, d = time.Date(y, m+1, 1, 0, 0, 0, 0, time.UTC).Date()
	return fmt.Sprintf(formats["RU"], y, int(m), d)
}

func PreviousMonth(date string) string {
	t, err := time.Parse(dateFormat, date)
	if err != nil {
		return date
	}
	y, m, d := t.Date()
	if m == 1 {
		y, m, d = time.Date(y-1, 12, 1, 0, 0, 0, 0, time.UTC).Date()
	} else {
		y, m, d = time.Date(y, m-1, 1, 0, 0, 0, 0, time.UTC).Date()
	}
	return fmt.Sprintf(formats["RU"], y, int(m), d)
}

func WeekDays(language string) []string {
	return weekDays[language]
}

func (c Calendar) Days() [][]string {
	return c.days
}

func (c Calendar) String() (s string) {
	s = fmt.Sprintf("%s\n", c.date)
	s += fmt.Sprintf("%s\n", strings.Join(weekDays["RU"], " "))
	for row := 0; row < Rows; row++ {
		for col := 0; col < Cols; col++ {
			s += fmt.Sprintf("%2s ", c.days[row][col])
		}
		s += "\n"
	}
	return s
}

func (c Calendar) print() {
	fmt.Println(c)
}

func (c *Calendar) emptyDays() {
	c.days = make([][]string, Rows)
	for i := range c.days {
		c.days[i] = make([]string, Cols)
		for j := range c.days[i] {
			c.days[i][j] = "  " // __"
		}
	}
}

func (c *Calendar) fillDays() {
	firstDay := time.Date(c.year, c.month, 1, 0, 0, 0, 0, time.UTC)
	// Weekday in Go: Sunday = 0, Monday = 1, ..., Saturday = 6
	goWeekday := firstDay.Weekday()
	// our week days numeration: 0 = Monday, 6 = Sunday
	col := (int(goWeekday) + 6) % Cols
	// find out days number in this month
	lastDay := time.Date(c.year, c.month+1, 0, 0, 0, 0, 0, time.UTC).Day()
	// fill the days
	row := 0
	for day := 1; day <= lastDay; day++ {
		c.days[row][col] = fmt.Sprintf("%02d", day)
		col++
		if col == Cols {
			col = 0
			row++
		}
	}
}

func CurrentDate() (currentYear, currentMonth, currentDay string) {
	year, month, day := time.Now().Date()
	currentYear, currentMonth, currentDay = fmt.Sprintf("%04d", year), fmt.Sprintf("%02d", month), fmt.Sprintf("%02d", day)
	return currentYear, currentMonth, currentDay
}

func YyyyMmDd(date string) (y, m, d string) {
	ymd := strings.Split(date, "-")
	return ymd[0], ymd[1], ymd[2]
}

func nextDate(givenDate string) string {
	t, err := time.Parse(dateFormat, givenDate)
	if err != nil {
		return zeroDate
	}
	next := t.AddDate(0, 0, 1)
	return next.Format(dateFormat)
}

func previousDate(givenDate string) string {
	t, err := time.Parse(dateFormat, givenDate)
	if err != nil {
		return zeroDate
	}
	prev := t.AddDate(0, 0, -1)
	return prev.Format(dateFormat)
}
