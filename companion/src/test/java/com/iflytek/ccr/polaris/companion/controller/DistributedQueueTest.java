package com.iflytek.ccr.polaris.companion.controller;

import com.iflytek.ccr.polaris.companion.common.Constants;
import com.iflytek.ccr.polaris.companion.common.QueueItemSerializer;
import com.iflytek.ccr.zkutil.ZkHelper;
import org.apache.curator.framework.CuratorFramework;
import org.apache.curator.framework.recipes.queue.DistributedQueue;
import org.apache.curator.framework.recipes.queue.QueueConsumer;
import org.apache.curator.framework.recipes.queue.QueueSerializer;
import org.apache.curator.framework.state.ConnectionState;
import org.junit.Test;

public class DistributedQueueTest {
    @Test
    public void test() {

        QueueConsumer consumer = new QueueConsumer<String>() {
            @Override
            public void stateChanged(CuratorFramework client, ConnectionState newState) {
                System.out.println(newState);
            }

            @Override
            public void consumeMessage(String message) throws Exception {
                Thread.sleep(10000);
                System.out.println(message);
            }
        };

        try {
            ZkHelper zkHelper = new ZkHelper("10.1.86.73:2181,10.1.86.74:2181,10.1.86.78:2181");
            QueueSerializer<String> serializer = new QueueItemSerializer();
            DistributedQueue distributedQueue = zkHelper.getDistributedQueue(consumer, serializer, Constants.QUEUE_PATH_SERVICE);
            distributedQueue.put("distributedQueuea");
            distributedQueue.put("distributedQueueb");
            distributedQueue.put("distributedQueuec");
            distributedQueue.put("distributedQueued");
            distributedQueue.put("assdssd");
            Thread.sleep(100000);
        } catch (Exception e) {
            e.printStackTrace();
        }

    }
}
