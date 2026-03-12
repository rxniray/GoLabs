package main

import (
	"fmt"
	"lab_03/calc"
)

func main() {
	fmt.Println("використання пакету calc ")

	// Використовуємо функції зі змінною кількістю аргументів
	fmt.Println("сума:", calc.Sum(1.5, 2.5, 3.0))
	fmt.Println("максимум:", calc.Max(10, -5, 42, 7))
	fmt.Println("мінімум:", calc.Min(10, -5, 42, 7))

	// Обробка помилок при діленні
	res, err := calc.Divide(10, 2)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("результат 10 / 2 =", res)
	}

	// Спробуємо поділити на нуль, щоб побачити помилку
	resZero, errZero := calc.Divide(5, 0)
	if errZero != nil {
		fmt.Println(errZero) // Виведе повідомлення про помилку
	} else {
		fmt.Println("результат:", resZero)
	}
}
