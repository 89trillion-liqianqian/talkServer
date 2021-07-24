package ctrlWs

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"talkServer/internal/handler"
	"talkServer/internal/myerr"
	"talkServer/protocol/protobuf"
	"talkServer/utils"
)

// 定义命令
const (
	LoginCmd    = "login"    // 绑定 username
	PingCmd     = "ping"     // ping/pong
	TalkCmd     = "talk"     // 聊天
	UserListCmd = "userlist" // 获取用户列表
	ExitCmd     = "exit"     // 退出
)

// ClientManager is a websocket manager
type ClientManager struct {
	Clients    map[string]*Client
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

// Client is a websocket client
type Client struct {
	ID       string
	UserName string
	Socket   *websocket.Conn
	Send     chan []byte
}

// Message is return msg
type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

// Manager define a ws server manager
var Manager = ClientManager{
	Broadcast:  make(chan []byte),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
	Clients:    make(map[string]*Client),
}

// Start is  项目运行前, 协程开启start -> go Manager.Start()
func (manager *ClientManager) Start() {
	for {
		log.Println("<---管道通信--->")
		select {
		case conn := <-Manager.Register:
			log.Printf("新用户加入:%v", conn.ID)
			Manager.Clients[conn.ID] = conn
			jsonMessage, _ := json.Marshal(&Message{Content: "Successful connection to socket service"})
			conn.Send <- jsonMessage
		case conn := <-Manager.Unregister:
			log.Printf("用户离开:%v", conn.ID)
			bData, _ := handler.GetExitData(conn.UserName)
			if _, ok := Manager.Clients[conn.ID]; ok {
				conn.Send <- bData
				close(conn.Send)
				delete(Manager.Clients, conn.ID)
			}

		case message := <-Manager.Broadcast:
			MessageStruct := protobuf.TaskInfo{}
			// 解析数据，广播其他玩家
			proto.Unmarshal(message, &MessageStruct)
			for _, conn := range Manager.Clients {
				if conn.UserName == MessageStruct.Name || conn.UserName == "" {
					continue
				}
				select {
				case conn.Send <- message:
				default:
					close(conn.Send)
					delete(Manager.Clients, conn.ID)
				}
			}
		}
	}
}

// 获取在线用户列表
func getNamelist() string {
	// 数组默认长度为map长度,后面append时,不需要重新申请内存和拷贝,效率很高
	keys := make([]string, 0, len(Manager.Clients))
	for _, c := range Manager.Clients {
		if c.UserName != "" {
			keys = append(keys, c.UserName)
		}
	}
	keysB, _ := json.Marshal(keys)
	return string(keysB)
}

// 监听读数据
func (c *Client) Read() {
	defer func() {
		Manager.Unregister <- c
		c.Socket.Close()
	}()

	for {
		c.Socket.PongHandler()
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			Manager.Unregister <- c
			c.Socket.Close()
			break
		}
		receivedData := protobuf.TaskInfo{}
		proto.Unmarshal(message, &receivedData)
		// 登陆
		if receivedData.Cmd == LoginCmd {
			c.UserName = receivedData.Name
			continue
		}
		// ping/pong
		if receivedData.Cmd == PingCmd {
			bData, _ := handler.GetPingData()
			select {
			case c.Send <- bData:
			default:
				close(c.Send)
				delete(Manager.Clients, c.ID)
			}
			continue
		}
		// 获取玩家列表 userlist
		if receivedData.Cmd == UserListCmd {
			// 获取用户列表
			userList := getNamelist()
			bData, _ := handler.GetUserListData(userList, UserListCmd)
			select {
			case c.Send <- bData:
			default:
				close(c.Send)
				delete(Manager.Clients, c.ID)
			}
			continue
		}
		// 退出
		if receivedData.Cmd == ExitCmd && receivedData.Name == c.UserName {
			Manager.Unregister <- c
			c.Socket.Close()
			// 广播其他人
			Manager.Broadcast <- message
			break
		}

		Manager.Broadcast <- message
	}
}

// 监听写数据
func (c *Client) Write() {
	defer func() {
		c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			log.Printf("发送到到客户端的信息:%s", string(message))

			c.Socket.WriteMessage(websocket.BinaryMessage, message)
		}
	}
}

//Ws Handler socket 连接 中间件 作用:升级协议,用户验证,自定义信息等
func WsHandler(c *gin.Context) {
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		myerr.SendHttpErr(c)
		return
	}
	// postman 使用设置username
	username := c.Request.Header.Get("username")
	client := &Client{
		ID:       utils.GetUID(),
		Socket:   conn,
		Send:     make(chan []byte),
		UserName: username,
	}
	Manager.Register <- client
	go client.Read()
	go client.Write()
}
