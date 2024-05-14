package plex

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/adamcreekroad/hooks-go/config"
	"github.com/adamcreekroad/hooks-go/discord"
	"github.com/gofrs/uuid"
)

const (
	Show    = "show"
	Episode = "episode"
	Movie   = "movie"
	Album   = "album"
)

func SendBulkLibraryNewMessage() {
	payloads := fetchCachedPayloads()

	if len(payloads) == 0 {
		return
	}

	message := discord.Payload{Content: "New on Plex:", Tts: false, Embeds: []discord.Embed{}}

	var events []event
	var thumbs []*os.File

OUTER:
	for _, payload := range payloads {
		for _, event := range events {

			// Ignore duplicates
			if event.Metadata.GUID == payload.Event.Metadata.GUID {
				continue OUTER
			}
		}

		if !validateMediaType(payload.Event) {
			return
		}

		events = append(events, payload.Event)
		thumbs = append(thumbs, payload.Thumb)

		appendLibraryNewItem(payload.Event, payload.Thumb, &message)
	}

	if len(events) > 0 {
		discord.SendMessage(discordChannelID, message, thumbs)
	}

	clearCachedThumbs(payloads)
}

func processLibraryNewHook(p string, t *multipart.FileHeader) {
	event := parsePayload(p)

	if !validateMediaType(event) {
		return
	}

	id, _ := uuid.NewV4()

	payload := Payload[*multipart.FileHeader]{ID: id, Event: event, Thumb: t}

	cachePayload(payload)
}

func validateMediaType(e event) bool {
	switch e.Metadata.Type {
	case Show, Episode, Movie:
		return true
	default:
		return false
	}
}

func clearCachedThumbs(payloads []Payload[*os.File]) {
	for _, payload := range payloads {
		os.Remove(payload.Thumb.Name())
	}
}

func appendLibraryNewItem(e event, t *os.File, p *discord.Payload) {
	switch e.Metadata.Type {
	case Show:
		buildShowMessage(e, p, t)
	case Episode:
		buildEpisodeMessage(e, p, t)
	case Movie:
		buildMovieMessage(e, p, t)
	}
}

func buildShowMessage(e event, message *discord.Payload, t *os.File) {
	filename, _ := filepath.Rel(config.CacheDir(), t.Name())
	url := fmt.Sprintf("attachment://%s", filename)
	description := fmt.Sprintf("**%d**\n**`%s`**\n\n%s", e.Metadata.Year, e.Metadata.ContentRating, e.Metadata.Summary)

	embed := discord.Embed{
		Author:      discord.Author{Name: e.Metadata.LibrarySectionTitle},
		Title:       e.Metadata.Title,
		Description: description,
		Thumbnail:   discord.Thumbnail{Url: url},
	}

	message.Embeds = append(message.Embeds, embed)
}

func buildEpisodeMessage(e event, message *discord.Payload, t *os.File) {
	filename, _ := filepath.Rel(config.CacheDir(), t.Name())
	url := fmt.Sprintf("attachment://%s", filename)

	var summary string
	if e.Metadata.Summary != "" {
		summary = fmt.Sprintf("||%s||", e.Metadata.Summary)
	} else {
		summary = "No summary available."
	}

	description := fmt.Sprintf("**Season %d**\n**Episode %d · %s**\n**`%s`**\n\n%s", e.Metadata.ParentIndex, e.Metadata.Index, e.Metadata.Title, e.Metadata.ContentRating, summary)

	embed := discord.Embed{
		Author:      discord.Author{Name: e.Metadata.LibrarySectionTitle},
		Title:       e.Metadata.GrandparentTitle,
		Description: description,
		Thumbnail:   discord.Thumbnail{Url: url},
	}

	message.Embeds = append(message.Embeds, embed)
}

func buildMovieMessage(e event, message *discord.Payload, t *os.File) {
	filename, _ := filepath.Rel(config.CacheDir(), t.Name())
	url := fmt.Sprintf("attachment://%s", filename)
	// TODO: Get the duration nicely formatted
	// duration := time.Duration(e.Metadata.Duration)*time.Millisecond
	description := fmt.Sprintf("**%d**\n**`%s`**\n\n%s", e.Metadata.Year, e.Metadata.ContentRating, e.Metadata.Summary)

	embed := discord.Embed{
		Author:      discord.Author{Name: e.Metadata.LibrarySectionTitle},
		Title:       e.Metadata.Title,
		Description: description,
		Thumbnail:   discord.Thumbnail{Url: url},
	}

	message.Embeds = append(message.Embeds, embed)
}

// func buildAlbumMessage(e event, message *discord.Payload, t *os.File) {
// 	filename, _ := filepath.Rel(config.CacheDir(), t.Name())
// 	url := fmt.Sprintf("attachment://%s", filename)
// 	description := fmt.Sprintf("**%d**  %s", e.Metadata.Year, e.Metadata.Genre[0].Tag)

// 	embed := discord.Embed{
// 		Author:      discord.Author{Name: e.Metadata.ParentTitle},
// 		Title:       e.Metadata.Title,
// 		Description: description,
// 		Thumbnail:   discord.Thumbnail{Url: url},
// 	}

// 	message.Embeds = append(message.Embeds, embed)
// }

// func seconds_to_human_time(seconds int32) string {
// 	hours := math.Floor(float64((seconds%31536000)%86400) / 3600)
// 	minutes := math.Floor(float64(((seconds%31536000)%86400)%3600) / 60)

// 	return fmt.Sprintf("%d hr %d min", int16(hours), int16(minutes))
// }
