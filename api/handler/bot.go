package handler

import (
	"NetlifyBot/pkg/utils"
	"fmt"
	"github.com/gofiber/fiber"
	"net/http"
	"strings"
)

func BOTBuild(c *fiber.Ctx) {
	body := utils.ParseQueryString(c.Body())
	if body["response_url"] == "" {
		c.Status(400)
		c.JSON(map[string] interface{} {
			"status": "Response URL not supplied",
		})
	} else if body["command"] == "/build" {
		var text string
		if body["text"] == "" {
			text = "Project name not supplied."
		}
		msg := strings.NewReader(fmt.Sprintf(`{"text":"%s"}`, text))
		r, _ := http.Post(body["response_url"],"application/json",msg)
		responseBody := make([]byte, 100)
		n, _ := r.Body.Read(responseBody)
		c.SendString(string(responseBody[:n]))
	}
}
