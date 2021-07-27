package ws

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"log"
	"talkServer/internal/handler"
	"talkServer/protocol/protobuf"
)

// Message is return msg
type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

// ClientManager is a websocket manager
type ClientManager struct {
	Clients    map[string]*Client
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

// 创建一个新的client管理
func NewManager() *ClientManager {
	var Manager = ClientManager{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[string]*Client),
	}

	return &Manager
}

// Start is  项目运行前, 协程开启start -> go Manager.Start()
func (m *ClientManager) Start() {
	for {
		log.Println("<---管道通信--->")
		select {
		case conn := <-m.Register:
			log.Printf("新用户加入:%v", conn.ID)
			m.Clients[conn.ID] = conn
			jsonMessage, _ := json.Marshal(&Message{Content: "Successful connection to socket service"})
			conn.Send <- jsonMessage
		case conn := <-m.Unregister:
			log.Printf("用户离开:%v", conn.ID)
			bData, _ := handler.GetExitData(conn.UserName)
			if _, ok := m.Clients[conn.ID]; ok {
				conn.Send <- bData
				close(conn.Send)
				delete(m.Clients, conn.ID)
			}

		case message := <-m.Broadcast:
			MessageStruct := protobuf.TaskInfo{}
			// 解析数据，广播其他玩家
			proto.Unmarshal(message, &MessageStruct)
			for _, conn := range m.Clients {
				if conn.UserName == MessageStruct.Name || conn.UserName == "" {
					continue
				}
				select {
				case conn.Send <- message:
				default:
					close(conn.Send)
					delete(m.Clients, conn.ID)
				}
			}
		}
	}
}
