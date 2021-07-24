import json
import logging
import re
import threading
import time

import gevent
import websocket
from locust import User, task, events
from protocol import talk_pb2


class SocketIOUser(User):
    abstract = True

    def __init__(self, environment):
        # super().__init__(environment)
        self.ws = websocket.WebSocket()
        self.ws.connect("ws://localhost:8000/ws")
        datas = talk_pb2.TaskInfo()
        datas.code = 1
        datas.name = "qq"
        datas.msg = "登陆"
        datas.cmd = "login"
        dataB = datas.SerializeToString()
        self.ws.send(dataB)
        timer = threading.Timer(1, self.fun_heart)
        timer.start()

    # ping/pong
    def fun_heart(self):
        datas = talk_pb2.PingInfo()
        datas.msg = "ping"
        dataB = datas.SerializeToString()
        self.ws.send(json.dumps(dataB))
        timer = threading.Timer(3, self.fun_heart)
        timer.start()


class MySocketIOUser(SocketIOUser):

    @task
    def test_ws(self, send_flag=False):
        # seng_flag = False
        # self.connect(self.host)
        # 这里可以使用这一个task来完成发送订阅和处理订阅
        time.sleep(0.5)
        if not send_flag:
            datas = talk_pb2.TaskInfo()
            datas.code = 1
            datas.name = "qq"
            datas.msg = "来聊几句"
            datas.cmd = "talk"
            dataB = datas.SerializeToString()
            self.ws.send(dataB)
            dataR = self.ws.recv()
            result = talk_pb2.TaskInfo
            result.ParseFromString(dataR)
            print(result)
            send_flag=True
            events.request_success.fire(
                request_type="ws",
                name="send_entrust",
                response_time=100,
                response_length=300)
        else:
            flag = True
            # 循环接收
            while flag:
                # self.ws本身是个迭代器，next和调用resv()是一样的，都可以
                start_time = time.time()
                resv = next(self.ws)
                # 推送间隔就可以当做统计数据进行显示
                total_time = int((time.time() - start_time) * 1000)
                if resv != '':
                    events.request_success.fire(
                        request_type="ws",
                        name="resv_entrust",
                        response_time=total_time,
                        response_length=59)
                else:
                    # 如果中断了可以增加重连方案
                    events.request_failure.fire(
                        request_type="ws",
                        name="resv_entrust",
                        response_time=total_time,
                        exception=Exception(""),
                    )
                    flag = False
            else:
                print("---")

    if __name__ == "__main__":
        host = "ws://127.0.0.1:8000/ws"