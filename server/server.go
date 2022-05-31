package server

import (
	"billing/backend/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func RegisterRoutes(app *fiber.App, handlers []handlers.Handler) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))
	for _, handler := range handlers {
		handler.RegisterRoutes(app)
	}
}
