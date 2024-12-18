package calculator

import (
	stdErr "errors"
	"fmt"
	"testing"

	calculatorErr "CalcService/internal/errors"

	"github.com/stretchr/testify/assert"
)

type CaseTokenize struct {
	nameCase    string
	actual      interface{}
	expected    []string
	expectedMsg string
	err         error
	errMsg      string
}

func TestTokenizeExpression(t *testing.T) {
	t.Run("Empty expression", testEmptyExpression)
	t.Run("Only digits", testDigits)
	t.Run("Only charaters", testChar)
	t.Run("Only brackets", testBrackets)

	assert := assert.New(t)
	caces := []CaseTokenize{
		{
			nameCase: "00",
			actual:   "1 + 1",
			expected: []string{"1", "+", "1"},
		},
		{
			nameCase: "01",
			actual:   "1 + 1 * 2",
			expected: []string{"1", "+", "1", "*", "2"},
		},
		{
			nameCase: "02",
			actual:   "1 + (1 * 2) / 34",
			expected: []string{"1", "+", "(", "1", "*", "2", ")", "/", "34"},
		},
		{
			nameCase: "03",
			actual:   "(((56) + 1) + 2) + 1",
			expected: []string{"(", "(", "(", "56", ")", "+", "1", ")", "+", "2", ")", "+", "1"},
		},
		{
			nameCase: "04",
			actual:   "1 2 3 4 5))))) * /",
			expected: []string{"12345", ")", ")", ")", ")", ")", "*", "/"},
		},
		{
			nameCase: "05",
			actual:   ") ()",
			err:      calculatorErr.ErrMissingBracket,
			errMsg:   "пропущена скобка",
		},
		{
			nameCase: "06",
			actual:   "( 3 + 33             ())",
			expected: []string{"(", "3", "+", "33", ")"},
		},
		{
			nameCase: "07",
			actual:   "() 3 + 3",
			expected: []string{"3", "+", "3"},
		},
		{
			nameCase: "08",
			actual:   "2 + 2 ((3))",
			expected: []string{"2", "+", "2", "(", "(", "3", ")", ")"},
		},
		{
			nameCase: "09",
			actual:   "1+1* 123455678",
			expected: []string{"1", "+", "1", "*", "123455678"},
		},
		{
			nameCase: "10",
			actual:   "1 1 1 1 1 1 1 1",
			expected: []string{"11111111"},
		},
		{
			nameCase: "11",
			actual:   "1 ////// 2",
			expected: []string{"1", "/", "/", "/", "/", "/", "/", "2"},
		},
		{
			nameCase: "12",
			actual:   "(22) * (22 * (31 / 4.29)) + ((1.11 - 3 + ((4512.12 / 29) * 3)) * 501)",
			expected: []string{"(", "22", ")", "*", "(", "22", "*", "(", "31", "/", "4.29", ")", ")", "+", "(", "(", "1.11", "-", "3", "+", "(", "(", "4512.12", "/", "29", ")", "*", "3", ")", ")", "*", "501", ")"},
		},
		{
			nameCase: "13",
			actual:   "(2 + 2) * (2 * 2)",
			expected: []string{"(", "2", "+", "2", ")", "*", "(", "2", "*", "2", ")"},
		},
	}

	for _, tc := range caces {
		t.Run(tc.nameCase, func(t *testing.T) {
			output, err := tokenizeExpression(tc.actual.(string))
			assert.Equal(output, tc.expected, tc.expectedMsg)
			if err != nil {
				assert.Equal(err, tc.err)
				assert.EqualError(err, tc.errMsg)
			}
		})
	}
}

func testEmptyExpression(t *testing.T) {
	assert := assert.New(t)
	cases := []CaseTokenize{
		{
			nameCase:    "Empty",
			actual:      "",
			expected:    nil,
			expectedMsg: "Пустое выражение должно ожидать nil",
			err:         calculatorErr.ErrEmptyExpression,
			errMsg:      "Пустое выражение должно возвращать ошибку",
		},
		{
			nameCase:    "One_space",
			actual:      " ",
			expected:    nil,
			expectedMsg: "Пустое выражение должно ожидать nil",
			err:         calculatorErr.ErrEmptyExpression,
			errMsg:      "Пустое выражение должно возвращать ошибку",
		},
		{
			nameCase:    "Several_space",
			actual:      "          ",
			expected:    nil,
			expectedMsg: "Пустое выражение должно ожидать nil",
			err:         calculatorErr.ErrEmptyExpression,
			errMsg:      "Пустое выражение должно возвращать ошибку",
		},
		{
			nameCase:    "Bracket groups",
			actual:      " ( (())()(()))",
			expected:    nil,
			expectedMsg: "Выражение с пустыми скобками должно возвращать nil",
			err:         calculatorErr.ErrEmptyExpression,
			errMsg:      "Выражение с пустыми скобками должно возвращать ошибку",
		},
	}

	for _, tc := range cases {
		t.Run(tc.nameCase, func(t *testing.T) {
			output, err := tokenizeExpression(tc.actual.(string))
			assert.Equal(output, tc.expected, tc.expectedMsg)
			if assert.Error(err) {
				assert.Equal(err, tc.err)
				assert.EqualError(err, "пустое выражение")
			}
		})
	}
}

func testDigits(t *testing.T) {
	assert := assert.New(t)
	cases := []CaseTokenize{
		{
			nameCase:    "One_digit_1",
			actual:      "        2   ",
			expected:    []string{"2"},
			expectedMsg: "Выражение c одним числом должно возвращать cлaйлc c этим числом",
		},
		{
			nameCase:    "One_digit_2",
			actual:      "2935   ",
			expected:    []string{"2935"},
			expectedMsg: "Выражение c одним числом должно возвращать cлaйлc c этим числом",
		},
		{
			nameCase:    "One_digit_3",
			actual:      "         9894000",
			expected:    []string{"9894000"},
			expectedMsg: "Выражение c одним числом должно возвращать cлaйлc c этим числом",
		},
		{
			nameCase: "One_digit_4",
			actual:   "         3,4",
			expected: nil,
			err:      fmt.Errorf("неккоректный символ: \",\", на позиции 1 (3,4)"),
		},
		{
			nameCase:    "Digits",
			actual:      "  1 2 3 4 5 6    ",
			expected:    []string{"123456"},
			expectedMsg: "Выражение только с числами, разделенными пробелам, должно возвращать слайс с один числом",
		},
		{
			nameCase:    "Floating_digit_1",
			actual:      " 3.1415",
			expected:    []string{"3.1415"},
			expectedMsg: "Выражение с одним вещественным число, должно возвращать одно число в слайсе",
		},
	}

	for _, tc := range cases {
		t.Run(tc.nameCase, func(t *testing.T) {
			output, err := tokenizeExpression(tc.actual.(string))
			assert.Equal(output, tc.expected, tc.expectedMsg)
			if err != nil {
				assert.Error(err)
				assert.Equal(err, tc.err)
			}
		})
	}
}

func testChar(t *testing.T) {
	assert := assert.New(t)
	cases := []CaseTokenize{
		{
			nameCase:    "One_invlid_char_1",
			actual:      "P",
			expected:    nil,
			expectedMsg: "Выражние с символами кроме (+ - * / . ) должно возвращать nil и ошибку",
			err:         stdErr.New("неккоректный символ: \"P\", на позиции 0 (P)"),
		},
		{
			nameCase:    "One_invlid_char_2",
			actual:      "E",
			expected:    nil,
			expectedMsg: "Выражние с символами кроме (+ - * / . ) должно возвращать nil и ошибку",
			err:         stdErr.New("неккоректный символ: \"E\", на позиции 0 (E)"),
		},
		{
			nameCase:    "One_invlid_char_3",
			actual:      "y",
			expected:    nil,
			expectedMsg: "Выражние с символами кроме (+ - * / . ) должно возвращать nil и ошибку",
			err:         stdErr.New("неккоректный символ: \"y\", на позиции 0 (y)"),
		},
		{
			nameCase:    "One_invlid_char_4",
			actual:      "x",
			expected:    nil,
			expectedMsg: "Выражние с символами кроме (+ - * / . ) должно возвращать nil и ошибку",
			err:         stdErr.New("неккоректный символ: \"x\", на позиции 0 (x)"),
		},
		{
			nameCase:    "One_valid_char_1",
			actual:      "+",
			expected:    []string{"+"},
			expectedMsg: "Выражние с символами (+ - * / . ) должно возвращать этот символ с слайсе",
		},
		{
			nameCase:    "One_valid_char_2",
			actual:      "-",
			expected:    []string{"-"},
			expectedMsg: "Выражние с символами (+ - * / . ) должно возвращать этот символ с слайсе",
		},
		{
			nameCase:    "One_valid_char_3",
			actual:      "*",
			expected:    []string{"*"},
			expectedMsg: "Выражние с символами (+ - * / . ) должно возвращать этот символ с слайсе",
		},
		{
			nameCase:    "One_valid_char_4",
			actual:      "/",
			expected:    []string{"/"},
			expectedMsg: "Выражние с символами (+ - * / . ) должно возвращать этот символ с слайсе",
		},
		{
			nameCase:    "One_vlid_char_5",
			actual:      ".",
			expected:    []string{"."},
			expectedMsg: "Выражние с символами (+ - * / . ) должно возвращать этот символ с слайсе",
		},
	}

	for _, tc := range cases {
		t.Run(tc.nameCase, func(t *testing.T) {
			output, err := tokenizeExpression(tc.actual.(string))
			assert.Equal(output, tc.expected, tc.expectedMsg)
			if err != nil {
				assert.Error(err)
				assert.Equal(err, tc.err)
			}
		})
	}
}

func testBrackets(t *testing.T) {
	assert := assert.New(t)
	cases := []CaseTokenize{
		{
			nameCase:    "Empty bracket group",
			actual:      " (  )    ",
			expected:    nil,
			expectedMsg: "Выражение с пустыми скобками должно возвращать nil",
			err:         calculatorErr.ErrEmptyExpression,
			errMsg:      "Выражение с пустыми скобками должно возвращать ошибку",
		},
		{
			nameCase:    "Invalid brackets 1",
			actual:      " ) (",
			expected:    nil,
			expectedMsg: "Выражение с неккоректными скобками должно возвращать nil",
			err:         calculatorErr.ErrMissingBracket,
			errMsg:      "Выражение с неккореткными скобками должно возвращать ошибку",
		},
		{
			nameCase:    "Invalid brackets 2",
			actual:      " ) (()  ))) (",
			expected:    nil,
			expectedMsg: "Выражение с неккоректными скобками должно возвращать nil",
			err:         calculatorErr.ErrMissingBracket,
			errMsg:      "Выражение с неккореткными скобками должно возвращать ошибку",
		},
		{
			nameCase:    "Invalid brackets 3",
			actual:      "()()()() )",
			expected:    nil,
			expectedMsg: "Выражение с неккоректными скобками должно возвращать nil",
			err:         calculatorErr.ErrMissingBracket,
			errMsg:      "Выражение с неккореткными скобками должно возвращать ошибку",
		},
	}

	for _, tc := range cases {
		t.Run(tc.nameCase, func(t *testing.T) {
			output, err := tokenizeExpression(tc.actual.(string))
			assert.Equal(output, tc.expected, tc.expectedMsg)
			if assert.Error(err) {
				assert.Equal(err, tc.err, tc.errMsg)
			}
		})
	}
}
