package main

import (
	"NetlifyBot/api/handler"
	"github.com/gofiber/fiber"
	"log"
)
func main() {
	app := fiber.New()
	api := app.Group("/api", Logger)
	v1 := api.Group("/v1")
	hooks := v1.Group("/hooks")
	hooks.Post("/:project", handler.HooksHandler)
	app.Listen(80)
}
func Logger(ctx *fiber.Ctx) {
	log.Println("Origin IP:",ctx.IP(), " PATH:",ctx.OriginalURL())
	ctx.Next()
}
