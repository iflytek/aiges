package com.iflytek.ccr.polaris.companion.utils;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.companion.cache.consumer.ConfigFeedbackQueueConsumer;
import com.iflytek.ccr.polaris.companion.cache.consumer.ServiceFeedbackQueueConsumer;
import com.iflytek.ccr.polaris.companion.common.Constants;
import com.iflytek.ccr.polaris.companion.common.QueueItemSerializer;
import org.apache.curator.framework.recipes.queue.DistributedQueue;
import org.apache.curator.framework.recipes.queue.QueueSerializer;
import org.apache.zookeeper.data.Stat;

import java.util.Iterator;
import java.util.Map;
import java.util.Set;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicInteger;

public class DistributedQueueManager {
    private static final EasyLogger logger = EasyLoggerFactory.getInstance(DistributedQueueManager.class);
    private static ConcurrentHashMap<String, DistributedQueue> configDistributedQueueMap = new ConcurrentHashMap<>();
    private static ConcurrentHashMap<String, DistributedQueue> serviceDistributedQueueMap = new ConcurrentHashMap<>();
    private static AtomicInteger currentConfigIndex = new AtomicInteger(0);
    private static AtomicInteger currentServiceIndex = new AtomicInteger(0);

    public static synchronized void clearQueueMap() {
        configDistributedQueueMap.clear();
        serviceDistributedQueueMap.clear();
    }

    public static synchronized DistributedQueue getConfigDistributedQueue() {
        ZkHelper zkHelper = ZkInstanceUtil.getInstance();
        DistributedQueue resultDistributedQueue = null;
        if (configDistributedQueueMap.isEmpty()) {
            try {
                QueueSerializer<String> serializer = new QueueItemSerializer();
                String queuePath = Constants.QUEUE_PATH_CONFIG + "/" + currentConfigIndex.get();
                if (!zkHelper.checkExists(queuePath)) {
                    zkHelper.addPersistent(queuePath, "");
                }
                DistributedQueue distributedQueue = zkHelper.getDistributedQueue(new ConfigFeedbackQueueConsumer(queuePath), serializer, queuePath);
                configDistributedQueueMap.put(queuePath, distributedQueue);
                resultDistributedQueue = distributedQueue;
            } catch (Exception e) {
                logger.error("", e);
            }
        } else {
            Set<Map.Entry<String, DistributedQueue>> set = configDistributedQueueMap.entrySet();
            Iterator<Map.Entry<String, DistributedQueue>> iterator = set.iterator();
            while (iterator.hasNext()) {
                Map.Entry<String, DistributedQueue> entry = iterator.next();
                String queuePath = entry.getKey();
                Stat stat = zkHelper.update(queuePath, "");
                if (stat.getNumChildren() < Constants.QUEUE_MAX_NUM) {
                    resultDistributedQueue = entry.getValue();
                    break;
                }
            }
            if (null == resultDistributedQueue) {
                String queuePath = Constants.QUEUE_PATH_CONFIG + "/" + currentConfigIndex.addAndGet(1);
                if (!zkHelper.checkExists(queuePath)) {
                    zkHelper.addPersistent(queuePath, "");
                }
                QueueSerializer<String> serializer = new QueueItemSerializer();
                try {
                    DistributedQueue distributedQueue = zkHelper.getDistributedQueue(new ConfigFeedbackQueueConsumer(queuePath), serializer, queuePath);
                    configDistributedQueueMap.put(queuePath, distributedQueue);
                    resultDistributedQueue = distributedQueue;
                } catch (Exception e) {
                    logger.error("", e);
                }

            }
        }
        return resultDistributedQueue;
    }

    public static synchronized DistributedQueue getServiceDistributedQueue() {
        ZkHelper zkHelper = ZkInstanceUtil.getInstance();
        DistributedQueue resultDistributedQueue = null;
        if (serviceDistributedQueueMap.isEmpty()) {
            try {
                QueueSerializer<String> serializer = new QueueItemSerializer();
                String queuePath = Constants.QUEUE_PATH_SERVICE + "/" + currentServiceIndex.get();
                if (!zkHelper.checkExists(queuePath)) {
                    zkHelper.addPersistent(queuePath, "");
                }
                DistributedQueue distributedQueue = zkHelper.getDistributedQueue(new ServiceFeedbackQueueConsumer(queuePath), serializer, queuePath);
                serviceDistributedQueueMap.put(queuePath, distributedQueue);
                resultDistributedQueue = distributedQueue;
            } catch (Exception e) {
                logger.error("", e);
            }
        } else {
            Set<Map.Entry<String, DistributedQueue>> set = serviceDistributedQueueMap.entrySet();
            Iterator<Map.Entry<String, DistributedQueue>> iterator = set.iterator();
            while (iterator.hasNext()) {
                Map.Entry<String, DistributedQueue> entry = iterator.next();
                String queuePath = entry.getKey();
                Stat stat = zkHelper.update(queuePath, "");
                if (stat.getNumChildren() < Constants.QUEUE_MAX_NUM) {
                    resultDistributedQueue = entry.getValue();
                    break;
                }
            }
            if (null == resultDistributedQueue) {
                String queuePath = Constants.QUEUE_PATH_SERVICE + "/" + currentServiceIndex.addAndGet(1);
                if (!zkHelper.checkExists(queuePath)) {
                    zkHelper.addPersistent(queuePath, "");
                }
                QueueSerializer<String> serializer = new QueueItemSerializer();
                try {
                    DistributedQueue distributedQueue = zkHelper.getDistributedQueue(new ServiceFeedbackQueueConsumer(queuePath), serializer, queuePath);
                    serviceDistributedQueueMap.put(queuePath, distributedQueue);
                    resultDistributedQueue = distributedQueue;
                } catch (Exception e) {
                    logger.error("", e);
                }

            }
        }
        return resultDistributedQueue;
    }
}
