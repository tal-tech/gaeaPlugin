package xesgin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func GetWebSocketConn(c *gin.Context, upgraders ...websocket.Upgrader) (*websocket.Conn, error) {
	if len(upgraders) > 0 {
		up := upgraders[0]
		return up.Upgrade(c.Writer, c.Request, nil)
	}
	return upGrader.Upgrade(c.Writer, c.Request, nil)
}
