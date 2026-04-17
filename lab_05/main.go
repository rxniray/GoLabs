package main

import (
	"fmt"
	"sync"
)

func main() {
	var counter int
	var mu sync.Mutex
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
					mu.Lock()
					counter++
					mu.Unlock() 
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
					mu.Lock()
					counter--
					mu.Unlock()
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

	fmt.Printf("Фінальне значення counter (Mutex): %d\n", counter)
}