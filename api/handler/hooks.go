package handler

import (
	"NetlifyBot/pkg/dbOperations"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber"
	"net/http"
	"strings"
)
func AddAttachment(dest *[]interface{}, srcTitle string, srcValue interface{}) {
	attachment := make(map[string] interface{})
	srcTitle=strings.ReplaceAll(srcTitle,"_"," ")
	attachment["title"] = strings.ToUpper(srcTitle)
	attachment["value"] = srcValue
	*dest = append(*dest,attachment)
}
func MessageBodyANDLogsInsertion(nm string,v string) string{
		var m map[string] interface{}
		json.Unmarshal([]byte(v),&m)
		responseMap := map[string] interface{} {}
		responseMap["text"] = m["title"].(string)
		var summaryAttachment = make([]interface{}, 0, 100)
		summaryBody := (m["summary"].(map[string] interface{}))["messages"].([] interface{})
		for _, v := range summaryBody {
			eventInfo := v.(map[string] interface{})
			attachment := make(map[string] interface{})
			attachment["title"] = eventInfo["title"]
			attachment["value"] = eventInfo["description"]
			attachment["short"] = false
			summaryAttachment = append(summaryAttachment, attachment)
		}
		status := m["state"]
		responseMap["attachments"] = []interface{} {
			map[string]interface{}{
				"text": "Summary",
				"color": "#",
				"fields": summaryAttachment,
			},
		}
		var projectInfoAttachment = make([]interface{}, 0, 100)

		for v, value := range m {
			switch v {
			case "name":
				AddAttachment(&projectInfoAttachment,v,value)
				break
			case "committer":
				if value != nil {
					AddAttachment(&projectInfoAttachment,v,value)
				}
				break
			case "created_at":
				AddAttachment(&projectInfoAttachment,v,value)
				break
			case "published_at":
				if value != nil {
					AddAttachment(&projectInfoAttachment,v,value)
				}
				break
			case "deploy_url":
				AddAttachment(&projectInfoAttachment,v,value)
				break
			case "deploy_ssl_url":
				AddAttachment(&projectInfoAttachment,v,value)
				break
			case "branch":
				AddAttachment(&projectInfoAttachment,v,value)
				break
			case "context":
				AddAttachment(&projectInfoAttachment,v,value)
				break
			case "deploy_time":
				if value != nil {
					AddAttachment(&projectInfoAttachment,v,value)
				}
				break
			case "error_message":
				if value != nil {
					AddAttachment(&projectInfoAttachment,v,value)
				}
				break
			case "Context":
				if value != nil {
					AddAttachment(&projectInfoAttachment,v,value)
				}
				break
			}
		}
		var statusColor string
		if status == "building" {
			statusColor = "#f7d488"
		} else if status == "ready" {
			statusColor = "72b381"
		} else {
			statusColor = "#ff0000"
		}
		if len(summaryAttachment) == 0 {
			responseMap["attachments"] = []interface{} {
				map[string]interface{}{
					"text": "Project Info",
					"color": statusColor,
					"fields": projectInfoAttachment,
				},
			}
		} else {
			responseMap["attachments"] = []interface{} {
				map[string]interface{}{
					"text": "Project Info",
					"color": statusColor,
					"fields": projectInfoAttachment,
				},
				map[string]interface{}{
					"text": "Summary",
					"color": "#464785",
					"fields": summaryAttachment,
				},
			}
		}

		b, _ := json.Marshal(responseMap)
	logs := m["log_access_attributes"].(map[string] interface{})
	dbOperations.InsertLogs(nm,logs["url"].(string)+".json")
	return string(b)
}
func HooksHandler(ctx *fiber.Ctx) {
	c := ctx.Params("project")
	v := dbOperations.ReadProject(c)
	if v.URL != "" {
		body := strings.NewReader(
				MessageBodyANDLogsInsertion(c,ctx.Body()))
		http.Post(v.SlackURL, "application/json",body)
		ctx.SendString(
			MessageBodyANDLogsInsertion(c,ctx.Body()))
	} else {
		ctx.Status(400)
		ctx.JSON(map[string] string {
			"status": fmt.Sprintf("Project [%s] not in the database.",c),
		})
	}
}
