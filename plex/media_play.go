package plex

// func send_media_play_message(e event, t *multipart.FileHeader) {
// 	message := discord.Payload{}

// 	switch e.Metadata.Type {
// 	case "episode":
// 		build_media_play_episode_message(e, &message, t)
// 	case "track":
// 		build_media_play_track_message(e, &message, t)
// 	}

// 	discord.SendMessage(channel_id, message, t)
// }

// func build_media_play_episode_message(e event, message *discord.Payload, t *multipart.FileHeader) {
// 	message.Content = fmt.Sprintf(
// 		"%s is watching S%dE%d of %s",
// 		e.Account.Title, e.Metadata.ParentIndex, e.Metadata.Index, e.Metadata.GrandparentTitle,
// 	)

// 	description := fmt.Sprintf("||%s||", e.Metadata.Summary)

// 	url := fmt.Sprintf("attachment://%s", t.Filename)

// 	message.Tts = false
// 	message.Embeds = []discord.Embed{{Title: e.Metadata.Title, Description: description, Thumbnail: discord.Thumbnail{Url: url}}}
// }

// func build_media_play_track_message(e event, message *discord.Payload, t *multipart.FileHeader) {
// 	message.Content = fmt.Sprintf(
// 		"%s is jammin' to %s by %s", e.Account.Title, e.Metadata.Title, e.Metadata.GrandparentTitle,
// 	)

// 	url := fmt.Sprintf("attachment://%s", t.Filename)

// 	message.Tts = false
// 	message.Embeds = []discord.Embed{{Author: discord.Author{Name: e.Account.Title, IconUrl: e.Account.Thumb}, Title: e.Metadata.Title, Description: e.Metadata.Summary, Thumbnail: discord.Thumbnail{Url: url}}}
// }
