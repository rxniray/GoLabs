package main

import (
	"fmt"
)

func generate() <-chan int {
	out := make(chan int, 10)
	go func() {
		for i := 1; i <= 100; i++ {
			out <- i
		}
		close(out)
	}()
	return out
}

func filterEven(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			if n%2 == 0 {
				out <- n
			}
		}
		close(out)
	}()
	return out
}

func square(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func sum(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		total := 0
		for n := range in {
			total += n
		}
		out <- total
		close(out)
	}()
	return out
}

func main() {
	fmt.Println("запуск конвеєра pipeline...")

	genCh := generate()
	filteredCh := filterEven(genCh)
	squaredCh := square(filteredCh)
	sumCh := sum(squaredCh)

	finalResult := <-sumCh

	fmt.Printf("фінальна сума квадратів парних чисел: %d\n", finalResult)
}
