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
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/net/websocket"
)

func get_response(message string, currentContext []int) string {
	request_url := "http://localhost:11434/api/generate"
	var jsonBytes []byte
	if len(currentContext) == 0 {
		jsonBytes = []byte(`{"model":"mistral", "prompt":"` + message + `"}`)
	} else {
		currentContextBytes, err := json.Marshal(currentContext)
		if err != nil {
			fmt.Println(err)
			return get_response(message, []int{})
		}
		jsonBytes = []byte(
			`{"model":"mistral", "prompt":"` + message +
				`", "context":` + string(currentContextBytes[:]) + `}`,
		)
	}
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
	// build context for next request
	chatContextBytes, err := json.Marshal(chatContext)
	if err != nil {
		log.Fatal(err)
	}
	var bytes bytes.Buffer
	paragraph := strings.Join(ret_arr, "")
	err = botMessage(paragraph, string(chatContextBytes)).Render(context.Background(), &bytes)
	if err != nil {
		log.Fatal(err)
	}
	return bytes.String()
}

type UserRequest struct {
	Entry   string `json:"entry"`
	Context []int  `json:"context"`
}

func wsHandler(c echo.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		for {
			userJson := ""
			err := websocket.Message.Receive(ws, &userJson)
			if err != nil {
				c.Logger().Error(err)
				break
			}
			var requestJson map[string]interface{}
			fmt.Println(userJson)
			err = json.Unmarshal([]byte(userJson), &requestJson)
			if err != nil {
				c.Logger().Error(err)
				break
			}
			userMsg := requestJson["entry"].(string)
			contextJson := requestJson["context"]
			var currentContextStr string
			if contextJson == nil {
				fmt.Println("context is nil")
				currentContextStr = "[]"
			} else {
				currentContextStr = contextJson.(string)
			}
			var currentContext []int
			err = json.Unmarshal([]byte(currentContextStr), &currentContext)
			if err != nil {
				c.Logger().Error(err)
			}
			var bytes bytes.Buffer
			err = userMessage(userMsg).Render(context.Background(), &bytes)
			if err != nil {
				log.Fatal(err)
			}
			userHTML := bytes.String()
			err = websocket.Message.Send(ws, userHTML)
			if err != nil {
				c.Logger().Error(err)
				break
			}
			chatReply := get_response(userMsg, currentContext)
			err = websocket.Message.Send(ws, chatReply)
			if err != nil {
				c.Logger().Error(err)
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

// TODO: make context so that we refer to the element that is oob.
func main() {
	e := echo.New()
	e.Static("/static", "static")
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// e.POST("/request", func(c echo.Context) error {
	// 	message := c.FormValue("entry")
	// 	currentContext := c.FormValue("context")
	// 	chatReply := get_response(message, currentContext)
	// 	return c.HTML(http.StatusOK, chatReply)
	// })
	e.GET("/", func(c echo.Context) error {
		return c.File("static/index.html")
	})
	e.GET("/ws", wsHandler)
	fmt.Println("Server started at port 1323")
	e.Logger.Fatal(e.Start(":1323"))
}
