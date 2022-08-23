package plex

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"os"
	"time"

	"github.com/adamcreekroad/hooks-go/config"
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

type Payload struct {
	Event event
	Thumb *os.File
}

var channel_id = os.Getenv("PLEX_DISCORD_CHANNEL_ID")

const WHITESPACE_CHAR = "\u200b"

func ProcessHook(p string, t *multipart.FileHeader) {
	event := parse_payload(p)

	switch event.Event {
	case "library.new":
		process_library_new_hook(p, t)
	}
}

func cache_payload(uuid string, p string) {

}

func fetch_cached_payload(uuid string) event {
	raw_event := config.RedisConn.Get(config.RedisConn.Context(), fmt.Sprintf("plex:event:%s", uuid)).Val()
	event := parse_payload(raw_event)

	return event
}

func fetch_cached_thumb(uuid string) *os.File {
	filename := config.RedisConn.Get(config.RedisConn.Context(), fmt.Sprintf("plex:thumb:%s", uuid)).Val()

	thumb, err := os.Open(filename)

	if err != nil {
		panic(err)
	}

	return thumb
}

func parse_payload(p string) event {
	var result event

	if err := json.Unmarshal([]byte(p), &result); err != nil {
		panic(err)
	}

	return result
}
