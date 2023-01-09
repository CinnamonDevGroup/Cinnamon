package websocket

import (
	"github.com/gin-gonic/gin"
)

func Init() {
	r := gin.Default()
	hub := newHub()
	go hub.run()
	r.GET("/minecraftSocket", func(c *gin.Context) {
		serveWs(hub, c.Writer, c.Request)
	})
	r.Run("localhost:8080") // listen and serve on 0.0.0.0:8080
}
