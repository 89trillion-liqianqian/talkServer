package ws

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"log"
	"talkServer/internal/handler"
	"talkServer/protocol/protobuf"
)

// 定义命令
const (
	LoginCmd    = "login"    // 绑定 username
	PingCmd     = "ping"     // ping/pong
	TalkCmd     = "talk"     // 聊天
	UserListCmd = "userlist" // 获取用户列表
	ExitCmd     = "exit"     // 退出
)

// Client is a websocket client
type Client struct {
	ID       string
	UserName string
	Socket   *websocket.Conn
	Send     chan []byte
	Manager  *ClientManager
}

// 监听读数据
func (c *Client) Read() {
	defer func() {
		c.Manager.Unregister <- c
		c.Socket.Close()
	}()

	for {
		c.Socket.PongHandler()
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			c.Manager.Unregister <- c
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
				delete(c.Manager.Clients, c.ID)
			}
			continue
		}
		// 获取玩家列表 userlist
		if receivedData.Cmd == UserListCmd {
			// 获取用户列表
			userList := getNamelist(c.Manager)
			bData, _ := handler.GetUserListData(userList, UserListCmd)
			select {
			case c.Send <- bData:
			default:
				close(c.Send)
				delete(c.Manager.Clients, c.ID)
			}
			continue
		}
		// 退出
		if receivedData.Cmd == ExitCmd && receivedData.Name == c.UserName {
			c.Manager.Unregister <- c
			c.Socket.Close()
			// 广播其他人
			c.Manager.Broadcast <- message
			break
		}

		c.Manager.Broadcast <- message
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

// 获取在线用户列表
func getNamelist(Manager *ClientManager) string {
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
