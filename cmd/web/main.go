package main

import (
	"fmt"
	"log"

	"github.com/adamcreekroad/hooks-go/config"
	"github.com/adamcreekroad/hooks-go/plex"
	"github.com/gin-gonic/gin"
)

func plexHook(c *gin.Context) {
	payload := c.PostForm("payload")
	thumb, _ := c.FormFile("thumb")

	log.Println(payload)

	plex.ProcessHook(payload, thumb)
}

func main() {
	config.Router.POST("/plex", plexHook)

	addr := fmt.Sprintf("%s:%s", config.Binding(), config.Port())

	config.Router.Run(addr)
}
