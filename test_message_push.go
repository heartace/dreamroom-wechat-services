package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// Test JSON message that WeChat would send
	testJSON := `{
		"ToUserName": "gh_1234567890",
		"FromUserName": "o_abcdefghijklmnop",
		"CreateTime": 1640995200,
		"MsgType": "text",
		"Content": "Hello World",
		"MsgId": 12345678901234567
	}`

	// Send POST request to the message-push endpoint
	resp, err := http.Post("http://localhost:80/message-push", "application/json", bytes.NewBufferString(testJSON))
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Read and print the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	fmt.Printf("Response Status: %s\n", resp.Status)
	fmt.Printf("Response Body:\n%s\n", string(body))
} 