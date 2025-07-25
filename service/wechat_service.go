package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// WeChatMessage represents the incoming message from WeChat (JSON format)
type WeChatMessage struct {
	ToUserName   string `json:"ToUserName"`
	FromUserName string `json:"FromUserName"`
	CreateTime   int64  `json:"CreateTime"`
	MsgType      string `json:"MsgType"`
	Content      string `json:"Content"`
	MsgId        int64  `json:"MsgId"`
}

// WeChatReply represents the reply message to WeChat (JSON format)
type WeChatReply struct {
	ToUserName   string `json:"ToUserName"`
	FromUserName string `json:"FromUserName"`
	CreateTime   int64  `json:"CreateTime"`
	MsgType      string `json:"MsgType"`
	Content      string `json:"Content"`
}

// MessagePushHandler handles WeChat message push events
func MessagePushHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Received message push request: %s %s\n", r.Method, r.URL.Path)
	
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Error reading request body: %v\n", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fmt.Printf("Received JSON body: %s\n", string(body))

	// Parse the JSON message
	var message WeChatMessage
	err = json.Unmarshal(body, &message)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	fmt.Printf("Parsed message - From: %s, To: %s, Content: %s\n", 
		message.FromUserName, message.ToUserName, message.Content)

	// Process the message and generate reply
	reply := processMessage(&message)

	// Convert reply to JSON
	replyJSON, err := json.MarshalIndent(reply, "", "  ")
	if err != nil {
		fmt.Printf("Error generating reply JSON: %v\n", err)
		http.Error(w, "Failed to generate reply", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Sending reply: %s\n", string(replyJSON))

	// Set response headers
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(replyJSON)
}

// processMessage processes the incoming message and returns a reply
func processMessage(message *WeChatMessage) *WeChatReply {
	var replyContent string
	
	// Handle different message types
	switch message.MsgType {
	case "text":
		// For text messages, append "!" to the content
		replyContent = message.Content + "!"
	case "image":
		replyContent = "收到图片消息！"
	case "voice":
		replyContent = "收到语音消息！"
	case "video":
		replyContent = "收到视频消息！"
	case "location":
		replyContent = "收到位置消息！"
	case "link":
		replyContent = "收到链接消息！"
	default:
		replyContent = "收到消息！"
	}

	reply := &WeChatReply{
		ToUserName:   message.FromUserName, // Swap ToUserName and FromUserName
		FromUserName: message.ToUserName,
		CreateTime:   time.Now().Unix(),
		MsgType:      "text",
		Content:      replyContent,
	}

	return reply
} 