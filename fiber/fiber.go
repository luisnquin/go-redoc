package fiberredoc

import (
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/luisnquin/go-redoc"
)

func FiberHandler(r redoc.Redoc) fiber.Handler {
	return func(c *fiber.Ctx) error {
		adaptor.HTTPHandlerFunc(r.Handler())

		return nil
	}
}
