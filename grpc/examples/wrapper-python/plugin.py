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

import logging
import sys
#  Copyright (c) 2022. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
#  Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
#  Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
#  Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
#  Vestibulum commodo. Ut rhoncus gravida arcu.
from concurrent import futures
from io import StringIO
from logging.handlers import QueueHandler, QueueListener
from queue import Queue

import grpc
from aiges.aiges_inner import aiges_inner_pb2
from aiges.aiges_inner import aiges_inner_pb2_grpc
from aiges.aiges_inner import grpc_stdio_pb2
from aiges.aiges_inner import grpc_stdio_pb2_grpc
from grpc_health.v1 import health_pb2, health_pb2_grpc
from grpc_health.v1.health import HealthServicer


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
        self.log.setLevel(logging.DEBUG)
        self.logFormatter = logging.Formatter('%(asctime)s %(levelname)s  %(name)s %(pathname)s:%(lineno)d - %('
                                              'message)s')
        self.handler.setFormatter(self.logFormatter)
        for handler in self.log.handlers:
            self.log.removeHandler(handler)
        self.log.addHandler(self.queue_handler)
        self.listener.start()

    def __del__(self):
        self.listener.stop()

    def read(self):
        self.handler.flush()
        ret = self.logFormatter.format(self.listener.queue.get()) + "\n"
        return ret.encode("utf-8")


logger = Logger()
log = logger.log


class WrapperServiceServicer(aiges_inner_pb2_grpc.WrapperServiceServicer):
    """Provides methods that implement functionality of route guide server."""

    def __init__(self):
        pass

    def wrapperInit(self, request, context):
        log.info(request)
        log.error(2222222222222)
        return aiges_inner_pb2.Ret(ret=1)

    def wrapperOnceExec(self, request, context):
        print(request)
        return aiges_inner_pb2.Response(list=[])
        pass

    def testStream(self, request_iterator, context):
        prev_notes = []
        for new_note in request_iterator:
            print(new_note.data)
            yield aiges_inner_pb2.Response(list=[])
            prev_notes.append(new_note)


def serve():
    # We need to build a health service to work with go-plugin
    health = HealthServicer()
    health.set("plugin", health_pb2.HealthCheckResponse.ServingStatus.Value('SERVING'))
    # Start the server.
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    aiges_inner_pb2_grpc.add_WrapperServiceServicer_to_server(
        WrapperServiceServicer(), server)
    # add stdio service
    grpc_stdio_pb2_grpc.add_GRPCStdioServicer_to_server(StdioService(logger), server)

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
