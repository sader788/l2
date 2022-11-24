package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func Extract(str string) (string, error) {
	if len(str) == 0 {
		return "", nil
	}

	runes := []rune(str)

	sb := strings.Builder{}

	for i := 0; i < len(str); i++ {
		if isDigit(runes[i]) {
			if i < 1 {
				return "", errors.New("Wrong string")
			}

			j := i
			for isDigit(runes[j]) {
				j++
				if j == len(str) {
					break
				}
			}

			cnt, _ := strconv.Atoi(string(runes[i:j]))
			for cnt > 1 {
				sb.WriteRune(runes[i-1])
				cnt--
			}
		} else {
			sb.WriteRune(runes[i])
		}
	}
	return sb.String(), nil
}

func main() {

	res, err := Extract("a11bbbc5")
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(res)
}
