package utils

import "fmt"

func EnterExpression(prompt string) (string, error) {
	var expression string
	fmt.Print(prompt + " ")
	_, err := fmt.Scan(&expression)
	return expression, err
}

func PrintResultExpression(log_message string, result interface{}) (int, error) {
	n, err := fmt.Printf(log_message+": %v\n", result)
	return n, err
}
