package main

import (
	"errors"
	"fmt"
)

type Calculator interface {
	Sum(nums ...float64) float64
	Max(nums ...float64) float64
	Min(nums ...float64) float64
	Divide(a, b float64) (float64, error)
}
type Calc struct{}

func (c Calc) Sum(nums ...float64) float64 {
	sum := 0.0
	for _, n := range nums {
		sum += n
	}
	return sum
}

func (c Calc) Max(nums ...float64) float64 {
	if len(nums) == 0 {
		return 0
	}
	max := nums[0]
	for _, n := range nums {
		if n > max {
			max = n
		}
	}
	return max
}

func (c Calc) Min(nums ...float64) float64 {
	if len(nums) == 0 {
		return 0
	}
	min := nums[0]
	for _, n := range nums {
		if n < min {
			min = n
		}
	}
	return min
}

func (c Calc) Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("помилка: ділення на нуль неможливе")
	}
	return a / b, nil
}

func main() {
	var myCalc Calculator = Calc{}

	fmt.Println("cума:", myCalc.Sum(4, 5, 6))
	fmt.Println("мінімум:", myCalc.Min(100, 20, 50))

	// Перевірка ділення
	result, err := myCalc.Divide(20, 4)
	if err != nil {
		fmt.Println("помилка:", err)
	} else {
		fmt.Println("20 / 4 =", result)
	}
}
