package config

import (
	"os"

	"github.com/joho/godotenv"
)

func init() {
	configureEnvironment()
}

func configureEnvironment() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func Binding() string {
	binding, present := os.LookupEnv("BINDING")

	if !present {
		binding = "127.0.0.1"
	}

	return binding
}

func Port() string {
	port, present := os.LookupEnv("PORT")

	if !present {
		port = "8080"
	}

	return port
}
