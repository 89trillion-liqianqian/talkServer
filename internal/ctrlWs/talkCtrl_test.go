package ctrlWs

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"net/url"
	"talkServer/protocol/protobuf"
	"testing"
)

func TestLogin(t *testing.T) {
	addr := "127.0.0.1:8000"
	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}
	var dialer *websocket.Dialer
	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	sendData := protobuf.TaskInfo{
		Code: 1,
		Msg:  "我是 001，登陆",
		Name: "001",
		Cmd:  LoginCmd,
	}
	bData, _ := proto.Marshal(&sendData)
	conn.WriteMessage(websocket.BinaryMessage, bData)

	_, message, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("log"+
			" 失败:", err)
		return
	}

	fmt.Printf("login ok : %s\n", message)
}
