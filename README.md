# hooks-go

Webhook handler/processor written in Go.

As of right now, it only ingests plex webhooks for the `library.new` event, and will send messages to the specified discord channel for all items when running the notifier command.
