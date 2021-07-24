package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"strings"
	"talkServer/protocol/protobuf"
)

var addr = flag.String("addr", "127.0.0.1:8000", "http service address")
var username = flag.String("username", "001", "this is user")

func timeWriter(conn *websocket.Conn) {
	// 聊天
	sendData := protobuf.TaskInfo{
		Code: 1,
		Msg:  "我是" + *username + "，哈喽！",
		Name: *username,
		Cmd:  "talk",
	}
	bData, _ := proto.Marshal(&sendData)
	conn.WriteMessage(websocket.BinaryMessage, bData)
}

func userList(conn *websocket.Conn) {
	//conn.WriteMessage(websocket.TextMessage, []byte("002_task_:我是002,发送时间："+time.Now().Format("2006-01-02 15:04:05")))
	// 获取在线用户列表
	sendData := protobuf.TaskInfo{
		Code: 1,
		Msg:  "我是" + *username + "，获取用户列表！",
		Name: *username,
		Cmd:  "userlist",
	}
	bData, _ := proto.Marshal(&sendData)
	conn.WriteMessage(websocket.BinaryMessage, bData)
}

func exit(conn *websocket.Conn) {
	//conn.WriteMessage(websocket.TextMessage, []byte("002_task_:我是002,发送时间："+time.Now().Format("2006-01-02 15:04:05")))
	// 发送退出
	sendData := protobuf.TaskInfo{
		Code: 1,
		Msg:  "我是" + *username + "，我要下线了！",
		Name: *username,
		Cmd:  "exit",
	}
	bData, _ := proto.Marshal(&sendData)
	conn.WriteMessage(websocket.BinaryMessage, bData)
}

func login(conn *websocket.Conn) {
	//conn.WriteMessage(websocket.TextMessage, []byte("002_task_:我是002,发送时间："+time.Now().Format("2006-01-02 15:04:05")))
	sendData := protobuf.TaskInfo{
		Code: 1,
		Msg:  "我是" + *username + "，来聊两句",
		Name: *username,
		Cmd:  "login",
	}
	bData, _ := proto.Marshal(&sendData)
	conn.WriteMessage(websocket.BinaryMessage, bData)
}

func ping(conn *websocket.Conn) {
	//conn.WriteMessage(websocket.TextMessage, []byte("002_task_:我是002,发送时间："+time.Now().Format("2006-01-02 15:04:05")))
	sendData := protobuf.TaskInfo{
		Code: 1,
		Msg:  "我是" + *username + "，ping",
		Name: *username,
		Cmd:  "ping",
	}
	bData, _ := proto.Marshal(&sendData)
	conn.WriteMessage(websocket.BinaryMessage, bData)
}

// 控制器
func action(conn *websocket.Conn) {
	print("请输入内容: ")
	for true {
		var talkStr string
		reader := bufio.NewReader(os.Stdin)
		talkStr, _ = reader.ReadString('\n')
		talkStr = strings.TrimSpace(talkStr)
		if talkStr == "login" {
			login(conn)
		} else if talkStr == "talk" {
			timeWriter(conn)
		} else if talkStr == "userlist" {
			userList(conn)
		} else if talkStr == "exit" {
			exit(conn)
		} else if talkStr == "ping" {
			ping(conn)
		}
	}
}

func main() {
	flag.Parse()
	log.Println("--start test")
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	var dialer *websocket.Dialer
	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	//login(conn)
	//go timeWriter(conn)
	go action(conn)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			return
		}

		fmt.Printf("received: %s\n", message)
	}
}
