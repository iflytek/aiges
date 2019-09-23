package com.iflytek.ccr.polaris.companion.task;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.companion.cache.consumer.ConfigFeedbackQueueConsumer;
import com.iflytek.ccr.polaris.companion.common.Constants;
import com.iflytek.ccr.polaris.companion.common.QueueItemSerializer;
import com.iflytek.ccr.polaris.companion.utils.ZkInstanceUtil;
import com.iflytek.ccr.zkutil.ZkHelper;
import org.apache.curator.framework.recipes.queue.QueueSerializer;
import org.apache.zookeeper.data.Stat;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class ConfigFeedbackQueueConsumerMonitorTask implements Runnable {

    private static final EasyLogger logger = EasyLoggerFactory.getInstance(ConfigFeedbackQueueConsumerMonitorTask.class);
    private Map<String, Integer> countMap = new HashMap<>();

    @Override
    public void run() {
        while (true) {
            try {
                Thread.sleep(Constants.MILLIS_1_MINUTES);
                ZkHelper zkHelper = ZkInstanceUtil.getInstance();
                //如果节点不存在，则创建
                if(!zkHelper.checkExists(Constants.QUEUE_PATH_CONFIG)){
                    zkHelper.addPersistent(Constants.QUEUE_PATH_CONFIG,"");
                }
                List<String> queuePathList = zkHelper.getChildren(Constants.QUEUE_PATH_CONFIG);
                if (!queuePathList.isEmpty()) {
                    for (String path : queuePathList) {
                        String queuePath = Constants.QUEUE_PATH_CONFIG + "/" + path;
                        Stat stat = zkHelper.update(queuePath, "");
                        if (countMap.containsKey(queuePath)) {
                            if (0 != stat.getNumChildren() && stat.getNumChildren() == countMap.get(queuePath)) {
                                QueueSerializer<String> serializer = new QueueItemSerializer();
                                try {
                                    zkHelper.getDistributedQueue(new ConfigFeedbackQueueConsumer(queuePath), serializer, queuePath);
                                } catch (Exception e) {
                                    logger.error("", e);
                                }
                            }
                        }
                        countMap.put(queuePath, stat.getNumChildren());
                    }
                }
            } catch (Exception e) {
                logger.error("", e);
            }
        }


    }
}
