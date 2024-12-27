package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestCalcHandler200(t *testing.T) {
	testCases200 := []struct {
		name           string
		Expression     string `json:"expression"`
		expectedResult float64
		expectedStatus int
		method         string
	}{
		{
			name:           "simple_1",
			Expression:     "1",
			expectedResult: 1,
			method:         http.MethodPost,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "simple_2",
			Expression:     "42  ",
			expectedResult: 42,
			method:         http.MethodPost,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "simple_3",
			Expression:     "    42",
			expectedResult: 42,
			method:         http.MethodPost,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "simple_4",
			Expression:     "   42      ",
			expectedResult: 42,
			method:         http.MethodPost,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "simple_5",
			Expression:     "1+1",
			expectedResult: 2,
			method:         http.MethodPost,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "simple_6",
			Expression:     "  1  +   1",
			expectedResult: 2,
			method:         http.MethodPost,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "priority_1",
			Expression:     "(2+2)*2",
			expectedResult: 8,
			method:         http.MethodPost,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "priority_2",
			Expression:     "2+2*2 - 3 * (100 - 35)",
			expectedResult: -189,
			method:         http.MethodPost,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "division_1",
			Expression:     "1/2",
			expectedResult: 0.5,
			method:         http.MethodPost,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "division_2",
			Expression:     "11/111",
			expectedResult: 0.0990990990990991,
			method:         http.MethodPost,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "expression_1",
			Expression:     "1/2 * 99 - (45 / 235 * 999)",
			expectedResult: -141.79787234042553,
			method:         http.MethodPost,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "expression_2",
			Expression:     "3984 * 9385 / 385983 - 384539485398583 + (35 / 23 * 348953) - 99 * 9999",
			expectedResult: -3.845394858573717e+14,
			method:         http.MethodPost,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "expression_3",
			Expression:     "42 * 42 * 42 * 42 * 42 * 42 * 42 * 42 * 42 * 42 * 42 * 42 * 42 * 42 * 42 * 42 * 42",
			expectedResult: 3.937657486715347e+27,
			method:         http.MethodPost,
			expectedStatus: http.StatusOK,
		},
	}

	for _, testCase := range testCases200 {
		t.Run(testCase.name, func(t *testing.T) {
			var b bytes.Buffer

			if err := json.NewEncoder(&b).Encode(testCase); err != nil {
				t.Error("failed encode test cases", err)
				return
			}

			e := echo.New()
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(testCase.method, "/api/v1/calculate", &b)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := echo.New().NewContext(req, rec)

			resp := Response{}
			server := Server{server: e}
			if err := server.calculate(c); err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil && err != io.EOF {
				t.Errorf("unexpected error: %v", err)
				return
			}

			assert.Equal(t, testCase.expectedResult, resp.Result)
		})
	}
}

func TestCalcHandler422(t *testing.T) {
	testCases422 := []struct {
		name           string
		Expression     string `json:"expression"`
		expectedStatus int
		expectedErrMsg string
	}{
		{
			name:           "some bad req",
			Expression:     "x(50-50-50)",
			expectedErrMsg: "code=422, message=Expression is not valid (not enogth operand)",
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "simple1",
			Expression:     "1+1*",
			expectedErrMsg: "code=422, message=Expression is not valid (not enogth operand)",
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "priority",
			Expression:     "2+2**2",
			expectedErrMsg: "code=422, message=Expression is not valid (not enogth operand)",
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "right paranthes",
			Expression:     "((2+2-*(2",
			expectedErrMsg: "code=422, message=Expression is not valid (miss right paranthesis)",
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "left paranthes",
			Expression:     "2+2)-2",
			expectedErrMsg: "code=422, message=Expression is not valid (miss left paranthesis)",
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "empty",
			Expression:     "",
			expectedErrMsg: "code=422, message=Expression is not valid (empty expression)",
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "division by zero",
			Expression:     "10/0",
			expectedErrMsg: "code=422, message=Expression is not valid (division by zero)",
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}
	for _, testCase := range testCases422 {
		t.Run(testCase.name, func(t *testing.T) {
			var b bytes.Buffer
			err := json.NewEncoder(&b).Encode(testCase)
			if err != nil {
				t.Error("failed encode test cases", err)
				return
			}

			e := echo.New()
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", &b)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := echo.New().NewContext(req, rec)

			server := Server{server: e}
			err = server.calculate(c) //nolint:errcheck
			resp := Response{}

			if testCase.expectedStatus != 200 {
				recCode := strings.Split(err.Error(), "=")
				recCode = strings.Split(recCode[1], ",")
				recievedCode, err2 := strconv.Atoi(recCode[0])

				if err2 != nil {
					t.Errorf("unexpected error: %v", err)
				}

				assert.Equal(t, testCase.expectedErrMsg, err.Error())
				assert.Equal(t, testCase.expectedStatus, recievedCode)
			}

			err = json.NewDecoder(rec.Body).Decode(&resp)

			if err != nil && err != io.EOF {
				t.Errorf("unexpected error: %v", err)
				return
			}
		})
	}
}

func TestCalcHandler500(t *testing.T) {
	testCases500 := []struct {
		name           string
		Expression     string `json:"expression"`
		expectedStatus int
		expectedErrMsg string
	}{
		{
			name:           "internal_case",
			Expression:     "internal  ",
			expectedErrMsg: "code=500, message=Internal server error",
			expectedStatus: http.StatusInternalServerError,
		},
	}
	for _, testCase := range testCases500 {
		t.Run(testCase.name, func(t *testing.T) {
			var b bytes.Buffer
			var err = json.NewEncoder(&b).Encode(testCase)
			if err != nil {
				t.Error("failed encode test cases", err)
				return
			}

			e := echo.New()
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", &b)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := echo.New().NewContext(req, rec)

			server := Server{server: e}
			err = server.calculate(c) //nolint:errcheck
			resp := Response{}

			if testCase.expectedStatus != http.StatusInternalServerError {
				recCode := strings.Split(err.Error(), "=")
				recCode = strings.Split(recCode[1], ",")
				recievedCode, err2 := strconv.Atoi(recCode[0])

				if err2 != nil {
					t.Errorf("unexpected error: %v", err)
				}

				assert.Equal(t, testCase.expectedErrMsg, err.Error())
				assert.Equal(t, testCase.expectedStatus, recievedCode)
			}

			err = json.NewDecoder(rec.Body).Decode(&resp)

			if err != nil && err != io.EOF {
				t.Errorf("unexpected error: %v", err)
				return
			}
		})
	}
}
