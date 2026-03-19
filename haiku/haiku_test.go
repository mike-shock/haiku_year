package haiku

import (
	//	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

const testText = `
Это многострочный текст
на русском языке,
в котором есть специальные вставки:

{25.04.1976}
[Ray]
<это коммментарий>
А дальше - ещё текст.
`

type TestDatum struct {
	input          string
	outputPositive int
	outputNegative error
	outputError    error
}

var (
	testData = []TestDatum{
		TestDatum{"", 0, EmptyDateError, nil},
		TestDatum{"-", 0, BadDelimiterError, nil},
		TestDatum{"1976.04.05", 0, BadDelimiterError, nil},
		TestDatum{"0000-00-00", 0, InvalidDateError, nil},
		TestDatum{"0-1-1", 0, InvalidDateError, nil},
		TestDatum{"1-0-1", 0, InvalidDateError, nil},
		TestDatum{"1-1-0", 0, InvalidDateError, nil},
		TestDatum{"1978-2", 0, BadDelimiterError, nil},
		TestDatum{"1978-2.3", 0, BadDelimiterError, nil},
		TestDatum{"2020-5-7", 3, nil, nil},
		TestDatum{"2026-01-01", 1, nil, nil},
		TestDatum{"2025-01-02", 0, nil, TextMissingError},
		TestDatum{"2026-02-05", 3, nil, nil},
		TestDatum{"2004-03-06", 1, nil, nil},
		TestDatum{"2023-10-16", 2, nil, nil},
		TestDatum{"2022-12-31", 1, nil, nil},
		TestDatum{"9999-12-32", 0, InvalidDateError, nil},
		TestDatum{"8888-13-01", 0, InvalidDateError, nil},
	}
	testDate = "2022-03-07"
)

func TestNewHaiku(t *testing.T) {
	ymd := strings.Split(testDate, "-")
	h := NewHaiku(testDate)
	if h.day != ymd[2] {
		t.Errorf("invalid day of %v: %v != %v", testDate, h.day, ymd[2])
	}
	if h.month != ymd[1] {
		t.Errorf("invalid month of %v: %v != %v", testDate, h.month, ymd[1])
	}
	wantYear := ""
	if h.year != wantYear {
		t.Errorf("invalid year of %v: %v != %v", testDate, h.year, wantYear)
	}
}

func TestPretextHaiku(t *testing.T) {
	h, err := pretext()
	if err != nil {
		t.Fatalf("default haiku reading error: '%v'", err)
	}
	if len(h) == 0 {
		t.Fatalf("default haiku not found: '%v'", h)
	}
}

func TestTodayHaiku(t *testing.T) {
	h := Today()
	if len(h) == 0 {
		t.Fatalf("no haiku found for today(): '%v'", h)
	} else {
		h[0].print()
	}
	if h[0].day == "" {
		t.Errorf("invalid day of from today(): '%v'", h[0].day)
	}
}

func TestSplitText(t *testing.T) {
	wantDay, wantMonth, wantYear := "25", "04", "1976"
	h := NewHaiku(testDate)
	h.splitText(testText)
	if h.day != wantDay {
		t.Errorf("day: want %v, got %v", wantDay, h.day)
	}
	if h.month != wantMonth {
		t.Errorf("month: want %v, got %v", wantMonth, h.month)
	}
	if h.year != wantYear {
		t.Errorf("year: want %v, got %v", wantYear, h.year)
	}
}

func TestCheckDate(t *testing.T) {
	for _, td := range testData {
		err := checkDate(td.input)
		if err != td.outputNegative {
			t.Errorf("wrong result while checking date of '%v': %v != %v", td.input, err, td.outputNegative)
		}
	}
}

func TestReadFile(t *testing.T) {
	testMonth1 := "05"
	testFile1 := "05-07.txt"
	filePath1 := filepath.Join(HAIKU_PATH, testMonth1, testFile1)
	content, err := readFile(filePath1)
	if err != nil {
		t.Errorf("can't read the file: %v", testFile1)
	}
	if len(content) == 0 {
		t.Errorf("empty content!")
	}
}

func TestLoadHaiku(t *testing.T) {
	for _, td := range testData {
		if td.outputNegative == nil {
			h, err := loadHaiku(td.input)
			if err != nil {
				if err != td.outputError {
					t.Errorf("unexpected error while reading by date: %v", err)
				}
			}
			if len(h) != td.outputPositive {
				t.Errorf("unexpected number of text variants: %v != %v", len(h), td.outputPositive)
			}
			for i := 0; i < len(h); i++ {
				h[i].print()
			}
		}
	}
}

func TestFindDate(t *testing.T) {
	wantDay, wantMonth, wantYear := "25", "04", "1976"
	day, month, year := findDate(testText)
	if day != wantDay {
		t.Errorf("day: want %v, got %v", wantDay, day)
	}
	if month != wantMonth {
		t.Errorf("month: want %v, got %v", wantMonth, month)
	}
	if year != wantYear {
		t.Errorf("year: want %v, got %v", wantYear, year)
	}
}

func TestFindAuthor(t *testing.T) {
	wantAuthor := "Ray"
	author := findAuthor(testText)
	if author != wantAuthor {
		t.Errorf("author: want %v, got %v", wantAuthor, author)
	}
}

func TestFindComment(t *testing.T) {
	wantComment := "это коммментарий"
	comment := findComment(testText)
	if comment != wantComment {
		t.Errorf("comment: want %v, got %v", wantComment, comment)
	}
}

type FileVariantVersion struct {
	file    string
	variant string
	version int
}

func TestFindVariant(t *testing.T) {
	fvv := []FileVariantVersion{
		FileVariantVersion{"10-16.txt", "FINAL", 0},
		FileVariantVersion{"10~16_0.txt", "DRAFT", 0},
		FileVariantVersion{"02-05_2.txt", "ALTERNATIVE", 2},
		FileVariantVersion{"05-07_M.txt", "FINAL", 0},
	}
	for i := 0; i < len(fvv); i++ {
		variant, version := findVariant(fvv[i].file)
		if variant != fvv[i].variant {
			t.Errorf("variant for %v: want %v, got %v", fvv[i].file, fvv[i].variant, variant)
		}
		if version != fvv[i].version {
			t.Errorf("version for %v: want %v, got %v", fvv[i].file, fvv[i].version, version)
		}
	}
}
