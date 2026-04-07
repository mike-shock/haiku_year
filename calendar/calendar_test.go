package calendar

import (
	"fmt"
	"strconv"
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

func TestNextMonth(t *testing.T) {
	y := 2026
	m := 3
	d := 1
	var d1, d2 = "", ""
	for m != 2 {
		d1 = fmt.Sprintf("%04d-%02d-%02d", y, m, d)
		d2 = NextMonth(d1)
		fmt.Println("NextMonth", d1, "-->", d2)
		_, mm, _ := YyyyMmDd(d2)
		m2, _ := strconv.Atoi(mm)
		if (m < 12) && (m2 != m+1) {
			t.Errorf("invalid next date: '%v'", d2)
		}
		if (m == 12) && (m2 != 1) {
			t.Errorf("invalid next date: '%v'", d2)
		}
		m = m2
	}
}

func TestPreviousMonth(t *testing.T) {
	y := 2026
	m := 5
	d := 1
	var d1, d2 = "", ""
	for m != 6 {
		d1 = fmt.Sprintf("%04d-%02d-%02d", y, m, d)
		d2 = PreviousMonth(d1)
		fmt.Println("PreviousMonth", d1, "-->", d2)
		yyyy, mm, _ := YyyyMmDd(d2)
		y2, _ := strconv.Atoi(yyyy)
		m2, _ := strconv.Atoi(mm)
		if (m > 1) && (m2 != m-1) {
			t.Fatalf("invalid previous date: '%v'", d2)
		}
		if (m == 1) && (m2 != 12) {
			t.Fatalf("invalid previous date: '%v'", d2)
		}
		m = m2
		y = y2
	}
}

func TestNextDate(t *testing.T) {
	tests := []struct {
		name     string
		given    string
		expected string
	}{
		// End of short months (30-day months)
		{"April 30 -> May 1", "2025-04-30", "2025-05-01"},
		{"June 30 -> July 1", "2025-06-30", "2025-07-01"},
		{"September 30 -> October 1", "2025-09-30", "2025-10-01"},
		{"November 30 -> December 1", "2025-11-30", "2025-12-01"},

		// End of February (ordinary year)
		{"Feb 28 (ordinary) -> Mar 1", "2023-02-28", "2023-03-01"},
		// End of February (leap year)
		{"Feb 28 (leap) -> Feb 29", "2024-02-28", "2024-02-29"},
		{"Feb 29 (leap) -> Mar 1", "2024-02-29", "2024-03-01"},

		// End of year
		{"Dec 31 -> Jan 1 next year", "2025-12-31", "2026-01-01"},

		// Middle of month (no edge)
		{"Middle of month", "2025-03-15", "2025-03-16"},

		// Start of month (Jan 1 -> Jan 2)
		{"Jan 1 -> Jan 2", "2025-01-01", "2025-01-02"},

		// End of 31-day months
		{"Jan 31 -> Feb 1", "2025-01-31", "2025-02-01"},
		{"Mar 31 -> Apr 1", "2025-03-31", "2025-04-01"},
		{"May 31 -> Jun 1", "2025-05-31", "2025-06-01"},
		{"Jul 31 -> Aug 1", "2025-07-31", "2025-08-01"},
		{"Aug 31 -> Sep 1", "2025-08-31", "2025-09-01"},
		{"Oct 31 -> Nov 1", "2025-10-31", "2025-11-01"},
		{"Dec 31 -> Jan 1 next year", "2025-12-31", "2026-01-01"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NextDate(tt.given)
			if got != tt.expected {
				t.Errorf("nextDate(%q) = %q, want %q", tt.given, got, tt.expected)
			}
		})
	}
}

func TestPreviousDate(t *testing.T) {
	tests := []struct {
		name     string
		given    string
		expected string
	}{
		// Going backwards from start of short months
		{"May 1 -> April 30", "2025-05-01", "2025-04-30"},
		{"July 1 -> June 30", "2025-07-01", "2025-06-30"},
		{"October 1 -> September 30", "2025-10-01", "2025-09-30"},
		{"December 1 -> November 30", "2025-12-01", "2025-11-30"},

		// February borders (ordinary year)
		{"Mar 1 (ordinary) -> Feb 28", "2023-03-01", "2023-02-28"},
		{"Feb 28 (ordinary) -> Feb 27", "2023-02-28", "2023-02-27"},

		// February borders (leap year)
		{"Mar 1 (leap) -> Feb 29", "2024-03-01", "2024-02-29"},
		{"Feb 29 (leap) -> Feb 28", "2024-02-29", "2024-02-28"},
		{"Feb 28 (leap) -> Feb 27", "2024-02-28", "2024-02-27"},

		// Year boundary
		{"Jan 1 -> Dec 31 previous year", "2022-01-01", "2021-12-31"},

		// Middle of month
		{"Middle of month", "2025-03-16", "2025-03-15"},

		// End of month (31-day)
		{"Feb 1 -> Jan 31", "2025-02-01", "2025-01-31"},
		{"Apr 1 -> Mar 31", "2025-04-01", "2025-03-31"},
		{"Jun 1 -> May 31", "2025-06-01", "2025-05-31"},
		{"Sep 1 -> Aug 31", "2025-09-01", "2025-08-31"},
		{"Nov 1 -> Oct 31", "2025-11-01", "2025-10-31"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PreviousDate(tt.given)
			if got != tt.expected {
				t.Errorf("prevDate(%q) = %q, want %q", tt.given, got, tt.expected)
			}
		})
	}
}
