package calendar

import (
	"fmt"
	"testing"
	"time"
)

type TestResult struct {
	row int
	col int
	day string
}

const format1 = "2006-01-02"

var (
	testData    = [...]string{"1976-04-25", "2022-03-07", "2026-03-20"}
	testResults = [...]TestResult{TestResult{3, 6, "25"}, TestResult{1, 0, "07"}, TestResult{3, 4, "20"}}
)

func TestNewCalendar(t *testing.T) {
	for i := 0; i < len(testData); i++ {
		c := NewCalendar(testData[i])
		row, col := testResults[i].row, testResults[i].col
		day := c.days[row][col]
		if day != testResults[i].day {
			t.Errorf("invalid day: want %v, got %v", testResults[i].day, day)
		}
	}
}

func TestPrintCalendar(t *testing.T) {
	for i := 0; i < len(testData); i++ {
		c := NewCalendar(testData[i])
		c.print()
		fmt.Println(c)
		fmt.Printf("%#v\n", c)
	}
}

func TestWeekDays(t *testing.T) {
	friday := WeekDays("EN")[5-1]
	if friday != "Fr" {
		t.Errorf("invalid English name for 'Friday': %v", friday)
	}
	friday = WeekDays("RU")[5-1]
	if friday != "Пт" {
		t.Errorf("invalid Russian name for 'Friday': %v", friday)
	}
	friday = WeekDays("JP")[5-1]
	if friday != "金" {
		t.Errorf("invalid Japanese name for 'Friday': %v", friday)
	}
}

func TestCurrentDate(t *testing.T) {
	date := time.Now()
	year := date.Year()
	month := date.Month()
	day := date.Day()
	wantYear, wantMonth, wantDay := fmt.Sprintf("%04d", year), fmt.Sprintf("%02d", month), fmt.Sprintf("%02d", day)
	gotYear, gotMonth, gotDay := CurrentDate()
	if gotYear != wantYear {
		t.Errorf("invalid year for '%v': want %v, got %v", time.Now(), wantYear, gotYear)
	}
	if gotMonth != wantMonth {
		t.Errorf("invalid year for '%v': want %v, got %v", time.Now(), wantMonth, gotMonth)
	}
	if gotDay != wantDay {
		t.Errorf("invalid year for '%v': want %v, got %v", time.Now(), wantDay, gotDay)
	}
}

func TestSeason(t *testing.T) {
	d1 := "1898-07-17"
	ru := "Лето"
	gotSeason := Season(d1, "RU")
	if gotSeason != ru {
		t.Errorf("invalid season for '%v': want %v, got %v", d1, ru, gotSeason)
	}
	d2 := "1927-01-01"
	en := "Winter"
	gotSeason = Season(d2, "EN")
	if gotSeason != en {
		t.Errorf("invalid season for '%v': want %v, got %v", d2, en, gotSeason)
	}
	d3 := "1927-03-26"
	jp := "春"
	gotSeason = Season(d3, "JP")
	if gotSeason != jp {
		t.Errorf("invalid season for '%v': want %v, got %v", d3, jp, gotSeason)
	}
}

func TestThisDay(t *testing.T) {
	date := "1981-11-28"
	jp := "1981年11月28日"
	gotDate := ThisDay(date, "JP")
	if gotDate != jp {
		t.Errorf("invalid date for '%v': want %v, got %v", date, jp, gotDate)
	}

}
