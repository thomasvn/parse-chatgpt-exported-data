package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type ChatGPTConversation struct {
	Title       string                       `json:"title"`
	CreateTime  float64                      `json:"create_time"`
	UpdateTime  float64                      `json:"update_time"`
	MessageData map[string]MessageDataObject `json:"mapping"`
	ID          string                       `json:"id"`
}

type MessageDataObject struct {
	ID       string   `json:"id"`
	Message  Message  `json:"message"`
	Parent   string   `json:"parent"`
	Children []string `json:"children"`
}

type Message struct {
	Content Content `json:"content"`
}

type Content struct {
	Parts []string `json:"parts"`
}

func main() {
	// Open the file
	file, err := os.Open("./data/conversations.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Read from the file
	bytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Decode the JSON
	var ChatGPTConversationData []ChatGPTConversation
	err = json.Unmarshal(bytes, &ChatGPTConversationData)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Pretty print the JSON
	// prettyJSON, err := json.MarshalIndent(ChatGPTConversationData, "", "  ")
	// if err != nil {
	// 	fmt.Println("Error encoding JSON:", err)
	// 	return
	// }
	// fmt.Println(string(prettyJSON))

	// Iterate through each conversation
	for _, conversation := range ChatGPTConversationData {

		fmt.Println("\n--------------------------------------------------------------------------------")
		fmt.Println("# ", conversation.Title)

		// Iterate through each message and save parent/child IDs
		IDToChildID := map[string]string{}
		var root string
		for key, value := range conversation.MessageData {
			if len(value.Children) != 0 {
				IDToChildID[key] = value.Children[0]
			}
			if value.Parent == "" {
				root = key
			}
		}

		// Start with the root message and continuously print child messages
		for {
			currentMessage := conversation.MessageData[root]
			if len(currentMessage.Message.Content.Parts) > 0 {
				if currentMessage.Message.Content.Parts[0] != "" {
					fmt.Println("\n## ", currentMessage.Message.Content.Parts[0])
				}
			}
			if currentMessage.Children == nil {
				break
			}
			root = IDToChildID[root]
		}
	}
}
