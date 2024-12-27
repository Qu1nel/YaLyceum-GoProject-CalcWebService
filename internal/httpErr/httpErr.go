package httpErr

import (
	"net/http"

	"github.com/labstack/echo"
)

var NotValidExpression = echo.NewHTTPError(http.StatusUnprocessableEntity, "Expression is not valid")
var InternalServer = echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
