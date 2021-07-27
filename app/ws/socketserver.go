package ws

import (
	"github.com/gin-gonic/gin"
	"talkServer/internal/router"
	"talkServer/internal/ws"
	//"talkServer/internal/ws"
)

// websocket 服务
func SocketServer() {
	////由于是外部调用包，所以必须含包名 gin. 作为前缀
	////Default 返回带有已连接 Logger 和 Recovery 中间件的 Engine 实例。
	r := gin.Default()
	// 创建ws管理
	manager := ws.NewManager()
	go manager.Start()
	// 运行ws 管理
	//go ctrlWs.Manager.Start()
	//// Engine 结构体中内嵌了 RouterGroup 结构体，即继承了 RouterGroup（其有成员方法 GET、POST、DELETE、PUT、ANY 等）
	//router.Router(r)
	router.Router(r, manager)
	//// 默认是 0.0.0.0:8080 端口，内部使用了 http.ListenAndServe(address, engine)
	r.Run("0.0.0.0:8000") // listen and serve on 0.0.0.0:8000
}
