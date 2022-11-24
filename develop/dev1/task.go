package main

import (
	"fmt"
	"github.com/beevik/ntp"
	"time"
)

func GetTime(host string) (*time.Time, error) {
	t, err := ntp.Time(host)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func main() {
	timeNTP, err := GetTime("host")
	if err != nil {
		panic(err.Error())
	}

	timePC := time.Now()

	fmt.Println(timeNTP)
	fmt.Println(timePC)
}
