package main

import (
	"fmt"
	"sync"
	"time"
)

func or(channels ...<-chan interface{}) <-chan interface{} {
	single := make(chan interface{})

	wg := &sync.WaitGroup{}

	for _, channel := range channels {
		wg.Add(1)
		go func(wg *sync.WaitGroup, c <-chan interface{}) {
			done := <-c
			single <- done
			wg.Done()
		}(wg, channel)
	}

	go func(wg *sync.WaitGroup) {
		wg.Wait()
		close(single)
	}(wg)

	return single
}

func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(2*time.Minute),
		sig(2*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("fone after %v", time.Since(start))

}
