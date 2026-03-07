package main

import (
	"fmt"
	"math/rand"
)

func main() {
	var y float64
	x := rand.Intn(150)

	if x > 100 {
		y = float64((2*x + 3) * (4 - x))
	} else {
		y = float64(x-4) / float64(x)
	}

	/*
		switch {
		case x > 100:
			y = (2*xFloat + 3) * (4 - xFloat)
		default:
			if x == 0 {
				fmt.Println("Кернес тільки навпаки")
			}
			y = (xFloat - 4) / xFloat
		}
	*/
	fmt.Printf("x = %d\n", x)
	fmt.Printf("Результат y = %f\n", y)
}
