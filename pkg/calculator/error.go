package calculator

import (
	stdErr "errors"
)

var (
	ErrDivisionByZero       = stdErr.New("division by zero")
	ErrEmptyExpression      = stdErr.New("empty expression")
	ErrNotEnogthOperand     = stdErr.New("not enogth operand")
	ErrInvalidExpression    = stdErr.New("invalid expression")
	ErrMissLeftParanthesis  = stdErr.New("miss left paranthesis")
	ErrMissRightParanthesis = stdErr.New("miss right paranthesis")
)
