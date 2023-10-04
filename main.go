package main

import (
	"bytes"
	// "context"
	"encoding/json"
	"fmt"
	"github.com/a-h/templ"
	"io"
	"log"
	"net/http"
	// "os"
	"github.com/labstack/echo/v4"
	"strings"
)

// take in string
//
//	run through ai
//
// return html
func main() {
	e := echo.New()
	request_url := "http://localhost:11434/api/generate"
	jsonStr := []byte(`{"model":"mistral", "prompt":"Why is the sky blue?"}`)
	resp, err := http.Post(request_url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	string_arr := strings.Split(string(body), "\n")
	arr := make([]map[string]interface{}, 0)
	for _, line := range string_arr {
		var data map[string]interface{}
		json.Unmarshal([]byte(line), &data)
		arr = append(arr, data)
	}
	resp_arr := make([]string, 0)
	for _, line := range arr {
		s, ok := line["response"].(string)
		if ok {
			resp_arr = append(resp_arr, s)
		}
	}
	ret := strings.Join(resp_arr, "")
	component := hello(ret)
	http.Handle("/", templ.Handler(component))
	http.ListenAndServe(":8080", nil)
	// var data map[string]interface{}
	// if err := json.Unmarshal(body, &data); err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(data)
}
