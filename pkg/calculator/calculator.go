package calculator

import (
	"strconv"
	"strings"
)

func getPrecedence(operator string) int {
	switch operator {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	default:
		return 0
	}
}

func Calc(expression string) (float64, error) {
	if strings.TrimSpace(expression) == "" {
		return 0, ErrEmptyExpression
	}

	tokens := tokenize(expression)
	postfixed, err := toPostfixNotation(tokens)

	if err != nil {
		return 0, err
	}

	calculated, err := calculate(postfixed)

	if err != nil {
		return 0, err
	}

	return calculated, nil
}

func tokenize(expr string) []string {
	var tokens []string
	var nums string

	for _, r := range expr {
		strRune := string(r)
		if strings.ContainsAny(strRune, "0123456789.") {
			nums += strRune
		} else {
			if nums != "" {
				tokens = append(tokens, nums)
				nums = ""
			}
			if strRune != " " {
				tokens = append(tokens, strRune)
			}
		}
	}

	if nums != "" {
		tokens = append(tokens, nums)
	}

	return tokens
}

func toPostfixNotation(tokens []string) ([]string, error) {
	var output []string
	var operators []string

	for _, token := range tokens {
		switch token {
		case "+", "-", "*", "/":
			for len(operators) > 0 && getPrecedence(operators[len(operators)-1]) >= getPrecedence(token) {
				output = append(output, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			operators = append(operators, token)
		case ")":
			for len(operators) > 0 && operators[len(operators)-1] != "(" {
				output = append(output, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			if len(operators) == 0 {
				return nil, ErrMissLeftParanthesis
			}
			operators = operators[:len(operators)-1]
		case "(":
			operators = append(operators, token)
		default:
			output = append(output, token)
		}
	}

	for len(operators) > 0 {
		if operators[len(operators)-1] == "(" {
			return nil, ErrMissRightParanthesis
		}
		output = append(output, operators[len(operators)-1])
		operators = operators[:len(operators)-1]
	}

	return output, nil
}

func calculate(tokens []string) (float64, error) {
	var stack []float64
	for _, token := range tokens {
		if num, err := strconv.ParseFloat(token, 64); err == nil {
			stack = append(stack, num)
		} else {
			if len(stack) < 2 {
				return 0, ErrNotEnogthOperand
			}
			b, a := stack[len(stack)-1], stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			switch token {
			case "-":
				stack = append(stack, a-b)
			case "+":
				stack = append(stack, a+b)
			case "*":
				stack = append(stack, a*b)
			case "/":
				if b == 0 {
					return 0, ErrDivisionByZero
				}
				stack = append(stack, a/b)
			default:
				return 0, ErrInvalidExpression
			}
		}
	}

	if len(stack) != 1 {
		return 0, ErrInvalidExpression
	}

	return stack[0], nil
}
