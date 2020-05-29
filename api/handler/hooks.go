package handler

import (
	"NetlifyBot/pkg/dbOperations"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber"
	"net/http"
	"strings"
)

func MessageBodyANDLogsInsertion(nm string,v string) string{
	var c map[string] interface{}
	json.Unmarshal([]byte(v), &c)
	summary := c["summary"].(map[string] interface{})
	var s = c["title"].(string) + "\n Status: " + summary["status"].(string)
	for v, value := range c {
		switch v {
		case "name":
			s += fmt.Sprintf("\n Name: %s", value)
			break
		case "committer":
			if value != nil {
				s += fmt.Sprintf("\n Commiter: %s", value)
			}
			break
		case "created_at":
			s += fmt.Sprintf("\n Created At: %s", value)
			break
		case "published_at":
			if value != nil {
				s += fmt.Sprintf("\n Published At: %s", value)
			}
			break
		case "deploy_url":
			s += fmt.Sprintf("\n Deploy URL: %s", value)
			break
		case "deploy_ssl_url":
			s += fmt.Sprintf("\n Deploy SSL URL: %s", value)
			break
		case "branch":
			s += fmt.Sprintf("\n Branch: %s", value)
			break
		case "context":
			s += fmt.Sprintf("\n Context: %s", value)
			break
		case "deploy_time":
			if value != nil {
				s += fmt.Sprintf("\n Deploy Time: %f", value.(float64))
			}
			break
		case "error_message":
			if value != nil {
				s += fmt.Sprintf("\n Error: %s", value)
			}
			break
		case "Context":
			if value != nil {
				s += fmt.Sprintf("\n Context: %s", value)
			}
			break
		}
	}
	message := summary["messages"].([] interface{})
	if len(message) > 0 {
		s += "\n SUMMARY:"
		for i, v := range message {
			value := v.(map[string]interface{})
			s += fmt.Sprintf("\n\t[%d] %s", i+1, value["title"])
		}
	}
	logs := c["log_access_attributes"].(map[string] interface{})
	dbOperations.InsertLogs(nm,logs["url"].(string)+".json")
	return s
}
func HooksHandler(ctx *fiber.Ctx) {
	c := ctx.Params("project")
	v := dbOperations.ReadProject(c)
	if v.URL != "" {
		body := strings.NewReader(
			fmt.Sprintf(`{"text": "%s"}`,
				MessageBodyANDLogsInsertion(c,ctx.Body())))
		http.Post(v.SlackURL, "application/json",body)
		ctx.SendString(fmt.Sprintf(`{"text": "%s"}`,
			MessageBodyANDLogsInsertion(c,ctx.Body())))
	} else {
		ctx.Status(400)
		ctx.JSON(map[string] string {
			"status": fmt.Sprintf("Project [%s] not in the database.",c),
		})
	}
}
