## 1. 整体框架
talkServer 是Websocket + protobuf 实现聊天服务,题5，

主要功能包括：接受客户端发起的连接请求，建立ws连接（header中带username），支持ping/pong心跳

2） 维护管理多用户的ws连接

3） 收到talk类型消息，记录日志，把收到消息转发给所有连接的用户

4） 收到exit类型消息，断开连接，清理用户信息

5） 收到userlist消息，返回所有在线用户username

## 2.目录结构

文件夹结构和各文件主要职责

```
liqianqian@liqianqian talkServer % tree
.
├── README.md
├── app
│   ├── main.go									#代码入口
│   └── ws
│       └── socketserver.go			#ws 服务启动
├── conf
├── go.mod
├── go.sum
├── internal
│   ├── ctrlWs									#控制层，聊天
│   │   ├── talkCtrl.go
│   │   └── talkCtrlz_test.go
│   ├── handler									#聊天handler
│   │   └── talkHandler.go
│   ├── myerr
│   │   └── err.go							#错误返回
│   └── router									#聊天路由
│       └── router.go
├── load.py											#压测脚本
├── locust											#压测报告
│   └── report_1627121507.6245902.html
├── pkg
├── protocol									#proto 文件
│   ├── protobuf
│   │   └── talk.pb.go
│   └── talk.proto
│   └── talk_pb2.py						#python proto生成文件
├── test											#聊天client 
│   └── test.go
└── utils											#工具
│    └── tool.go
└── websocket?\201\212天?\201?\213?\233?.jpg   #流程图
14 directories, 17 files
liqianqian@liqianqian talkServer % 
```

## 3.逻辑代码分层

|    层     | 文件夹                           | 主要职责                                                     | 调用关系                  | 其它说明     |
| :-------: | :------------------------------- | ------------------------------------------------------------ | ------------------------- | ------------ |
|  应用层   | /app/ws/socketserver.go          | Ws 服务器启动                                                | 调用路由层                | 不可同层调用 |
|  路由层   | /internal/router/router.go       | 路由转发                                                     | 被应用层调用，调用控制层  | 不可同层调用 |
|  控制层   | /internal/ctrlWs/talkCtrl,go     | Ws client连接管理，响应,talk聊天，userlist 用户列表，exit,ping | 被路由层调用，调用handler | 不可同层调用 |
| handler层 | /internal/handler/talkHandler.go | 处理具体业务                                                 | 被控制层调用              | 不可同层调   |
|           |                                  |                                                              |                           |              |
| 压力测试  | Locust/load.py                   | 进行压力测试                                                 | 无调用关系                | 不可同层调用 |

## 4.存储设计

无存储

## 5.接口设计供客户端调用的接口

5.1client连接登陆

| 信息     | 说明              |
| -------- | ----------------- |
| 接口方式 | websocket         |
| 事件名称 | login             |
| 请求消息 | Request\TaskInfo  |
| 响应消息 | Response\TaskInfo |

请求消息中定义

| Proto 字段 | 说明                                                |
| ---------- | --------------------------------------------------- |
| code       | 响应码，1成功                                       |
| msg        | 聊天信息                                            |
| name       | 发送人name                                          |
| cmd        | cmd命令标示，talk聊天、userlist玩家列表、  exit退出 |

错误码

| 字段 | 类型 | 说明     |
| ---- | ---- | -------- |
| 0    | 通用 | 响应错误 |

## 6.第三方库

gin

```
用于api服务，go web 框架
代码： github.com/gin-gonic/gin

```

proto

```
用于消息数据协议
包含：proto.Unmarshal，proto.Marshal 数据序列化
代码："github.com/golang/protobuf/proto"

```

websocket

```
用于建立socket 长连接
代码："github.com/gorilla/websocket"
```

## 7.如何编译执行

```
#切换主目录下
cd ./app/
#编译
go build
```

## 8.todo 

```
后续优化，连接验证
```





