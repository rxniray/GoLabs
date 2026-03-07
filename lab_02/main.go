package main

import "fmt"

func main() {

	var a = [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	b := []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}

	result := make([]int, 10)

	for i := 0; i < 10; i++ {
		result[i] = a[i] + b[i]
	}

	fmt.Println("Масив a:", a)
	fmt.Println("Слайс b:", b)
	fmt.Println("Результат (сума):", result)
}
