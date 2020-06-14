package main

import (
	"NetlifyBot/api/handler"
	"github.com/gofiber/fiber"
	"log"
	"os"
)
func main() {
	app := fiber.New()
	api := app.Group("/api", Logger)
	v1 := api.Group("/v1")
	hooks := v1.Group("/hooks")
	hooks.Post("/:project", handler.HooksHandler)
	bot := v1.Group("/bot")
	bot.Post("/build", handler.BOTBuild)
	bot.Post("/remove-project", handler.BOTRemoveProject)
	bot.Post("/read-project", handler.BOTReadProject)
	bot.Post("/show-projects",handler.BOTShowProjects)
	bot.Post("/add-project",handler.BOTAddProject)
	bot.Post("/update-project",handler.BOTUpdateProject)
	bot.Post("/project-logs", handler.BOTProjectLogs)
	app.Listen(os.Getenv("PORT"))
}
func Logger(ctx *fiber.Ctx) {
	log.Println("Origin IP:",ctx.IP(), " PATH:",ctx.OriginalURL())
	ctx.Next()
}
