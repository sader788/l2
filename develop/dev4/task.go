package main

import (
	"fmt"
	"sort"
)

type AnagramSet struct {
	annagrams []string
	letters   string
	letLen    int //letters len
}

func (as *AnagramSet) IsAnagram(str string) bool {
	if len(str) != as.letLen {
		return false
	}

	strRunes := []rune(str)
	sort.Slice(strRunes, func(i, j int) bool {
		return strRunes[i] < strRunes[j]
	})
	strLetters := string(strRunes)

	if strLetters != as.letters {
		return false
	}

	return true
}

func (as *AnagramSet) AddAnagram(str string) {
	as.annagrams = append(as.annagrams, str)
}

func (as *AnagramSet) GetAnagramList() (string, []string) {
	sort.Strings(as.annagrams)
	return as.annagrams[0], as.annagrams
}

func NewAnagramSet(str string) *AnagramSet {
	annagrams := []string{str}

	strRunes := []rune(str)
	sort.Slice(strRunes, func(i, j int) bool {
		return strRunes[i] < strRunes[j]
	})
	letters := string(strRunes)

	return &AnagramSet{annagrams, letters, len(letters)}
}

func AnagramAppend(sets []*AnagramSet, str string) []*AnagramSet {
	for _, set := range sets {
		if set.IsAnagram(str) {
			set.AddAnagram(str)
			return sets
		}
	}

	return append(sets, NewAnagramSet(str))
}

func NewAnagramMap(strs []string) *map[string][]string {
	sets := []*AnagramSet{}

	for _, str := range strs {
		sets = AnagramAppend(sets, str)
	}

	res := map[string][]string{}
	for _, set := range sets {
		s, l := set.GetAnagramList()
		res[s] = l
	}

	return &res
}

func main() {
	anagramMap := NewAnagramMap([]string{`пятак`, `пятка`, `тяпка`, `листок`, `слиток`, `столик`, `рофлинка`, `калинроф`})

	fmt.Println(anagramMap)
}
