package haiku

import (
	"embed"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	HAIKU_PATH = "year"

	DRAFT = iota
	ALTERNATIVE
	FINAL
)

type Haiku struct { // 俳句
	day     string
	month   string
	year    string
	text    string // 句
	author  string
	comment string
	variant string // DRAFT .. FINAL
	version int
}

var (
	EmptyDateError    = errors.New("empty date")
	InvalidDateError  = errors.New("invalid date")
	BadDelimiterError = errors.New("bad date delimiter")
	TextMissingError  = errors.New("haiku text missing")
)

var variants = []string{
	"%s-%s.txt", "%s-%s_M.txt", "%s-%s_R.txt",
	"%s-%s_0.txt", "%s-%s_1.txt", "%s-%s_2.txt", "%s-%s_3.txt",
	"%s~%s_0.txt", "%s~%s_1.txt", "%s~%s_2.txt", "%s~%s_3.txt",
}

//go:embed year
var haikuDir embed.FS

func iota2string(i int) (s string) {
	switch i {
	case DRAFT:
		s = "DRAFT"
	case ALTERNATIVE:
		s = "ALTERNATIVE"
	case FINAL:
		s = "FINAL"
	}
	return s
}

func Today() (today []Haiku) {
	kyou := time.Now().Format("2006-01-02") // 今日
	today, err := loadHaiku(kyou)
	if err != nil || len(today) == 0 {
		today, err = pretext()
		if err != nil {
			log.Printf("Today: %v", err)
		}
	}
	return today
}

func NewHaiku(date string) *Haiku {
	ymd := strings.Split(date, "-")
	h := Haiku{day: fmt.Sprintf("%02s", ymd[2]), month: fmt.Sprintf("%02s", ymd[1])}
	return &h
}

func (h Haiku) Verse() string {
	return h.text
}

func (h Haiku) print() {
	fmt.Printf("%s\tDate: %s.%s.%s\n\tAuthor: %s\n\tComment: %s\n\tVariant: %v\n\tVersion: %v\n-------------------------\n",
		h.text, h.day, h.month, h.year, h.author, h.comment, h.variant, h.version)
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

func pretext() (unwritten []Haiku, err error) {
	unwritten = []Haiku{}
	h, err := readHaiku("0000-00-00", filepath.Join(HAIKU_PATH, "00"), "00-00.txt")
	if err != nil {
		log.Printf("pretext: %v", err)
	}
	if err == nil {
		unwritten = append(unwritten, *h)
	}
	return unwritten, err
}

func loadHaiku(date string) (list []Haiku, err error) {
	list = []Haiku{}
	err = checkDate(date)
	if err != nil {
		return list, err
	}

	for i := 0; i < len(variants); i++ {
		err = nil
		h := NewHaiku(date)
		fileName := fmt.Sprintf(variants[i], h.month, h.day)
		filePath := filepath.Join(HAIKU_PATH, h.month)
		h, err = readHaiku(date, filePath, fileName)
		if err != nil {
			continue
		}
		list = append(list, *h)
	}
	return list, nil
}

func readHaiku(date, filePath, fileName string) (h *Haiku, err error) {
	h = NewHaiku(date)
	t, err := readFile(filepath.Join(filePath, fileName))
	if err != nil {
		return nil, err
	}
	h.text = t
	h.variant, h.version = findVariant(fileName)
	h.splitText(t)
	return h, nil
}

func readFile(filePath string) (content string, err error) {
	file, err := haikuDir.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", TextMissingError
		}
		return "", err
	}
	defer file.Close()
	data, err := haikuDir.ReadFile(filePath)
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

func findVariant(fileName string) (variant string, version int) {
	dateAndVersion := strings.Split(fileName, "_")
	if len(dateAndVersion) == 2 {
		n := strings.Replace(dateAndVersion[1], ".txt", "", -1)
		v, err := strconv.Atoi(n)
		if err == nil {
			version = v
		}
	}
	if strings.Contains(fileName, "-") {
		if version == 0 {
			variant = iota2string(FINAL)
		} else {
			variant = iota2string(ALTERNATIVE)
		}
	}
	if strings.Contains(fileName, "~") {
		variant = iota2string(DRAFT)
	}
	return variant, version
}
