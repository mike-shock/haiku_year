package calendar

import (
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	dateFormat = "2006-01-02"
)

var (
	weekDaysRU = []string{"Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"}
	weekDaysEN = []string{"Mo", "Tu", "We", "Th", "Fr", "Sa", "Su"}
	weekDaysJP = []string{"月", "火", "水", "木", "金", "土", "日"} // 曜日
)

type Calendar struct { // 暦
	date  string     // 年月日
	year  int        // 年
	month time.Month // 月
	days  [][]string // cols = week days: 0 - Mon, ..., 6 - Sun; rows = days (日)
}

func NewCalendar(date string) Calendar {
	log.Println(date)
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

func (c Calendar) Days() [][]string {
	return c.days
}

func (c Calendar) String() (s string) {
	s = fmt.Sprintf("%s\n", c.date)
	s += fmt.Sprintf("%s\n", strings.Join(weekDaysRU, " "))
	for row := 0; row < 6; row++ {
		for col := 0; col < 7; col++ {
			s += fmt.Sprintf("%2s ", c.days[row][col])
		}
		s += "\n"
	}
	return s
}

func (c Calendar) print() {
	fmt.Println(c.date)
	fmt.Println(strings.Join(weekDaysRU, " "))
	for row := 0; row < 6; row++ {
		for col := 0; col < 7; col++ {
			fmt.Printf("%2s ", c.days[row][col])
		}
		fmt.Println()
	}
}

func (c *Calendar) emptyDays() {
	c.days = make([][]string, 6)
	for i := range c.days {
		c.days[i] = make([]string, 7)
		for j := range c.days[i] {
			c.days[i][j] = "__"
		}
	}
}

func (c *Calendar) fillDays() {
	firstDay := time.Date(c.year, c.month, 1, 0, 0, 0, 0, time.UTC)
	// Weekday in Go: Sunday = 0, Monday = 1, ..., Saturday = 6
	goWeekday := firstDay.Weekday()
	// our week days numeration: 0 = Monday, 6 = Sunday
	col := (int(goWeekday) + 6) % 7
	// find out days number in this month
	lastDay := time.Date(c.year, c.month+1, 0, 0, 0, 0, 0, time.UTC).Day()
	// fill the days
	row := 0
	for day := 1; day <= lastDay; day++ {
		c.days[row][col] = fmt.Sprintf("%02d", day)
		col++
		if col == 7 {
			col = 0
			row++
		}
	}
}
