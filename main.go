package main

import (
	"bytes"
	"context"
	"fmt"

	// "github.com/a-h/templ"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func get_response(message string, currentContext string) string {
	request_url := "http://localhost:11434/api/generate"
	var jsonBytes []byte
	msg := strings.ReplaceAll(message, "\"", "")
	if currentContext == "" {
		jsonBytes = []byte(`{"model":"mistral", "prompt":"` + msg + `"}`)
	} else {
		jsonBytes = []byte(`{"model":"mistral", "prompt":"` + msg + `", "context":[` + currentContext + `]}`)
	}

	fmt.Println(string(jsonBytes))
	resp, err := http.Post(request_url, "application/json", bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	body_arr := strings.Split(string(body), "\n")
	json_arr := make([]map[string]interface{}, 0)
	for _, line := range body_arr {
		var data map[string]interface{}
		json.Unmarshal([]byte(line), &data)
		json_arr = append(json_arr, data)
	}
	ret_arr := make([]string, 0)
	var chatContext []int
	for _, line := range json_arr {
		s, ok := line["response"].(string)
		if ok {
			ret_arr = append(ret_arr, s)
			continue
		}
		contextValue, ok := line["context"].([]interface{})
		if ok {
			for _, val := range contextValue {
				chatContext = append(chatContext, int(val.(float64)))
			}
		}
	}
	chatCtxString := strings.Trim(strings.Join(strings.Split(fmt.Sprint(chatContext), " "), ","), "[]")
	var bytes bytes.Buffer
	paragraph := strings.Join(ret_arr, "")
	err = postReply(paragraph, chatCtxString).Render(context.Background(), &bytes)
	if err != nil {
		log.Fatal(err)
	}
	return bytes.String()
}
func main() {
	e := echo.New()
	// keep context somehow
	e.Static("/static", "static")
	e.POST("/request", func(c echo.Context) error {
		message := c.FormValue("entry")
		currentContext := c.FormValue("context")
		chatReply := get_response(message, currentContext)
		return c.HTML(http.StatusOK, chatReply)
	})
	e.GET("/", func(c echo.Context) error {
		return c.File("static/index.html")
	})
	fmt.Println("Server started at port 1323")
	e.Logger.Fatal(e.Start(":1323"))
}
