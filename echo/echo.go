package echoredoc

import (
	"github.com/labstack/echo/v4"
	"github.com/luisnquin/go-redoc"
)

// Handler sets some defaults and returns a HandlerFunc.
func EchoHandler(r redoc.Redoc) echo.HandlerFunc {
	return func(c echo.Context) error {
		r.Handler()(c.Response(), c.Request())

		return nil
	}
}
