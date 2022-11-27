package minecraft

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Init(s *discordgo.Session, DB *gorm.DB) {
	r := gin.Default()
	hub := newHub()
	go hub.run(s, DB)
	r.GET("/minecraftSocket", func(c *gin.Context) {
		serveWs(hub, c.Writer, c.Request, s)
	})
	r.Run("localhost:8080") // listen and serve on 0.0.0.0:8080
}
