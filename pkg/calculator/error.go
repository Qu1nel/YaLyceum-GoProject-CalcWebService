package calculator

import (
	stdErr "errors"
)

var (
	ErrEmptyExpression      = stdErr.New("emptry expression")
	ErrInvalidExpression    = stdErr.New("invalid expression")
	ErrNotEnogthOperand     = stdErr.New("not enogth operand")
	ErrDivisionByZero       = stdErr.New("division by zero")
	ErrMissLeftParanthesis  = stdErr.New("miss left paranthesis")
	ErrMissRightParanthesis = stdErr.New("miss right paranthesis")
)
