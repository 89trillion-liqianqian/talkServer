package ctrlWs

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"talkServer/internal/myerr"
	"talkServer/internal/ws"
	"talkServer/utils"
)

//Ws Handler socket 连接 中间件 作用:升级协议,用户验证,自定义信息等
func WsHandler(c *gin.Context, manager *ws.ClientManager) {
	log.Println("-------aaa")
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		myerr.SendHttpErr(c)
		return
	}
	// postman 使用设置username
	username := c.Request.Header.Get("username")
	client := &ws.Client{
		ID:       utils.GetUID(),
		Socket:   conn,
		Send:     make(chan []byte),
		UserName: username,
		Manager:  manager,
	}
	manager.Register <- client
	go client.Read()
	go client.Write()
}
