package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	httpErr "CalcService/internal/error"

	"CalcService/pkg/calculator"
	"CalcService/pkg/logger"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.uber.org/zap"
)

type Request struct {
	Expression string `json:"expression"`
}
type Response struct {
	Result float64 `json:"result"`
}

type Server struct {
	server *echo.Echo
}

func New(ctx context.Context, port int, pattern_url string) (*Server, error) {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{AllowMethods: []string{"POST", "OPTION"}}))
	e.Use(middleware.LoggerWithConfig(
		middleware.LoggerConfig{
			Format:           `{time: ${time_rfc3339}, host: ${host}, method: ${method}, uri: ${uri}, status: ${status}, error: ${error}}` + "\n",
			CustomTimeFormat: "2006-05-08 15:55:05",
		},
	))

	server := &Server{e}
	server.server.POST(pattern_url, server.calculate)

	return server, nil
}

func (s *Server) Start(port string) error {
	return s.server.Start(port)
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *Server) calculate(c echo.Context) error {
	if c.Request().Header.Get("Content-Type") != "application/json" {
		logger.New().Info(context.Background(), "content type not allowed", zap.String("Content-Type", c.Request().Header.Get("Content-Type")), zap.String("required Content-Type", "application/json"))
		return httpErr.NotValidExpression
	}

	data, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return httpErr.NotValidExpression
	}

	req := Request{}
	err = json.Unmarshal(data, &req)

	if err != nil {
		logger.New().Info(context.Background(), "req", zap.String("expression", err.Error()))
		return httpErr.NotValidExpression
	}

	if req.Expression == "internal" {
		return httpErr.InternalServer
	}

	res, err := calculator.Calc(req.Expression)
	if err != nil {
		return httpErr.NotValidExpression
	}

	return c.JSON(http.StatusOK, Response{Result: res})
}
