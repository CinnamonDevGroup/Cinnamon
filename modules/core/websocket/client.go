package websocket

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AngelFluffyOokami/Cinnamon/modules/core/commonutils"

	"github.com/google/uuid"

	"github.com/gorilla/websocket"
)

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.

func (c *Client) readPump() {
	defer commonutils.RecoverPanic("")
	config := <-commonutils.GetConfig
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var data IncomingData
		err := c.Conn.ReadJSON(&data)
		if err != nil {
			if config.Debugging {
				commonutils.LogEvent("WS Unexpected Disconnect Event: "+fmt.Sprint(err), commonutils.LogWarning)
			}
			break
		}

		userMessage, err := json.Marshal(data)
		if err != nil {
			if config.Debugging {
				commonutils.LogEvent("JSON Marshal Error Event: "+fmt.Sprint(err), commonutils.LogError)
			}
		}
		c.Hub.Broadcast <- userMessage
		c.Hub.Client <- c
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	defer commonutils.RecoverPanic("")
	config := <-commonutils.GetConfig
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				if config.Debugging {
					commonutils.LogEvent("WritePump NextWriter Error Event: "+fmt.Sprint(err), commonutils.LogError)
				}
				return
			}

			_, err = w.Write(message)

			if err != nil {
				if config.Debugging {
					commonutils.LogEvent("WritePump Write Error: "+fmt.Sprint(err), commonutils.LogError)
				}
			}

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				_, err = w.Write(newline)
				if err != nil {
					if config.Debugging {
						commonutils.LogEvent("WritePump Write Error: "+fmt.Sprint(err), commonutils.LogError)
					}
				}
				_, err = w.Write(<-c.Send)
				if err != nil {
					if config.Debugging {
						commonutils.LogEvent("WritePump Write Error: "+fmt.Sprint(err), commonutils.LogError)
					}
				}
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				if config.Debugging {
					commonutils.LogEvent("Set Write Deadline Error Event: "+fmt.Sprint(err), commonutils.LogError)
				}
				return
			}
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				if config.Debugging {
					commonutils.LogEvent("WritePump Write Message Error Event: "+fmt.Sprint(err), commonutils.LogError)
				}
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	defer commonutils.RecoverPanic("")
	config := <-commonutils.GetConfig
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if config.Debugging {
			commonutils.LogEvent("serverWS Upgrader Error Event: "+fmt.Sprint(err), commonutils.LogError)
		}
		return
	}
	client := &Client{Hub: hub, Conn: conn, Send: make(chan []byte, 256)}
	client.Hub.Register <- client
	client.Addr = conn.RemoteAddr().String()

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()

}

func GenUserId() string {
	uid := uuid.NewString()
	return uid
}
