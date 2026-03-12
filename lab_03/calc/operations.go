package calc

import (
	"errors"
	"fmt"
)

func init() {
	fmt.Println("[init++]")
}

func Sum(nums ...float64) float64 {
	sum := 0.0
	for _, n := range nums {
		sum += n
	}
	return sum
}

func Max(nums ...float64) float64 {
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

func Min(nums ...float64) float64 {
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

func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("помилка: ділення на нуль неможливе")
	}
	return a / b, nil
}
