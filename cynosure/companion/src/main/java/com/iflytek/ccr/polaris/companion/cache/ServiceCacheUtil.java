package com.iflytek.ccr.polaris.companion.cache;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.companion.common.ServiceFeedBackValue;
import com.iflytek.ccr.polaris.companion.utils.DistributedQueueManager;
import com.iflytek.ccr.polaris.companion.utils.JacksonUtils;

public class ServiceCacheUtil {

    private static final EasyLogger logger = EasyLoggerFactory.getInstance(ServiceCacheUtil.class);

//    private static DistributedQueue distributedQueue = null;
//    private static ZkHelper zkHelper = ZkInstanceUtil.getInstance();
//    private static QueueConsumer<String> consumer = null;
//    private ExecutorService executorService = Executors.newFixedThreadPool(Constants.NUM_3);
//    static {
//        //消费者处理方法
//        consumer = new QueueConsumer<String>() {
//            @Override
//            public void consumeMessage(String message) throws Exception {
//                ServiceFeedBackValue serviceValue = jsonUtils.toObject(message, ServiceFeedBackValue.class);
//                String url = Program.CONFIG_VALUE.getWebsiteUrl() + Constants.PUSH_CONFIG_SERVICE_SITE_URI;
//                logger.info("consumeMessage:"+ url);
//                String reqJson = jsonUtils.toJson(serviceValue);
//                StringEntity entity = new StringEntity(reqJson);
//                try {
//                    String result = Request.Post(url).setHeader("Content-type", "application/json")
//                            .body(entity).socketTimeout(5000)
//                            .execute().returnContent().asString();
//                    WebsitResult response = jsonUtils.toObject(result, WebsitResult.class);
//                    if (!Constants.SUCCESS.equals(response.getCode())) {
//                        logger.error("request:" + reqJson + ",result:" + result);
//                    }
//                } catch (Exception e) {
//                    logger.error(e);
//                    ServiceCacheUtil.getInstance().add(serviceValue);
//                }
//
//            }
//
//            @Override
//            public void stateChanged(CuratorFramework client, ConnectionState newState) {
//                logger.warn(String.format("new state:s%", newState));
//                System.out.println("new state: " + newState);
//            }
//        };
//
//        QueueSerializer<String> serializer = new QueueItemSerializer();
//        try {
//            if(!zkHelper.checkExists(Constants.QUEUE_PATH_SERVICE)){
//                zkHelper.addPersistent(Constants.QUEUE_PATH_SERVICE,"");
//            }
//            distributedQueue = zkHelper.getDistributedQueue(consumer, serializer, Constants.QUEUE_PATH_SERVICE);
//        } catch (Exception e) {
//            logger.error(e);
//        }
//        ServiceCacheUtil.getInstance();
//    }


    private ServiceCacheUtil() {
//        TreeCacheListener treeCacheListener = new TreeCacheListener() {
//
//            @Override
//            public void childEvent(CuratorFramework client, TreeCacheEvent event) throws Exception {
//                switch (event.getType()) {
//                    case CONNECTION_SUSPENDED:
//                        break;
//                    case CONNECTION_RECONNECTED:
//                        try {
//                            executorService.shutdownNow();
//                            executorService = Executors.newFixedThreadPool(Constants.NUM_3);
//                            for (int i = 0; i < Constants.NUM_6; i++) {
//                                executorService.submit(new ServiceFeedbackQueueConsumerTask());
//                            }
//                        } catch (Exception e) {
//                            logger.error(e);
//                        }
//                        break;
//                    case CONNECTION_LOST:
//                        Thread.sleep(1000);
//                        executorService.shutdownNow();
//                        break;
//                    case INITIALIZED:
//                        executorService = Executors.newFixedThreadPool(Constants.NUM_3);
//                        for (int i = 0; i < Constants.NUM_3; i++) {
//                            executorService.submit(new ServiceFeedbackQueueConsumerTask());
//                        }
//                }}
//        };
//        zkHelper.addListener(treeCacheListener, Constants.QUEUE_PATH_CONFIG);
    }

    public static final ServiceCacheUtil getInstance() {
        return CacheUtilHolder.INSTANCE;
    }

    public boolean add(ServiceFeedBackValue serviceValue) {
        boolean flag = true;
        try {
//            if (null == distributedQueue) {
//                QueueSerializer<String> serializer = new QueueItemSerializer();
//                distributedQueue = zkHelper.getDistributedQueue(null, serializer, Constants.QUEUE_PATH_SERVICE);
//            }
            DistributedQueueManager.getServiceDistributedQueue().put(JacksonUtils.toJson(serviceValue));
        } catch (Exception e) {
            logger.error(e);
            try{
                DistributedQueueManager.clearQueueMap();
                DistributedQueueManager.getServiceDistributedQueue().put(JacksonUtils.toJson(serviceValue));
            }catch (Exception ee){
                LocalCacheUtil.feedBackQueue.add(serviceValue);
                flag = false;
                logger.error(ee);
            }
        }
        return flag;
    }

    private static class CacheUtilHolder {
        private static final ServiceCacheUtil INSTANCE = new ServiceCacheUtil();
    }
}
