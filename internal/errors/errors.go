package errors

import (
	stdErr "errors"
)

var (
	ErrDivisionByZero    = stdErr.New("деление на ноль")
	ErrEmptyExpression   = stdErr.New("пустое выражение")
	ErrMissingBracket    = stdErr.New("пропущена скобка")
	ErrInvalidExpression = stdErr.New("неккоректное выражение")
	ErrNotEnogthOperand  = stdErr.New("не достаточно операндов")
)
