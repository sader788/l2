package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Config struct {
	fields     []int
	det        string
	isSeparate bool
}

type Cutter struct {
	lines []string
	cfg   Config
}

func getFieldsIdx(str string) ([]int, error) {
	splited := strings.Split(str, ",")

	fields := make([]int, 0, len(splited))
	for _, f := range splited {
		i, err := strconv.Atoi(f)

		if err != nil {
			return []int{}, err
		}

		fields = append(fields, i-1)
	}

	sort.Ints(fields)

	return fields, nil
}

func (c *Cutter) Print() {

	for _, line := range c.lines {
		if !strings.Contains(line, c.cfg.det) && c.cfg.isSeparate {
			continue
		}

		splited := strings.Split(line, c.cfg.det)

		for _, i := range c.cfg.fields {
			if i >= len(splited) {
				break
			}
			fmt.Print(splited[i], c.cfg.det)
		}
		fmt.Println()
	}
}

func NewCutter(lines []string, cfg Config) *Cutter {
	return &Cutter{lines, cfg}
}

func ReadLines() ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func ParseConfig() (Config, error) {
	c := Config{}

	var fieldList string
	flag.StringVar(&fieldList, "f", "", "fields")
	flag.StringVar(&c.det, "d", "	", "delimeter")
	isSeparate := flag.Bool("s", true, "only string with delimeter")

	flag.Parse()

	c.isSeparate = *isSeparate

	var err error
	c.fields, err = parseFieldList(fieldList)
	if err != nil {
		log.Fatal("Bad field list")
	}

	return c, nil
}

func getSplitedNums(s string) ([]int, error) {
	str := strings.Split(s, "-")
	if len(str) != 3 {
		return nil, errors.New("Bad fields")
	}

	first, err := strconv.Atoi(str[0])
	if err != nil {
		return nil, errors.New("Bad fields")
	}
	second, err := strconv.Atoi(str[2])
	if err != nil {
		return nil, errors.New("Bad fields")
	}

	nums := []int{first, second}

	return nums, nil
}

func getFieldsMap(split []string) (map[int]struct{}, error) {
	mapFields := make(map[int]struct{})

	for _, s := range split {
		if strings.Contains(s, "-") {
			nums, err := getSplitedNums(s)
			if err != nil {
				return nil, err
			}

			for _, num := range nums {
				mapFields[num] = struct{}{}
			}
		}

		atoi, err := strconv.Atoi(s)
		if err != nil {
			return nil, errors.New("Bad fields")
		}
		mapFields[atoi] = struct{}{}
	}

	return mapFields, nil
}

func parseFieldList(str string) ([]int, error) {
	list := strings.TrimSpace(str)
	list = strings.Replace(str, " ", "", -1)

	split := strings.Split(list, ",")

	mapFields, err := getFieldsMap(split)
	if err != nil {
		return nil, err
	}

	fields := []int{}
	for k, _ := range mapFields {
		fields = append(fields, k)
	}

	return fields, nil
}

func main() {
	cfg, err := ParseConfig()
	if err != nil {
		log.Fatal(err)
	}

	lines, err := ReadLines()
	if err != nil {
		log.Fatal(err)
	}

	cutter := NewCutter(lines, cfg)
	cutter.Print()
}
