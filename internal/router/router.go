package router

import (
	"github.com/gin-gonic/gin"
	"talkServer/internal/ctrlWs"
)

// 路由管路i
func Router(r *gin.Engine) {
	r.GET("/ws", ctrlWs.WsHandler)
}
