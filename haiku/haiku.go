package haiku

import (
	"embed"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	HAIKU_PATH = "year"

	DRAFT       = iota // 下書き
	ALTERNATIVE        // 代替
	FINAL              // 完了
)

type Haiku struct { // 俳句
	day     string // 日
	month   string // 月
	year    string // 年
	text    string // 句
	author  string // 詩人, 著
	comment string // 言い草
	variant string // 変異体: DRAFT .. FINAL
	version int    // 稿
}
type Verses []Haiku // 詩

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

func iota2string(i int) (s string) { // 号を文に化
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

func Today() (list Verses) { // 今日
	kyou := time.Now().Format("2006-01-02")
	list, _ = loadHaiku(kyou) // ThisDay(kyou)
	//if err != nil {
	//log.Printf("Today(): %v", err)
	//}
	for _, h := range list {
		h.print()
	}
	return list
}

func ThisDay(date string) (haiku Verses, err error) { // この日
	haiku, err = loadHaiku(date)
	if err != nil || len(haiku) == 0 {
		haiku, err = loadHaiku("0000-00-00") // substitute = 代わり
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	rng.Shuffle(len(haiku), func(i, j int) { haiku[i], haiku[j] = haiku[j], haiku[i] })
	fmt.Printf("ThisDay(): %v", haiku)
	return haiku, err
}

func IsHaiku(date string) (ok bool) { // 俳句ですか
	if err := checkDate(date); err != nil {
		return false
	}
	filePath := date2path(date)
	if err := checkFile(filePath); err == nil {
		ok = true
	}
	return ok
}

func NewHaiku(date string) *Haiku { // 新しい俳句
	ymd := strings.Split(date, "-")
	h := Haiku{day: fmt.Sprintf("%02s", ymd[2]), month: fmt.Sprintf("%02s", ymd[1])}
	return &h
}

func (h Haiku) Date() string { // 日付
	return fmt.Sprintf("%04s-%02s-%02s", h.year, h.month, h.day)
}

func (h Haiku) Verse() string { // 詩
	return h.text
}

func (h Haiku) Author() string { // 詩人
	return h.author
}

func (h Haiku) Comment() string { // 言い草
	return h.comment
}

func (h Haiku) print() { // 刷る
	fmt.Printf("%s\tDate: %s.%s.%s\n\tAuthor: %s\n\tComment: %s\n\tVariant: %v\n\tVersion: %v\n-------------------------\n",
		h.text, h.day, h.month, h.year, h.author, h.comment, h.variant, h.version)
}

func (h *Haiku) splitText(content string) { // 本書を分
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

func (v Verses) sort() {
	sort.Slice(v, func(i, j int) bool {
		return v[i].variant < v[j].variant
	})
}

func loadHaiku(date string) (list Verses, err error) { // 俳句を引
	list = Verses{}
	err = checkDate(date)
	if err != nil {
		return list, err
	}
	for i := 0; i < len(variants); i++ {
		err = nil
		h := NewHaiku(date)
		fileName := fmt.Sprintf(variants[i], h.month, h.day)
		filePathAndName := fullFileName(h.month, fileName)
		h, err = readHaiku(date, filePathAndName)
		if err != nil {
			continue
		}
		list = append(list, *h)
	}
	return list, nil
}

func readHaiku(date, fileName string) (h *Haiku, err error) { // 俳句を読む
	h = NewHaiku(date)
	t, err := readFile(fileName)
	if err != nil {
		return nil, err
	}
	h.text = t
	h.variant, h.version = findVariant(fileName)
	h.splitText(t)
	return h, nil
}

func readFile(fileName string) (content string, err error) { // ファイルを読む
	if err = checkFile(fileName); err != nil {
		return "", err
	}
	data, err := haikuDir.ReadFile(fileName)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func checkFile(fileName string) (err error) { // ファイルを質す
	file, err := haikuDir.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	return err
}

func date2path(date string) (filePath string) { // 日付を道に化す
	ymd := strings.Split(date, "-")
	fn := fmt.Sprintf("%02s-%02s.txt", ymd[1], ymd[2])
	filePath = fullFileName(ymd[1], fn)
	return filePath
}

func fullFileName(month, file string) (filePathAndName string) { //
	filePathAndName = fmt.Sprintf("%s/%s/%s", HAIKU_PATH, month, file)
	return filePathAndName
}

func checkDate(date string) (err error) { // 日付を質す
	if date == "0000-00-00" { // special "default" date
		return nil
	}
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

func findDate(content string) (day, month, year string) { // 日付を探す
	re := regexp.MustCompile(`{(\d+)[.](\d+)[.](\d+)}`) // {DD.MM.YYYY}
	matches := re.FindStringSubmatch(content)
	if len(matches) == 4 {
		day = matches[1]
		month = matches[2]
		year = matches[3]
	}
	return day, month, year
}

func findAuthor(content string) (author string) { // 詩人を探す
	re := regexp.MustCompile(`\[([^[]+)\]`) // [Author]
	matches := re.FindStringSubmatch(content)
	if len(matches) == 2 {
		author = matches[1]
	}
	if author == "" {
		author = "Mike"
	}
	return author
}

func findComment(content string) (comment string) { // 言い草を探す
	re := regexp.MustCompile(`<([^>]+)>`) // <Comment>
	matches := re.FindStringSubmatch(content)
	if len(matches) == 2 {
		comment = matches[1]
	}
	return comment
}

func findVariant(fileName string) (variant string, version int) { //  変異体を探す
	dateAndVersion := strings.Split(fileName, "_")
	if len(dateAndVersion) == 1 {
		return iota2string(FINAL), 0
	}
	if len(dateAndVersion) == 2 {
		n := strings.Replace(dateAndVersion[1], ".txt", "", -1)
		v, err := strconv.Atoi(n)
		if err == nil {
			version = v
		}
	}
	if strings.Contains(fileName, "-") { // MM-DD_V.txt
		variant = iota2string(ALTERNATIVE)
	}
	if strings.Contains(fileName, "~") { // MM~DD_V.txt
		variant = iota2string(DRAFT)
	}
	return variant, version
}
