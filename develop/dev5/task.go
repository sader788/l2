package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

type Config struct {
	after    int
	before   int
	ctx      int
	isCnt    bool
	isIgnore bool
	isInvert bool
	isFixed  bool
	isNum    bool
	fileName string
	expr     string
}

type Line struct {
	num     int
	str     string
	idxs    [][]int
	isMatch bool
}

func NewLine(num int, str string) *Line {
	return &Line{num, str, [][]int{}, false}
}

type Parser struct {
	lines        []*Line
	matchedLines int
	len          int
	cfg          Config
}

func (p *Parser) Parse(expr string) {
	if p.cfg.isIgnore {
		expr = strings.ToLower(expr)
	}

	rExp := regexp.MustCompile(expr)

	for _, line := range p.lines {
		str := line.str

		if p.cfg.isIgnore {
			str = strings.ToLower(str)
		}
		idxs := rExp.FindAllStringIndex(str, -1)

		if len(idxs) > 0 {
			if p.cfg.isFixed && str[idxs[0][0]:idxs[0][1]] != expr {
				continue
			}
			line.isMatch = true
			line.idxs = idxs
			p.matchedLines++
		}
	}
}

func coloredIndex(str string, idxs [][]int) (string, error) {
	last := 0
	sb := strings.Builder{}

	for _, idx := range idxs {
		if len(idx) != 2 {
			return "", errors.New("There must be a pair of indexes pointing to the start and end")
		}

		sb.WriteString(str[last:idx[0]])

		sb.WriteString("\033[31m") //red color
		sb.WriteString(str[idx[0]:idx[1]])
		sb.WriteString("\033[0m") //reset color

		last = idx[1]
	}

	if last < len(str) {
		sb.WriteString(str[last:len(str)])
	}

	return sb.String(), nil
}

func numberedLine(str string, num int) string {
	return fmt.Sprintf("\033[34m%d:\033[0m%s", num, str)
}

func (p *Parser) getCount() int {
	if p.cfg.isInvert {
		return p.len - p.matchedLines
	}

	return p.matchedLines
}

func (p *Parser) getLinesByIndex(lineIdxs []int) ([]string, error) {
	matchedLines := []string{}

	for _, i := range lineIdxs {
		strLine := p.lines[i].str

		if p.lines[i].isMatch {
			coloredLine, err := coloredIndex(strLine, p.lines[i].idxs)
			if err != nil {
				return nil, err
			}
			strLine = coloredLine
		}
		if p.cfg.isNum {
			strLine = numberedLine(strLine, p.lines[i].num)
		}

		matchedLines = append(matchedLines, strLine)
	}

	return matchedLines, nil
}

func (p *Parser) getLinesIndex() []int {
	linesMatched := map[int]struct{}{}

	after, before := p.cfg.after, p.cfg.before

	if p.cfg.ctx > after {
		after = p.cfg.ctx
	}
	if p.cfg.ctx > before {
		before = p.cfg.ctx
	}

	for i, line := range p.lines {
		if (p.cfg.isInvert && line.isMatch) || (!p.cfg.isInvert && !line.isMatch) {
			continue
		}
		linesMatched[i] = struct{}{}

		for j := 1; j <= before && i-j >= 0; j++ {
			linesMatched[i-j] = struct{}{}
		}
		for j := 1; j <= after && i+j < p.len; j++ {
			linesMatched[i+j] = struct{}{}
		}
	}

	lineIdxs := make([]int, 0, len(linesMatched))
	for k, _ := range linesMatched {
		lineIdxs = append(lineIdxs, k)
	}

	sort.Ints(lineIdxs)
	return lineIdxs
}

func (p *Parser) Print() {
	if p.cfg.isCnt {
		fmt.Println(p.getCount())
		return
	}

	lines, err := p.getLinesByIndex(p.getLinesIndex())
	if err != nil {
		panic(err.Error())
	}

	for _, line := range lines {
		fmt.Println(line)
	}
}

func ParseConfig() Config {
	c := Config{}

	flag.IntVar(&c.after, "A", 0, "Print NUM lines of trailing context after matching lines.")
	flag.IntVar(&c.before, "B", 0, "Print NUM lines of leading context before matching lines.")
	flag.IntVar(&c.ctx, "C", 0, "Print NUM lines of output context.")
	flag.StringVar(&c.expr, "r", "", "Regular expression.")
	flag.StringVar(&c.fileName, "f", "", " Obtain patterns from FILE, one per line.")

	isCnt := flag.Bool("c", false, "Suppress normal output; instead print a count of matching lines for each input file.")
	isIgnore := flag.Bool("i", false, "Ignore case distinctions, so that characters that differ only in case match each other.")
	isInvert := flag.Bool("v", false, "Invert the sense of matching, to select non-matching lines.")
	isFixed := flag.Bool("F", false, "Interpret  PATTERN as a list of fixed strings (instead of regular expressions), separated by newlines, any of which is to be matched.")
	isNum := flag.Bool("n", false, "Prefix each line of output with the 1-based line number within its input file.")

	flag.Parse()

	c.isCnt = *isCnt
	c.isIgnore = *isIgnore
	c.isInvert = *isInvert
	c.isFixed = *isFixed
	c.isNum = *isNum

	return c
}

func NewParser(strs []string, cfg Config) *Parser {
	lines := []*Line{}

	for i, str := range strs {
		lines = append(lines, NewLine(i+1, str))
	}

	return &Parser{lines, 0, len(lines), cfg}
}

func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func main() {
	cfg := ParseConfig()

	lines, err := ReadLines(cfg.fileName)
	if err != nil {
		return
	}

	parser := NewParser(lines, cfg)

	parser.Parse(cfg.expr)

	parser.Print()
}
