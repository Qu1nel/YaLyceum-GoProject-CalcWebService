package calculator

import (
	myErr "CalcService/internal/errors"
	"fmt"
	"strconv"
	"strings"
)

var OPERATORS = map[string]struct {
	precedence int
	do         func(float64, float64) float64
}{
	"+": {1, func(x, y float64) float64 { return x + y }},
	"-": {1, func(x, y float64) float64 { return x - y }},
	"*": {2, func(x, y float64) float64 { return x * y }},
	"/": {2, func(x, y float64) float64 { return x / y }},
}

func Calc(expression string) (float64, error) {
	tokens, err := tokenizeExpression(expression)

	if err != nil {
		return 0, err
	}

	fmt.Println(tokens, expression)

	result, err := evaluateTokens(tokens)

	if err != nil {
		return 0, err
	}

	return result, nil
}

func tokenizeExpression(expression string) ([]string, error) {
	var tokens []string
	var buffer strings.Builder

	expression = strings.ReplaceAll(expression, " ", "")

	for idx, char := range expression {
		switch char {
		case ' ':
			continue
		case '+', '-', '*', '/', '(', ')':
			if buffer.Len() > 0 { // Если до этого было число - добавляем в токены, очищаем буфер для последующих токенов
				tokens = append(tokens, buffer.String())
				buffer.Reset()
			}
			if char == ')' {
				n := len(tokens)
				if n == 0 { // Неправильный порядок скобок " ) ("
					return nil, myErr.ErrMissingBracket
				}
				if tokens[n-1] == "(" { // Убираем пустые скобки
					tokens = tokens[:n-1]
					continue
				}
			}
			if char == ')' && len(tokens) != 0 && tokens[len(tokens)-1] == "(" {
				tokens = tokens[:len(tokens)-1]
				continue
			}
			tokens = append(tokens, string(char))
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.':
			buffer.WriteRune(char)
		default:
			return nil, fmt.Errorf("неккоректный символ: \"%c\", на позиции %d (%s)", rune(char), idx, expression)
		}
	}

	if buffer.Len() > 0 { // Запись последнего идентификатора в токены (например числа)
		tokens = append(tokens, buffer.String())
	}

	if len(tokens) == 0 {
		return nil, myErr.ErrEmptyExpression
	}

	return tokens, nil
}

func evaluateTokens(tokens []string) (float64, error) {
	var stack []float64
	var operators []string

	for _, token := range tokens {
		// подразумеваем, что проверяемый токен - число
		if num, err := strconv.ParseFloat(token, 64); err == nil { // если не было ошибок, то это число - добавляем в стек и на следующую итерацию
			stack = append(stack, num)
			continue
		}

		// проверяем не число [ + - * / ) ( ]
		switch token {
		case "(":
			operators = append(operators, token)
		case ")":
			for len(operators) > 0 && operators[len(operators)-1] != "(" {
				last_operator := operators[len(operators)-1]
				operators = operators[:len(operators)-1]
				if err := applyOperator(&stack, last_operator, OPERATORS[last_operator].do); err != nil {
					return 0, err
				}
			}
			if len(operators) == 0 {
				return 0, myErr.ErrMissingBracket
			}
			operators = operators[:len(operators)-1] // Убираем '('
		default:
			operator := token

			for len(operators) > 0 && OPERATORS[operators[len(operators)-1]].precedence >= OPERATORS[operator].precedence {
				last_operator := operators[len(operators)-1]
				operators = operators[:len(operators)-1]
				if err := applyOperator(&stack, last_operator, OPERATORS[last_operator].do); err != nil {
					return 0, err
				}
			}
			operators = append(operators, token)
		}
	}

	for len(operators) > 0 {
		last_operator := operators[len(operators)-1]
		operators = operators[:len(operators)-1]
		if err := applyOperator(&stack, last_operator, OPERATORS[last_operator].do); err != nil {
			return 0, err
		}
	}

	if len(stack) != 1 {
		return 0, myErr.ErrInvalidExpression
	}

	return stack[0], nil
}

func applyOperator(stack *[]float64, operator string, fn func(float64, float64) float64) error {
	if len(*stack) < 2 {
		return myErr.ErrNotEnogthOperand
	}

	b := (*stack)[len(*stack)-1]
	*stack = (*stack)[:len(*stack)-1]
	a := (*stack)[len(*stack)-1]
	*stack = (*stack)[:len(*stack)-1]

	if operator == "/" && b == 0 {
		return myErr.ErrDivisionByZero
	}

	*stack = append(*stack, fn(a, b))

	return nil
}
