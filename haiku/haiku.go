package haiku

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	HAIKU_PATH = "../year"
	DRAFT      = iota
	MIKE
	RAY
	FINAL
)

type Haiku struct {
	day     string
	month   string
	year    string
	text    string
	author  string
	comment string
	variant int // DRAFT .. FINAL
}

var (
	EmptyDateError    = errors.New("empty date")
	InvalidDateError  = errors.New("invalid date")
	BadDelimiterError = errors.New("bad date delimiter")
	TextMissingError  = errors.New("haiku text missing")
)

func NewHaiku(date string) *Haiku {
	ymd := strings.Split(date, "-")
	h := Haiku{day: fmt.Sprintf("%02s", ymd[2]), month: fmt.Sprintf("%02s", ymd[1])}
	return &h
}

func (h Haiku) print() {
	fmt.Printf("%s\tDate: %s.%s.%s\n\tAuthor: %s\n\tComment: %s\n-------------------------\n",
		h.text, h.day, h.month, h.year, h.author, h.comment)
}

func (h *Haiku) splitText(content string) {
	lines := strings.Split(content, "\n")
	text := ""
	for i := 0; i < len(lines); i++ {
		if strings.Contains(lines[i], "{") {
			break
		}
		text += lines[i] + "\n"
	}
	h.text = text
	h.day, h.month, h.year = findDate(content)
	h.author = findAuthor(content)
	h.comment = findComment(content)
}

func readHaiku(date string) (today []Haiku, err error) {
	today = []Haiku{}
	err = checkDate(date)
	if err != nil {
		return today, err
	}
	h := NewHaiku(date)

	variants := []string{
		"%s-%s.txt", "%s-%s_M.txt", "%s-%s_R.txt",
		"%s-%s_0.txt", "%s-%s_1.txt", "%s-%s_2.txt", "%s-%s_3.txt",
		"%s~%s_0.txt", "%s~%s_1.txt", "%s~%s_2.txt", "%s~%s_3.txt",
	}
	for i := 0; i < len(variants); i++ {
		err = nil
		file := fmt.Sprintf(variants[i], h.month, h.day)
		filePath := filepath.Join(HAIKU_PATH, h.month, file)
		t, err := readFile(filePath)
		if err != nil {
			continue
		}
		h.text = t
		h.variant = FINAL
		h.splitText(t)
		today = append(today, *h)
	}
	return today, nil
}

func readFile(filePath string) (content string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", TextMissingError
		}
		return "", err
	}
	defer file.Close()
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func checkDate(date string) (err error) {
	if date == "" {
		return EmptyDateError
	}
	if !strings.Contains(date, "-") {
		return BadDelimiterError
	}
	if strings.Count(date, "-") != 2 {
		return BadDelimiterError
	}
	ymd := strings.Split(date, "-")
	if len(ymd) < 3 {
		return InvalidDateError
	}
	if ymd[0] == "" || ymd[1] == "" || ymd[2] == "" {
		return EmptyDateError
	}
	y, err := strconv.Atoi(ymd[0])
	if err != nil {
		return InvalidDateError
	}
	m, err := strconv.Atoi(ymd[1])
	if err != nil {
		return InvalidDateError
	}
	d, err := strconv.Atoi(ymd[2])
	if err != nil {
		return InvalidDateError
	}
	if d <= 0 || m <= 0 || y <= 0 {
		return InvalidDateError
	}
	if d > 31 || m > 12 {
		return InvalidDateError
	}
	return nil
}

func findDate(content string) (day, month, year string) {
	re := regexp.MustCompile(`{(\d+)[.](\d+)[.](\d+)}`)
	matches := re.FindStringSubmatch(content)
	if len(matches) == 4 {
		day = matches[1]
		month = matches[2]
		year = matches[3]
	}
	return day, month, year
}

func findAuthor(content string) (author string) {
	re := regexp.MustCompile(`\[([^[]+)\]`)
	matches := re.FindStringSubmatch(content)
	if len(matches) == 2 {
		author = matches[1]
	}
	if author == "" {
		author = "Mike"
	}
	return author
}

func findComment(content string) (comment string) {
	re := regexp.MustCompile(`<([^>]+)>`)
	matches := re.FindStringSubmatch(content)
	if len(matches) == 2 {
		comment = matches[1]
	}
	return comment
}
