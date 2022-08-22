package main

import (
	"fmt"

	"github.com/adamcreekroad/hooks-go/config"
	"github.com/adamcreekroad/hooks-go/plex"
	"github.com/gin-gonic/gin"
)

func plex_hook(c *gin.Context) {
	payload := c.PostForm("payload")
	thumb, _ := c.FormFile("thumb")

	fmt.Println(payload)

	plex.ProcessHook(payload, thumb)
}

func main() {
	config.Router.POST("/plex", plex_hook)

	addr := fmt.Sprintf("%s:%s", config.Binding(), config.Port())

	config.Router.Run(addr)
}
