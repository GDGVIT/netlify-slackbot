package handler

import (
	"NetlifyBot/pkg/dbOperations"
	ut "NetlifyBot/pkg/utils"
	"fmt"
	"github.com/gofiber/fiber"
	"net/http"
	"strings"
)
func BOTProjectLogs(c *fiber.Ctx) {
	m := ut.ParseQueryString(c.Body())
	if m["text"] == "" {
		ut.SlackMessage(m["response_url"], "Project name not supplied.")
		c.SendString("Project name not supplied")
	} else {
		text := m["text"]
		subs := strings.Split(text,"+")
		requiredProject := dbOperations.ReadProject(subs[0])
		if requiredProject.URL == "" {
			ut.SlackMessage(m["response_url"],fmt.Sprintf(
				"Project [%s] not found.", subs[0]))
			c.SendString("Project name not found")
		} else {
			ut.SlackMessage(m["response_url"],dbOperations.GetLogs(subs[0]))
			c.SendString("Project Logs Sent.")
		}
	}
}
func BOTUpdateProject(c *fiber.Ctx) {
	c.SendString("")
	m := ut.ParseQueryString(c.Body())
	if m["text"] == "" {
		ut.SlackMessage(m["response_url"], "Project name not supplied.")
		//c.SendString("Project name not supplied")
	} else {
		text := m["text"]
		subs := strings.Split(text,"+")
		requiredProject := dbOperations.ReadProject(subs[0])
		if requiredProject.URL == "" {
			ut.SlackMessage(m["response_url"],fmt.Sprintf(
				"Project [%s] not found.", subs[0]))
			//c.SendString("Project name not found")
		} else {
			if len(subs) != 3 {
				ut.SlackMessage(m["response_url"], fmt.Sprintf(
					"Project [%s] doesn't have required number of parameters.\n Follow the format \n PROJECT_NAME NETLIFY_HOOK_URL SLACK_HOOK_URL", subs[0]))
			} else {
				if dbOperations.UpdateProject(subs[0], subs[1], subs[2]) == 0 {
					ut.SlackMessage(m["response_url"], fmt.Sprintf(
						"Project [%s] been updated.", subs[0]))
					//c.SendString("Project Updated")
				} else {
					ut.SlackMessage(m["response_url"], fmt.Sprintf(
						"Project [%s] update failed.", subs[0]))
					//c.SendString("Project update failed ")
				}
			}
		}
	}
}

func BOTAddProject(c *fiber.Ctx) {
	c.SendString("")
	m := ut.ParseQueryString(c.Body())
	if m["text"] == "" {
		ut.SlackMessage(m["response_url"], "Project name not supplied.")
		//c.SendString("Project name not supplied")
	} else {
		text := m["text"]
		subs := strings.Split(text,"+")
		if len(subs) != 3 {
			ut.SlackMessage(m["response_url"], fmt.Sprintf(
				"Project [%s] doesn't have required number of parameters.\n Follow the format \n PROJECT_NAME NETLIFY_HOOK_URL SLACK_HOOK_URL", subs[0]))
		} else {
			if dbOperations.InsertProject(subs[0], subs[1], subs[2]) == 0 {
				ut.SlackMessage(m["response_url"], fmt.Sprintf(
					"Project [%s] been added.", subs[0]))
				//c.SendString("Project Added")
			} else {
				ut.SlackMessage(m["response_url"], fmt.Sprintf(
					"Project [%s] already exists.", subs[0]))
				//c.SendString("Project already exists")
			}
		}
	}
}

func BOTShowProjects(c *fiber.Ctx) {
	c.SendString("")
	m := ut.ParseQueryString(c.Body())
	l := dbOperations.AllProjects()
	var s = "Projects"
	for i,v := range l {
		s += fmt.Sprintf("\n[%d] %s", i+1, v.Name)
	}
	ut.SlackMessage(m["response_url"],s)
	//c.SendString("values sent.")
}
func BOTBuild(c *fiber.Ctx) {
	c.SendString("")
	m := ut.ParseQueryString(c.Body())
	if m["text"] == "" {
		//c.SendString("")
		ut.SlackMessage(m["response_url"], "Project name not supplied.")

	} else {
		text := m["text"]
		subs := strings.Split(text,"+")
		requiredProject := dbOperations.ReadProject(subs[0])
		if requiredProject.URL == "" {
			c.Send()
			ut.SlackMessage(m["response_url"],fmt.Sprintf(
				"Project [%s] not found.", subs[0]))
		} else {
			ut.SlackMessage(m["response_url"],fmt.Sprintf(
				"Project [%s] build started.", subs[0]))
			http.Post(requiredProject.URL,"application/json",strings.NewReader("{}"))

		}
	}
}
func BOTReadProject(c *fiber.Ctx) {
	c.SendString("")
	m := ut.ParseQueryString(c.Body())
	if m["text"] == "" {
		ut.SlackMessage(m["response_url"], "Project name not supplied.")
		//c.SendString("Project name not supplied")
	} else {
		text := m["text"]
		subs := strings.Split(text,"+")
		requiredProject := dbOperations.ReadProject(subs[0])
		if requiredProject.URL == "" {
			ut.SlackMessage(m["response_url"],fmt.Sprintf(
				"Project [%s] not found.", subs[0]))
			//c.SendString("Project name not found")
		} else {
			ut.SlackMessage(m["response_url"],fmt.Sprintf(
				"Project [%s]\nNetlify Hook: %s\nSlack Hook: %s", subs[0], requiredProject.URL, requiredProject.SlackURL))
			//c.SendString("Project Returned")
		}
	}
}
func BOTRemoveProject(c *fiber.Ctx) {
	c.SendString("")
	m := ut.ParseQueryString(c.Body())
	if m["text"] == "" {
		ut.SlackMessage(m["response_url"], "Project name not supplied.")
		//c.SendString("Project name not supplied")
	} else {
		text := m["text"]
		subs := strings.Split(text,"+")
		requiredProject := dbOperations.ReadProject(subs[0])
		if requiredProject.URL == "" {
			ut.SlackMessage(m["response_url"],fmt.Sprintf(
				"Project [%s] not found.", subs[0]))
			//c.SendString("Project name not found")
		} else {
			ut.SlackMessage(m["response_url"],fmt.Sprintf(
				"Project [%s] Removed.", subs[0]))
			dbOperations.DeleteProjects(subs[0])
			//c.SendString("Project Removed")
		}
	}
}
func _ (c *fiber.Ctx) {
	c.SendString("")
	body := ut.ParseQueryString(c.Body())
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
		r.Body.Read(responseBody)
		//c.SendString(string(responseBody[:n]))
	}
}
