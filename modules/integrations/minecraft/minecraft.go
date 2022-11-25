package minecraft

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
)

func Run(s *discordgo.Session) {
	r := gin.Default()
	hub := newHub()
	go hub.run(s)
	r.GET("/minecraftSocket", func(c *gin.Context) {
		serveWs(hub, c.Writer, c.Request, s)
	})
	r.Run("localhost:8080") // listen and serve on 0.0.0.0:8080
}
