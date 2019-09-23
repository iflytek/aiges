package com.iflytek.ccr.polaris.companion.cache;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.companion.common.ConfigFeedBackValue;
import com.iflytek.ccr.polaris.companion.utils.DistributedQueueManager;
import com.iflytek.ccr.polaris.companion.utils.JacksonUtils;
import com.iflytek.ccr.polaris.companion.utils.ZkInstanceUtil;
import com.iflytek.ccr.zkutil.ZkHelper;
import org.apache.curator.framework.recipes.queue.DistributedQueue;

public class FeedbackCacheUtil {

    private static final EasyLogger logger = EasyLoggerFactory.getInstance(FeedbackCacheUtil.class);

    private static DistributedQueue distributedQueue = null;
    private static ZkHelper zkHelper = ZkInstanceUtil.getInstance();
//    ExecutorService executorService = Executors.newFixedThreadPool(Constants.NUM_3);

    private FeedbackCacheUtil() {

//        TreeCacheListener treeCacheListener = new TreeCacheListener() {
//
//            @Override
//            public void childEvent(CuratorFramework client, TreeCacheEvent event) throws Exception {
//                switch (event.getType()) {
//                    case CONNECTION_SUSPENDED:
//                    case CONNECTION_RECONNECTED:
//                        try {
//                            executorService.shutdownNow();
//                            executorService = Executors.newFixedThreadPool(Constants.NUM_3);
//                            for (int i = 0; i < Constants.NUM_6; i++) {
//                                executorService.submit(new FeedbackQueueConsumerMonitorTask());
//                            }
//                        } catch (Exception e) {
//                            logger.error(e);
//                        }
//                        break;
//                    case CONNECTION_LOST:
//                        executorService.shutdownNow();
//                        break;
//                    case INITIALIZED:
//                        executorService = Executors.newFixedThreadPool(Constants.NUM_3);
//                        for (int i = 0; i < Constants.NUM_3; i++) {
//                            executorService.submit(new FeedbackQueueConsumerMonitorTask());
//                        }
//                }
//            }
//        };
//        zkHelper.addListener(treeCacheListener, Constants.QUEUE_PATH_CONFIG);
    }

    public static final FeedbackCacheUtil getInstance() {
        return CacheUtilHolder.INSTANCE;
    }

    public boolean add(ConfigFeedBackValue feedBackValue) {
        boolean flag = true;
        try {
            DistributedQueueManager.getConfigDistributedQueue().put(JacksonUtils.toJson(feedBackValue));
        } catch (Exception e) {
            logger.error(e);
            try{
                DistributedQueueManager.clearQueueMap();
                DistributedQueueManager.getConfigDistributedQueue().put(JacksonUtils.toJson(feedBackValue));
            }catch (Exception ee){
                LocalCacheUtil.feedBackQueue.add(feedBackValue);
                flag = false;
                logger.error(ee);
            }
        }
        return flag;
    }

    private static class CacheUtilHolder {
        private static final FeedbackCacheUtil INSTANCE = new FeedbackCacheUtil();
    }
}
