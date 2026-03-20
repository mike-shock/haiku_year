package calendar

import (
	"fmt"
	"testing"
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
