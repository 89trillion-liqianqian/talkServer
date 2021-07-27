package router

import (
	"github.com/gin-gonic/gin"
	"log"
	"talkServer/internal/ctrlWs"
	"talkServer/internal/ws"
)

// 路由管路i
func Router(r *gin.Engine, manager *ws.ClientManager) {
	r.GET("/ws", func(context *gin.Context) {
		log.Println("-----ws")
		ctrlWs.WsHandler(context, manager)
	})
}
