package middlewares

import (
	"net/http"

	"github.com/labstack/echo"
)

func Auth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			cc := c.(*Context)
			if cc.SessionUserId() == "" {
				return echo.NewHTTPError(http.StatusUnauthorized)
			}
			return next(c)
		}
	}
}
