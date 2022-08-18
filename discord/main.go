package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

type Embed struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Thumbnail   Thumbnail `json:"thumbnail"`
	Author      Author    `json:"author"`
	Fields      []Field   `json:"fields"`
}

type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

type Payload struct {
	Content string  `json:"content"`
	Tts     bool    `json:"tts"`
	Embeds  []Embed `json:"embeds"`
}

type Thumbnail struct {
	Url string `json:"url"`
}

type Author struct {
	Name    string `json:"name"`
	IconUrl string `json:"icon_url"`
}

const api_url = "https://discord.com/api/v10"
const content_type = "application/json; charset=UTF-8"

var bot_token = os.Getenv("DISCORD_BOT_TOKEN")

func SendMessage(channel_id string, payload Payload, t *multipart.FileHeader) {
	url := fmt.Sprintf("%s/channels/%s/messages", api_url, channel_id)
	encoded, _ := json.Marshal(payload)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, _ := writer.CreateFormField("payload_json")
	part.Write(encoded)

	field, _ := writer.CreateFormFile("files[0]", t.Filename)

	file, _ := t.Open()

	_, err := io.Copy(field, file)

	if err != nil {
		panic(err)
	}

	writer.Close()

	request, _ := http.NewRequest("POST", url, body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", fmt.Sprintf("Bot %s", bot_token))

	client := &http.Client{}
	response, error := client.Do(request)

	if error != nil {
		panic(error)
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		body, _ := ioutil.ReadAll(response.Body)

		fmt.Println("Failed to send Discord message: ", string(body))
	}
}

type Media struct {
	File   multipart.File
	Header *multipart.FileHeader
}

func FileUrl(m Media) string {
	return fmt.Sprintf("attachment://%s", m.Header.Filename)
}
