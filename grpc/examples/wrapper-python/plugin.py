#!/usr/bin/env python
# coding:utf-8
"""
@author: nivic ybyang7
@license: Apache Licence
@file: server
@time: 2022/10/28
@contact: ybyang7@iflytek.com
@site:
@software: PyCharm

# code is far away from bugs with the god animal protecting
    I love animals. They taste delicious.
              ┏┓      ┏┓
            ┏┛┻━━━┛┻┓
            ┃      ☃      ┃
            ┃  ┳┛  ┗┳  ┃
            ┃      ┻      ┃
            ┗━┓      ┏━┛
                ┃      ┗━━━┓
                ┃  神兽保佑    ┣┓
                ┃　永无BUG！   ┏┛
                ┗┓┓┏━┳┓┏┛
                  ┃┫┫  ┃┫┫
                  ┗┻┛  ┗┻┛
"""

import datetime
import importlib
import logging
import sys
import time
#  Copyright (c) 2022. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
#  Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
#  Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
#  Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
#  Vestibulum commodo. Ut rhoncus gravida arcu.
from concurrent import futures
from io import StringIO
from logging.handlers import QueueHandler, QueueListener
# from queue import Queue
from multiprocessing import Queue

import grpc
from aiges.aiges_inner import aiges_inner_pb2
from aiges.aiges_inner import aiges_inner_pb2_grpc
from aiges.aiges_inner import grpc_stdio_pb2
from aiges.aiges_inner import grpc_stdio_pb2_grpc
from aiges.dto import DataListCls
from aiges.errors import *
from aiges.utils.log import getLogger
from grpc_health.v1 import health_pb2, health_pb2_grpc
from grpc_health.v1.health import HealthServicer

log = getLogger(fmt=" %(name)s:%(funcName)s:%(lineno)s - %(levelname)s:  %(message)s", name="wrapper")
wrapper_module = "wrapper"
wrapper_class = "Wrapper"


class StdioService(grpc_stdio_pb2_grpc.GRPCStdioServicer):
    def __init__(self, log):
        self.log = log

    def StreamStdio(self, request, context):
        while True:
            sd = grpc_stdio_pb2.StdioData(channel=1, data=self.log.read())
            yield sd


class Logger:
    def __init__(self):
        self.stream = StringIO()  #
        que = Queue(-1)  # no limit on size
        self.queue_handler = QueueHandler(que)
        self.handler = logging.StreamHandler()
        self.listener = QueueListener(que, self.handler)
        self.log = logging.getLogger('python-plugin')
        self.log.setLevel(logging.INFO)
        self.logFormatter = logging.Formatter('%(asctime)s %(levelname)s  %(name)s %(pathname)s:%(lineno)d - %('
                                              'message)s')
        # self.handler.setFormatter(self.logFormatter)
        for handler in self.log.handlers:
            self.log.removeHandler(handler)
        self.log.addHandler(self.handler)
        self.listener.start()

    def __del__(self):
        pass
        # self.listener.stop()

    def read(self):
        self.handler.flush()
        ret = self.logFormatter.format(self.listener.queue.get()) + "\n"
        return ret.encode("utf-8")


class WrapperServiceServicer(aiges_inner_pb2_grpc.WrapperServiceServicer):
    """Provides methods that implement functionality of route guide server."""

    def __init__(self, q):
        self.response_queue = q
        self.count = 0
        self.userWrapperObject = None
        pass

    def wrapperInit(self, request, context):
        log.info("Importing module from wrapper.py: %s", wrapper_module)
        try:
            interface_file = importlib.import_module(wrapper_module)
            user_wrapper_cls = getattr(interface_file, wrapper_class)
            self.userWrapperObject = user_wrapper_cls()
            log.info("User Wrapper newed Success.. starting call user init functions...")
            ret = self.userWrapperObject.wrapperInit(request.config)
            if ret != 0:
                log.error("User wrapperInit function failed.. ret: %s" % str(ret))
                return aiges_inner_pb2.Ret(ret=USER_INIT_ERROR)

        except Exception as e:
            log.error(e)
            ret = INIT_ERROR
            return aiges_inner_pb2.Ret(ret=ret)

        return aiges_inner_pb2.Ret(ret=OK)

    def wrapperOnceExec(self, request, context):
        if not self.userWrapperObject:
            return aiges_inner_pb2.Response(ret=USER_EXEC_ERROR)
        self.count += 1
        user_resp = self.userWrapperObject.wrapperOnceExec(request.params, self.convertPbReq2Req(request))
        if not user_resp or not user_resp.list:
            return aiges_inner_pb2.Response(ret=USER_EXEC_ERROR)
        d_list = []
        for ur in user_resp.list:
            d = aiges_inner_pb2.ResponseData(key=ur.key, data=ur.data, len=ur.len, status=ur.status)
            d_list.append(d)
        r = aiges_inner_pb2.Response(list=d_list, tag=request.tag)
        call_back(self.response_queue, r)
        return aiges_inner_pb2.Response(list=[])

    def convertPbReq2Req(self, req):
        r = DataListCls()
        r.list = req.list
        return r

    def testStream(self, request_iterator, context):
        prev_notes = []
        for new_note in request_iterator:
            print(new_note.data)
            yield aiges_inner_pb2.Response(list=[])
            prev_notes.append(new_note)

    def communicate(self, request_iterator, context):
        # 这里无需双向似乎，如有必要，需要在加载器中传入相关信息
        while True:
            data = self.response_queue.get()
            yield data


def call_back(response_queue, r):
    response_queue.put(r)


def send_to_queue(q):
    x = 0
    while True:
        x += 1
        time.sleep(1)
        # print("sending... {}".format(x))
        msg = "count: {} . now : {}".format(x, datetime.datetime.now())
        d = aiges_inner_pb2.ResponseData(key=str(x), data=msg.encode("utf-8"), len=x, status=3)
        r = aiges_inner_pb2.Response(list=[d])
        # q.put(r)


def serve():
    work_q = Queue()
    # w = threading.Thread(target=send_to_queue, args=(work_q,))
    # w.start()

    # We need to build a health service to work with go-plugin
    health = HealthServicer()
    health.set("plugin", health_pb2.HealthCheckResponse.ServingStatus.Value('SERVING'))
    # Start the server.
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    aiges_inner_pb2_grpc.add_WrapperServiceServicer_to_server(
        WrapperServiceServicer(work_q), server)
    # add stdio service
    # 这里没有必要，因为go-plugin似乎已经捕捉了 标准输出
    # grpc_stdio_pb2_grpc.add_GRPCStdioServicer_to_server(StdioService(logger), server)

    health_pb2_grpc.add_HealthServicer_to_server(health, server)

    server.add_insecure_port('[::]:50055')
    server.start()
    # Output information
    print("1|1|tcp|127.0.0.1:50055|grpc")
    sys.stdout.flush()

    server.wait_for_termination()


def run():
    logging.basicConfig()
    serve()


if __name__ == '__main__':
    run()
