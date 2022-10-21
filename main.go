package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"smartcardwifi/handle"
)

func main() {
	app := fiber.New()

	app.Use(cors.New())

	app.Static("/public", "./imgs")

	app.Get("/data", func(c *fiber.Ctx) error {
		data, err := handle.HandleReader(c)
		if err != nil {
			return c.Status(503).JSON(fiber.Map{
				"ok":    false,
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"ok":   true,
			"data": data,
		})
	})

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Status(503).JSON(fiber.Map{
			"message": "Hi",
		})
	})

	app.Listen(":9090")
}
