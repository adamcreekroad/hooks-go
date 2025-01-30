package config

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/joho/godotenv"
)

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
	cacheDir, present := os.LookupEnv("CACHE_DIR")

	if !present {
		cwd, _ := os.Getwd()

		cacheDir = fmt.Sprintf("%s/cache", cwd)
	}

	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		os.Mkdir(cacheDir, 0755)
	}

	return cacheDir
}

func init() {
	configureEnvironment()
	configureRedis()
	configureRouter()
}

func configureEnvironment() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	_, filename, _, _ := runtime.Caller(1)
	file := path.Join(path.Dir(filename), "../.env."+env)
	godotenv.Load(file)

	godotenv.Load()
}
