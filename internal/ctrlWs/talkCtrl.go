package ctrlWs

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"talkServer/protocol/protobuf"
	"talkServer/utils"
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
			sendData := protobuf.TaskInfo{
				Code: 1,
				Msg:  "我是" + conn.UserName + "，下线了",
				Name: conn.UserName,
			}
			bData, _ := proto.Marshal(&sendData)
			if _, ok := Manager.Clients[conn.ID]; ok {
				conn.Send <- bData
				close(conn.Send)
				delete(Manager.Clients, conn.ID)
			}

		case message := <-Manager.Broadcast:
			MessageStruct := protobuf.TaskInfo{}
			proto.Unmarshal(message, &MessageStruct)
			//json.Unmarshal(message, &MessageStruct)
			for id, conn := range Manager.Clients {
				log.Println("---conn.ID", conn.ID, id, conn.UserName, MessageStruct.Name)
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
		if receivedData.Cmd == "login" {
			c.UserName = receivedData.Name
			continue
		}
		// 退出
		if receivedData.Cmd == "exit" && receivedData.Name == c.UserName {
			Manager.Unregister <- c
			c.Socket.Close()
			// 广播其他人
			Manager.Broadcast <- message
			break
		}
		// 获取玩家列表 userlist
		if receivedData.Cmd == "userlist" {
			// 获取用户列表
			userList := getNamelist()
			sendData := protobuf.TaskInfo{
				Code: 1,
				Msg:  userList,
				Cmd:  "userlist",
			}
			bData, _ := proto.Marshal(&sendData)
			select {
			case c.Send <- bData:
			default:
				close(c.Send)
				delete(Manager.Clients, c.ID)
			}
			continue
		}

		Manager.Broadcast <- message
	}
}

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

//TestHandler socket 连接 中间件 作用:升级协议,用户验证,自定义信息等
func WsHandler(c *gin.Context) {
	uid := c.Query("uid")
	username := c.Request.Header["username"]
	log.Println("---username", uid, username, c.Request.Header)
	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	//可以添加用户信息验证
	client := &Client{
		ID:     utils.GetUID(),
		Socket: conn,
		Send:   make(chan []byte),
	}
	Manager.Register <- client
	go client.Read()
	go client.Write()
}
