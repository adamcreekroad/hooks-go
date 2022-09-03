package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/adamcreekroad/hooks-go/config"
)

type Embed struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Thumbnail   Thumbnail `json:"thumbnail"`
	Image       Image     `json:"image"`
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

type Image struct {
	Url string `json:"url"`
}

type Author struct {
	Name    string `json:"name"`
	IconUrl string `json:"icon_url"`
}

const apiUrl = "https://discord.com/api/v10"

var botToken = os.Getenv("DISCORD_BOT_TOKEN")

func SendMessage(channelID string, payload Payload, f []*os.File) {
	url := fmt.Sprintf("%s/channels/%s/messages", apiUrl, channelID)
	encoded, _ := json.Marshal(payload)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, _ := writer.CreateFormField("payload_json")
	part.Write(encoded)

	for index, file := range f {
		filename, _ := filepath.Rel(config.CacheDir(), file.Name())
		field, _ := writer.CreateFormFile(fmt.Sprintf("files[%d]", index), filename)

		_, err := io.Copy(field, file)

		if err != nil {
			panic(err)
		}
	}

	writer.Close()

	request, _ := http.NewRequest("POST", url, body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Set("Authorization", fmt.Sprintf("Bot %s", botToken))

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	if response.StatusCode >= 400 {
		body, _ := io.ReadAll(response.Body)

		fmt.Println("Failed to send Discord message: ", string(body))
	}
}
