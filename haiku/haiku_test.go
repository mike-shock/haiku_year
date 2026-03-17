package haiku

import (
	//	"fmt"
	"testing"
)

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
)

func TestCheckDate(t *testing.T) {
	for _, td := range testData {
		err := checkDate(td.input)
		if err != td.outputNegative {
			t.Errorf("wrong result while checking date of '%v': %v != %v", td.input, err, td.outputNegative)
		}
	}
}

func TestReadHaiku(t *testing.T) {
	for _, td := range testData {
		if td.outputNegative == nil {
			h, err := readHaiku(td.input)
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
