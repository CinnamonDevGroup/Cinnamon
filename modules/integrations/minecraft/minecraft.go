package minecraft

import (
	databaseHelper "github.com/AngelFluffyOokami/Cinnamon/modules/core/database"
	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
)

func Init(s *discordgo.Session, DB databaseHelper.DBstruct) {
	r := gin.Default()
	hub := newHub()
	go hub.run(s, DB)
	r.GET("/minecraftSocket", func(c *gin.Context) {
		serveWs(hub, c.Writer, c.Request, s)
	})
	r.Run("localhost:8080") // listen and serve on 0.0.0.0:8080
}
