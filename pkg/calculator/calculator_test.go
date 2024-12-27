package calculator

import (
	"errors"
	"strconv"
	"testing"
)

func TestCalcSuccess(t *testing.T) {
	testCasesSuccess := []struct {
		name           string
		expression     string
		expectedResult float64
	}{
		{
			name:           "simple_1",
			expression:     "1",
			expectedResult: 1,
		},
		{
			name:           "simple_2",
			expression:     "   42",
			expectedResult: 42,
		},
		{
			name:           "simple_3",
			expression:     "42  ",
			expectedResult: 42,
		},
		{
			name:           "simple_4",
			expression:     "   42   ",
			expectedResult: 42,
		},
		{
			name:           "simple_5",
			expression:     "2 + 2 * 2",
			expectedResult: 6,
		},
		{
			name:           "simple_5",
			expression:     "2 + 2 * 2",
			expectedResult: 6,
		},
		{
			name:           "paranthesis_1",
			expression:     "(1 + 2) * (10 - 4)",
			expectedResult: 18,
		},
		{
			name:           "paranthesis_2",
			expression:     "18 * (1 + 2)",
			expectedResult: 54,
		},
		{
			name:           "priority_1",
			expression:     "(2+2)*2",
			expectedResult: 8,
		},
		{
			name:           "priority_2",
			expression:     "2+2 * 2",
			expectedResult: 6,
		},
		{
			name:           "division_1",
			expression:     "1 /2",
			expectedResult: 0.5,
		},
		{
			name:           "division_2",
			expression:     "11 / (0-2)",
			expectedResult: -5.500000,
		},
		{
			name:           "expression_1",
			expression:     "11 * 99 / 35 - (235 * 42 * (5 - 4) - 555)",
			expectedResult: -9283.885714285714,
		},
		{
			name:           "expression_2",
			expression:     "35 + 35 * (35 / 35)",
			expectedResult: 70,
		},
		{
			name:           "expression_3",
			expression:     "11 / 111",
			expectedResult: 0.0990990990990991,
		},
	}
	for _, testCase := range testCasesSuccess {
		t.Run(testCase.name, func(t *testing.T) {
			val, err := Calc(testCase.expression)
			if err != nil {
				t.Errorf("successful case %s returns error %s", testCase.expression, err)
			}
			output := strconv.FormatFloat(val, 'f', -1, 64)
			expected := strconv.FormatFloat(testCase.expectedResult, 'f', -1, 64)
			if output != expected {
				t.Errorf("%s should be equal %s", output, expected)
			}
		})

	}
	testCasesFail := []struct {
		name        string
		expression  string
		expectedErr error
	}{
		{
			name:        "empty_0",
			expression:  "",
			expectedErr: ErrEmptyExpression,
		},
		{
			name:        "empty_1",
			expression:  " ",
			expectedErr: ErrEmptyExpression,
		},
		{
			name:        "empty_2",
			expression:  "   ",
			expectedErr: ErrEmptyExpression,
		},
		{
			name:        "bad_simple_1",
			expression:  "1+1*",
			expectedErr: ErrNotEnogthOperand,
		},
		{
			name:        "bad_priority_1",
			expression:  "2+2**2",
			expectedErr: ErrNotEnogthOperand,
		},
		{
			name:        "right paranthes",
			expression:  "((2+2-*(2",
			expectedErr: ErrMissRightParanthesis,
		},
		{
			name:        "left paranthes",
			expression:  "2+2)-2",
			expectedErr: ErrMissLeftParanthesis,
		},
		{
			name:        "division by zero",
			expression:  "10/0",
			expectedErr: ErrDivisionByZero,
		},
	}

	for _, testCase := range testCasesFail {
		t.Run(testCase.name, func(t *testing.T) {
			val, err := Calc(testCase.expression)
			if err == nil {
				t.Errorf("expression %s is invalid but result  %f was obtained", testCase.expression, val)
			}
			if err != nil {
				if !errors.Is(err, testCase.expectedErr) {
					t.Errorf("expected err: %v, recieved %v", testCase.expectedErr, err)
				}
			}
		})
	}
}

func TestCalcFail(t *testing.T) {
	testCasesFail := []struct {
		name        string
		expression  string
		expectedErr error
	}{
		{
			name:        "empty_0",
			expression:  "",
			expectedErr: ErrEmptyExpression,
		},
		{
			name:        "empty_1",
			expression:  " ",
			expectedErr: ErrEmptyExpression,
		},
		{
			name:        "empty_2",
			expression:  "   ",
			expectedErr: ErrEmptyExpression,
		},
		{
			name:        "bad_simple_1",
			expression:  "1+1*",
			expectedErr: ErrNotEnogthOperand,
		},
		{
			name:        "bad_priority_1",
			expression:  "2+2**2",
			expectedErr: ErrNotEnogthOperand,
		},
		{
			name:        "right paranthes",
			expression:  "((2+2-*(2",
			expectedErr: ErrMissRightParanthesis,
		},
		{
			name:        "left paranthes",
			expression:  "2+2)-2",
			expectedErr: ErrMissLeftParanthesis,
		},
		{
			name:        "division by zero",
			expression:  "10/0",
			expectedErr: ErrDivisionByZero,
		},
	}

	for _, testCase := range testCasesFail {
		t.Run(testCase.name, func(t *testing.T) {
			val, err := Calc(testCase.expression)
			if err == nil {
				t.Errorf("expression %s is invalid but result  %f was obtained", testCase.expression, val)
			}
			if err != nil {
				if !errors.Is(err, testCase.expectedErr) {
					t.Errorf("expected err: %v, recieved %v", testCase.expectedErr, err)
				}
			}
		})
	}
}
