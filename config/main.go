package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	configureEnvironment()
	configureRedis()
	configureRouter()
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

func CacheDir() string {
	cache_dir, present := os.LookupEnv("CACHE_DIR")

	if !present {
		cwd, _ := os.Getwd()

		cache_dir = fmt.Sprintf("%s/cache", cwd)
	}

	if _, err := os.Stat(cache_dir); os.IsNotExist(err) {
		os.Mkdir(cache_dir, 0755)
	}

	return cache_dir
}
