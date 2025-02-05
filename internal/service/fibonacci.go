package service

import (
	"fmt"
	"strconv"
)

func RespondFibonacci(message string) string {
	number, err := strconv.Atoi(message)
	if err != nil {
		return fmt.Sprintf("%s is not a number", message)
	}
	return strconv.Itoa(CalculateFibonacci(number))
}

func CalculateFibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return CalculateFibonacci(n-1) + CalculateFibonacci(n-2)
}
