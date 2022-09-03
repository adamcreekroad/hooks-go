package plex

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/adamcreekroad/hooks-go/config"
	"github.com/gofrs/uuid"
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
		LibrarySectionType    string        `json:"librarySectionType"`
		RatingKey             string        `json:"ratingKey"`
		Key                   string        `json:"key"`
		ParentRatingKey       string        `json:"parentRatingKey"`
		GrandparentRatingKey  string        `json:"grandparentRatingKey"`
		GUID                  string        `json:"guid"`
		LibrarySectionID      int           `json:"librarySectionID"`
		Type                  string        `json:"type"`
		Title                 string        `json:"title"`
		GrandparentKey        string        `json:"grandparentKey"`
		ParentKey             string        `json:"parentKey"`
		GrandparentTitle      string        `json:"grandparentTitle"`
		ParentTitle           string        `json:"parentTitle"`
		Summary               string        `json:"summary"`
		Index                 int           `json:"index"`
		ParentIndex           int           `json:"parentIndex"`
		RatingCount           int           `json:"ratingCount"`
		Thumb                 string        `json:"thumb"`
		Art                   string        `json:"art"`
		ParentThumb           string        `json:"parentThumb"`
		GrandparentThumb      string        `json:"grandparentThumb"`
		GrandparentArt        string        `json:"grandparentArt"`
		AddedAt               int           `json:"addedAt"`
		UpdatedAt             int           `json:"updatedAt"`
		Studio                string        `json:"studio"`
		LibrarySectionTitle   string        `json:"librarySectionTitle"`
		LibrarySectionKey     string        `json:"librarySectionKey"`
		ContentRating         string        `json:"contentRating"`
		Rating                float32       `json:"rating"`
		AudienceRating        float32       `json:"audienceRating"`
		Year                  int16         `json:"year"`
		Tagline               string        `json:"tagline"`
		Duration              time.Duration `json:"duration"`
		OriginallyAvailableAt string        `json:"originallyAvailableAt"`
		AudienceRatingImage   string        `json:"audienceRatingImage"`
		PrimaryExtraKey       string        `json:"primaryExtraKey"`
		RatingImage           string        `json:"ratingImage"`
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

type Payload[T any] struct {
	ID    uuid.UUID
	Event event
	Thumb T
}

type payloadJSON struct {
	ID        uuid.UUID `json:"id"`
	Event     event     `json:"event"`
	ThumbPath string    `json:"thumb_path"`
}

var discordChannelID = os.Getenv("PLEX_DISCORD_CHANNEL_ID")

func ProcessHook(p string, t *multipart.FileHeader) {
	event := parsePayload(p)

	switch event.Event {
	case "library.new":
		processLibraryNewHook(p, t)
	}
}

func cachePayload(p Payload[*multipart.FileHeader]) {
	file, err := p.Thumb.Open()

	if err != nil {
		panic(err)
	}

	filename := fmt.Sprintf("%s/plex-thumb-%s%s", config.CacheDir(), p.ID, filepath.Ext(p.Thumb.Filename))

	bytes, _ := io.ReadAll(file)

	if err := os.WriteFile(filename, bytes, 0644); err != nil {
		panic(err)
	}

	data, _ := json.Marshal(payloadJSON{ID: p.ID, Event: p.Event, ThumbPath: filename})

	config.RedisConn.RPush(config.RedisConn.Context(), "plex:library.new", data)
}

func fetchCachedPayloads() []Payload[*os.File] {
	var payloads []Payload[*os.File]

	payloadCount := config.RedisConn.LLen(config.RedisConn.Context(), "plex:library.new").Val()

	for i := 0; i < int(payloadCount); i++ {
		var payload payloadJSON

		rawPayload := config.RedisConn.LPop(config.RedisConn.Context(), "plex:library.new").Val()

		if err := json.Unmarshal([]byte(rawPayload), &payload); err != nil {
			panic(err)
		}

		thumb, err := os.Open(payload.ThumbPath)

		if err != nil {
			panic(err)
		}

		payloads = append(payloads, Payload[*os.File]{ID: payload.ID, Event: payload.Event, Thumb: thumb})
	}

	return payloads
}

func parsePayload(p string) event {
	var result event

	if err := json.Unmarshal([]byte(p), &result); err != nil {
		panic(err)
	}

	return result
}
