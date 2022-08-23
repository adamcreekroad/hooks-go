package plex

import (
	"fmt"
	"io"
	"log"
	"math"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/adamcreekroad/hooks-go/config"
	"github.com/adamcreekroad/hooks-go/discord"
	"github.com/gofrs/uuid"
)

func SendBulkLibraryNewMessage() {
	ids := fetch_event_ids()

	if len(ids) == 0 {
		return
	}

	message := discord.Payload{Content: "Recently added to Plex:", Tts: false, Embeds: []discord.Embed{}}

	var events []event
	var thumbs []*os.File

OUTER:
	for _, id := range ids {
		event := fetch_cached_payload(id)

		for _, e := range events {

			// Ignore duplicates
			if e.Metadata.GUID == event.Metadata.GUID {
				continue OUTER
			}
		}

		thumb := fetch_cached_thumb(id)

		events = append(events, event)
		thumbs = append(thumbs, thumb)

		append_library_new_item(event, thumb, &message)
	}

	discord.SendMessage(channel_id, message, thumbs)

	for _, id := range ids {
		result := config.RedisConn.Get(config.RedisConn.Context(), fmt.Sprintf("plex:thumb:%s", id))

		filename := result.Val()

		config.RedisConn.Del(config.RedisConn.Context(), fmt.Sprintf("plex:event:%s", id), fmt.Sprintf("plex:thumb:%s", id))
		config.RedisConn.SRem(config.RedisConn.Context(), "plex:library.new", id)

		if err := os.Remove(filename); err != nil {
			log.Println("Failed to remove file: ", err)
		}
	}
}

func process_library_new_hook(p string, t *multipart.FileHeader) {
	id, _ := uuid.NewV4()

	file, err := t.Open()

	if err != nil {
		panic(err)
	}

	filename := fmt.Sprintf("%s/plex-thumb-%s%s", config.CacheDir(), id, filepath.Ext(t.Filename))

	bytes, _ := io.ReadAll(file)

	if err := os.WriteFile(filename, bytes, 0644); err != nil {
		panic(err)
	}

	config.RedisConn.Set(config.RedisConn.Context(), fmt.Sprintf("plex:event:%s", id.String()), p, 0)
	config.RedisConn.Set(config.RedisConn.Context(), fmt.Sprintf("plex:thumb:%s", id.String()), filename, 0)
	config.RedisConn.SAdd(config.RedisConn.Context(), "plex:library.new", id.String())
}

func fetch_event_ids() []string {
	var ids []string

	config.RedisConn.SMembers(config.RedisConn.Context(), "plex:library.new").ScanSlice(&ids)

	return ids
}

func append_library_new_item(e event, t *os.File, p *discord.Payload) {
	switch e.Metadata.Type {
	case "show":
		build_library_new_show_message(e, p, t)
	case "episode":
		build_library_new_episode_message(e, p, t)
	case "movie":
		build_library_new_movie_message(e, p, t)
	}
}

func build_library_new_show_message(e event, message *discord.Payload, t *os.File) {
	filename, _ := filepath.Rel(config.CacheDir(), t.Name())
	url := fmt.Sprintf("attachment://%s", filename)
	title := fmt.Sprintf("%s - Season %d", e.Metadata.Title, e.Metadata.Index)
	description := fmt.Sprintf("**%d**  `%s`\n%s", e.Metadata.Year, e.Metadata.ContentRating, e.Metadata.Summary)

	embed := discord.Embed{
		Author:      discord.Author{Name: e.Metadata.GrandparentTitle},
		Title:       title,
		Description: description,
		Thumbnail:   discord.Thumbnail{Url: url},
	}

	message.Embeds = append(message.Embeds, embed)
}

func build_library_new_episode_message(e event, message *discord.Payload, t *os.File) {
	filename, _ := filepath.Rel(config.CacheDir(), t.Name())
	url := fmt.Sprintf("attachment://%s", filename)
	summary := fmt.Sprintf("||%s||", e.Metadata.Summary)
	title := fmt.Sprintf("S%d E%d - %s", e.Metadata.ParentIndex, e.Metadata.Index, e.Metadata.Title)
	description := fmt.Sprintf("%s - `%s`\n%s", time.Duration(e.Metadata.Duration)*time.Millisecond, e.Metadata.ContentRating, summary)

	embed := discord.Embed{
		Author:      discord.Author{Name: e.Metadata.GrandparentTitle},
		Title:       title,
		Description: description,
		Thumbnail:   discord.Thumbnail{Url: url},
	}

	message.Embeds = append(message.Embeds, embed)
}

func build_library_new_movie_message(e event, message *discord.Payload, t *os.File) {
	filename, _ := filepath.Rel(config.CacheDir(), t.Name())
	url := fmt.Sprintf("attachment://%s", filename)
	fields := []discord.Field{{Name: fmt.Sprintf("`%s`", e.Metadata.ContentRating), Value: e.Metadata.Summary}}
	description := fmt.Sprintf("**%d**  %s", e.Metadata.Year, time.Duration(e.Metadata.Duration)*time.Millisecond)

	embed := discord.Embed{
		Title: e.Metadata.Title, Description: description, Thumbnail: discord.Thumbnail{Url: url}, Fields: fields,
	}

	message.Embeds = append(message.Embeds, embed)
}

func seconds_to_human_time(seconds int32) string {
	hours := math.Floor(float64((seconds%31536000)%86400) / 3600)
	minutes := math.Floor(float64(((seconds%31536000)%86400)%3600) / 60)

	return fmt.Sprintf("%d hr %d min", int16(hours), int16(minutes))
}
