# hooks-go

Webhook handler/processor written in Go.

As of right now, it only ingests plex webhooks for the `library.new` event, and will send messages to the specified discord channel for all items when running the notifier command.

## Installation

Installation requires Ruby for scripting, and the application itself requires Go to run.

1. Clone or download this repo into `/opt/hooks`
2. Copy the `.env.example` into `.env` and configure it to your needs
3. Run `./bin/install`

This will create and start a systemd service for the webserver, as well as add a cron job for the notifier.

## Local Development

### Setup

You will need golang to run the application itself, and ruby for local scripts.

Run the setup script to prepare your local environment:
```shell-script
./bin/setup
```

This will create a `.env.development` file to configure and use for local development. Be sure to set the required variables within there.

### Running the app

This service consists of two processes:
- notifier
- web

The `web` process is a running webserver that will handle incoming webhooks, storing the payload within Redis for the `notifier` process to pick up. To run it:
```shell-script
go run cmd/web/main.go
```

The `notifier` process checks for any payloads that may be stored in Redis, sending a message to the configured Discord channel for each one. To run it:
```shell-script
go run cmd/notifier/main.go
```

