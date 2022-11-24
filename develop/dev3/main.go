package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Config struct {
	lines     [][]string
	column    int
	isNum     bool
	isReverse bool
	isUniqe   bool
	fileName  string
}

func (c Config) Len() int {
	return len(c.lines)
}

func (c Config) Less(i, j int) bool {
	if c.isNum {
		a, _ := strconv.Atoi(c.lines[i][c.column])
		b, _ := strconv.Atoi(c.lines[j][c.column])

		if c.isReverse {
			return a >= b
		}
		return a < b
	}

	if c.isReverse {
		return c.lines[i][c.column] >= c.lines[j][c.column]
	}

	return c.lines[i][c.column] < c.lines[j][c.column]
}

func (c Config) Swap(i, j int) {
	c.lines[i], c.lines[j] = c.lines[j], c.lines[i]
}

func (c Config) GetLines() []string {
	res := []string{}

	uniqe := map[string]struct{}{}

	sb := strings.Builder{}
	for _, line := range c.lines {
		for _, str := range line {
			sb.WriteString(str)
			sb.WriteRune(' ')
		}
		if _, existed := uniqe[sb.String()]; !existed || !c.isUniqe {
			res = append(res, sb.String())
			uniqe[sb.String()] = struct{}{}
		}

		sb.Reset()
	}

	return res
}

func (c *Config) ReadLines() error {
	strs, err := readFile(c.fileName)
	if err != nil {
		return err
	}

	for _, str := range strs {
		c.lines = append(c.lines, strings.Split(str, " "))
	}

	return nil
}

func NewConfig() Config {
	c := Config{}

	flag.StringVar(&c.fileName, "f", "", "file name")
	flag.IntVar(&c.column, "k", 0, "column for sort")
	isNum := flag.Bool("n", false, "sort by num")
	isReverse := flag.Bool("r", false, "reverse sort")
	isUniqe := flag.Bool("u", false, "uniqe")

	flag.Parse()

	c.isNum = *isNum
	c.isReverse = *isReverse
	c.isUniqe = *isUniqe

	return c
}

func readFile(path string) ([]string, error) {
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

func WriteLines(path string, lines []string) { //TODO args
	f, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
	}

	for i, line := range lines {
		if i != 0 {
			f.WriteString("\n")
		}

		f.WriteString(line)
	}

	defer f.Close()
}

func main() {
	config := NewConfig()

	config.ReadLines()
	sort.Sort(config)

	WriteLines("sorted.txt", config.GetLines())
}
