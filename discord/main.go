package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Embed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Payload struct {
	Content string  `json:"content"`
	Tts     bool    `json:"tts"`
	Embeds  []Embed `json:"embeds"`
}

const api_url = "https://discord.com/api/v10"
const content_type = "application/json; charset=UTF-8"

var bot_token = os.Getenv("DISCORD_BOT_TOKEN")

func SendMessage(channel_id string, payload Payload) {
	url := fmt.Sprintf("%s/channels/%s/messages", api_url, channel_id)
	encoded, _ := json.Marshal(payload)
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(encoded))
	request.Header.Set("Content-Type", content_type)
	request.Header.Set("Authorization", fmt.Sprintf("Bot %s", bot_token))

	client := &http.Client{}
	response, error := client.Do(request)

	if error != nil {
		panic(error)
	}

	response.Body.Close()

	if response.StatusCode >= 300 {
		fmt.Println("Failed to send Discord message: ", response.Body)
	}
}
