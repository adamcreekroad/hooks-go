package plex

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/adamcreekroad/hooks-go/discord"
)

type event struct {
	Event   string `json:"event"`
	User    bool   `json:"user"`
	Owner   bool   `json:"owner"`
	Account struct {
		ID    int    `json:"id"`
		Thumb string `json:"thumb"`
		Title string `json:"title"`
	} `json:"account"`
	Server struct {
		Title string `json:"title"`
		UUID  string `json:"uuid"`
	} `json:"server"`
	Player struct {
		Local         bool   `json:"local"`
		PublicAddress string `json:"publicAddress"`
		Title         string `json:"title"`
		UUID          string `json:"uuid"`
	} `json:"player"`
	Metadata struct {
		LibrarySectionType    string  `json:"librarySectionType"`
		RatingKey             string  `json:"ratingKey"`
		Key                   string  `json:"key"`
		ParentRatingKey       string  `json:"parentRatingKey"`
		GrandparentRatingKey  string  `json:"grandparentRatingKey"`
		GUID                  string  `json:"guid"`
		LibrarySectionID      int     `json:"librarySectionID"`
		Type                  string  `json:"type"`
		Title                 string  `json:"title"`
		GrandparentKey        string  `json:"grandparentKey"`
		ParentKey             string  `json:"parentKey"`
		GrandparentTitle      string  `json:"grandparentTitle"`
		ParentTitle           string  `json:"parentTitle"`
		Summary               string  `json:"summary"`
		Index                 int     `json:"index"`
		ParentIndex           int     `json:"parentIndex"`
		RatingCount           int     `json:"ratingCount"`
		Thumb                 string  `json:"thumb"`
		Art                   string  `json:"art"`
		ParentThumb           string  `json:"parentThumb"`
		GrandparentThumb      string  `json:"grandparentThumb"`
		GrandparentArt        string  `json:"grandparentArt"`
		AddedAt               int     `json:"addedAt"`
		UpdatedAt             int     `json:"updatedAt"`
		Studio                string  `json:"studio"`
		LibrarySectionTitle   string  `json:"librarySectionTitle"`
		LibrarySectionKey     string  `json:"librarySectionKey"`
		ContentRating         string  `json:"contentRating"`
		Rating                float32 `json:"rating"`
		AudienceRating        float32 `json:"audienceRating"`
		Year                  int16   `json:"year"`
		Tagline               string  `json:"tagline"`
		Duration              int32   `json:"duration"`
		OriginallyAvailableAt string  `json:"originallyAvailableAt"`
		AudienceRatingImage   string  `json:"audienceRatingImage"`
		PrimaryExtraKey       string  `json:"primaryExtraKey"`
		RatingImage           string  `json:"ratingImage"`
		Genre                 []struct {
			Id     int32  `json:"id"`
			Filter string `json:"filter"`
			Tag    string `json:"tag"`
			Count  int32  `json:"count"`
		} `json:"Genre"`
		Director []struct {
			Id     int32  `json:"id"`
			Filter string `json:"filter"`
			Tag    string `json:"tag"`
		} `json:"Director"`
		Writer []struct {
			Id     int32  `json:"id"`
			Filter string `json:"filter"`
			Tag    string `json:"tag"`
		} `json:"Writer"`
		Producer []struct {
			Id     int32  `json:"id"`
			Filter string `json:"filter"`
			Tag    string `json:"tag"`
			Count  int32  `json:"count"`
		} `json:"Producer"`
		Country []struct {
			Id     int32  `json:"id"`
			Filter string `json:"filter"`
			Tag    string `json:"tag"`
			Count  int32  `json:"count"`
		} `json:"Country"`
		Guid []struct {
			Id string `json:"id"`
		} `json:"Guid"`
		Ratings []struct {
			Image string  `json:"image"`
			Value float32 `json:"value"`
			Type  string  `json:"type"`
			Count int32   `json:"count"`
		} `json:"Rating"`
		Role []struct {
			Id     int32  `json:"id"`
			Filter string `json:"filter"`
			Tag    string `json:"tag"`
			TagKey string `json:"tagKey"`
			Count  int32  `json:"count"`
			Role   string `json:"role"`
			Thumb  string `json:"thumb"`
		} `json:"Role"`
	} `json:"metadata"`
}

var channel_id = os.Getenv("PLEX_DISCORD_CHANNEL_ID")

func ProcessHook(p string, t *multipart.FileHeader) {
	event := parse_payload(p)

	message := discord.Payload{}

	switch event.Event {
	case "library.new":
		build_library_new_message(event, &message, t)
	case "media.play":
		build_media_play_message(event, &message, t)
	}

	discord.SendMessage(channel_id, message, t)
}

func parse_payload(p string) event {
	var result event

	if err := json.Unmarshal([]byte(p), &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	return result
}

func build_library_new_message(e event, message *discord.Payload, t *multipart.FileHeader) {
	switch e.Metadata.Type {
	case "episode":
		build_library_new_episode_message(e, message, t)
	case "movie":
		build_library_new_movie_message(e, message, t)
	}
}

func build_library_new_episode_message(e event, message *discord.Payload, t *multipart.FileHeader) {
	message.Content = fmt.Sprintf(
		"S%d E%d of %s is now on Plex!",
		e.Metadata.ParentIndex,
		e.Metadata.Index,
		e.Metadata.GrandparentTitle,
	)

	description := fmt.Sprintf("||%s||", e.Metadata.Summary)
	url := fmt.Sprintf("attachment://%s", t.Filename)

	message.Tts = false
	message.Embeds = []discord.Embed{
		{Title: e.Metadata.Title, Description: description, Thumbnail: discord.Thumbnail{Url: url}},
	}
}

func build_library_new_movie_message(e event, message *discord.Payload, t *multipart.FileHeader) {
	message.Content = fmt.Sprintf("**%s** is now on Plex!", e.Metadata.Title)

	url := fmt.Sprintf("attachment://%s", t.Filename)

	fields := []discord.Field{{Name: "Year", Value: fmt.Sprintf("%d", e.Metadata.Year)}}

	message.Tts = false
	message.Embeds = []discord.Embed{
		{Title: e.Metadata.Title, Description: e.Metadata.Tagline, Thumbnail: discord.Thumbnail{Url: url}, Fields: fields},
	}
}

func build_media_play_message(e event, message *discord.Payload, t *multipart.FileHeader) {
	switch e.Metadata.Type {
	case "episode":
		build_media_play_episode_message(e, message, t)
	case "track":
		build_media_play_track_message(e, message, t)
	}
}

func build_media_play_episode_message(e event, message *discord.Payload, t *multipart.FileHeader) {
	message.Content = fmt.Sprintf(
		"%s is watching S%dE%d of %s",
		e.Account.Title, e.Metadata.ParentIndex, e.Metadata.Index, e.Metadata.GrandparentTitle,
	)

	description := fmt.Sprintf("||%s||", e.Metadata.Summary)

	url := fmt.Sprintf("attachment://%s", t.Filename)

	message.Tts = false
	message.Embeds = []discord.Embed{{Title: e.Metadata.Title, Description: description, Thumbnail: discord.Thumbnail{Url: url}}}
}

func build_media_play_track_message(e event, message *discord.Payload, t *multipart.FileHeader) {
	message.Content = fmt.Sprintf(
		"%s is jammin' to %s by %s", e.Account.Title, e.Metadata.Title, e.Metadata.GrandparentTitle,
	)

	url := fmt.Sprintf("attachment://%s", t.Filename)

	message.Tts = false
	message.Embeds = []discord.Embed{{Author: discord.Author{Name: e.Account.Title, IconUrl: e.Account.Thumb}, Title: e.Metadata.Title, Description: e.Metadata.Summary, Thumbnail: discord.Thumbnail{Url: url}}}
}
