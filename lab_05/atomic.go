package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var counter int64
	var wg sync.WaitGroup
	evenCh := make(chan int)
	oddCh := make(chan int)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case val, ok := <-evenCh:
				if !ok {
					return
				}
				if val%3 == 0 {
					atomic.AddInt64(&counter, 1)
				}
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case val, ok := <-oddCh:
				if !ok {
					return
				}
				if val%33 == 0 {
					atomic.AddInt64(&counter, -1)
				}
			}
		}
	}()

	for i := 1; i <= 1000; i++ {
		if i%2 == 0 {
			evenCh <- i
		} else {
			oddCh <- i
		}
	}

	close(evenCh)
	close(oddCh)
	wg.Wait()

	finalCounter := atomic.LoadInt64(&counter)
	fmt.Printf("Фінальне значення counter (Atomic): %d\n", finalCounter)
}