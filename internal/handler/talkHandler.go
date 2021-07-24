package handler

import (
	"github.com/golang/protobuf/proto"
	"log"
	"talkServer/protocol/protobuf"
)

// 获取下线数据
func GetExitData(userName string) (bData []byte, err error) {
	sendData := protobuf.TaskInfo{
		Code: 1,
		Msg:  "我是" + userName + "，下线了",
		Name: userName,
	}
	bData, err = proto.Marshal(&sendData)
	if err != nil {
		log.Println(err)
	}
	return
}

// 获取玩家列表数据
func GetUserListData(userList, cmd string) (bData []byte, err error) {
	sendData := protobuf.TaskInfo{
		Code: 1,
		Msg:  userList,
		Cmd:  cmd,
	}
	bData, err = proto.Marshal(&sendData)
	if err != nil {
		log.Println(err)
	}
	return
}

// 获取ping数据
func GetPingData() (bData []byte, err error) {
	sendData := protobuf.PingInfo{
		Msg: "pong",
	}
	bData, err = proto.Marshal(&sendData)
	if err != nil {
		log.Println(err)
	}
	return
}
