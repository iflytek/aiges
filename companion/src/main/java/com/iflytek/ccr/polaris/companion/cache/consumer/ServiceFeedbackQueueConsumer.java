package com.iflytek.ccr.polaris.companion.cache.consumer;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.companion.cache.ServiceCacheUtil;
import com.iflytek.ccr.polaris.companion.common.Constants;
import com.iflytek.ccr.polaris.companion.common.QueueItemSerializer;
import com.iflytek.ccr.polaris.companion.common.ServiceFeedBackValue;
import com.iflytek.ccr.polaris.companion.common.WebsitResult;
import com.iflytek.ccr.polaris.companion.main.Program;
import com.iflytek.ccr.polaris.companion.utils.HttpsClientUtil;
import com.iflytek.ccr.polaris.companion.utils.JacksonUtils;
import com.iflytek.ccr.polaris.companion.utils.ZkHelper;
import com.iflytek.ccr.polaris.companion.utils.ZkInstanceUtil;
import org.apache.curator.framework.CuratorFramework;
import org.apache.curator.framework.recipes.queue.DistributedQueue;
import org.apache.curator.framework.recipes.queue.QueueConsumer;
import org.apache.curator.framework.recipes.queue.QueueSerializer;
import org.apache.curator.framework.state.ConnectionState;
import org.apache.http.client.fluent.Request;
import org.apache.http.entity.StringEntity;

public class ServiceFeedbackQueueConsumer implements QueueConsumer<String> {

    private static final EasyLogger logger = EasyLoggerFactory.getInstance(ServiceFeedbackQueueConsumer.class);

    private static ZkHelper zkHelper = ZkInstanceUtil.getInstance();

    private String queuePath;

    public ServiceFeedbackQueueConsumer(String queuePath) {
        this.queuePath = queuePath;
    }

    @Override
    public void consumeMessage(String message) throws Exception {
        long start = System.currentTimeMillis();
        ServiceFeedBackValue serviceValue = JacksonUtils.toObject(message, ServiceFeedBackValue.class);
        String url = Program.CONFIG_VALUE.getWebsiteUrl() + Constants.PUSH_CONFIG_SERVICE_SITE_URI;
        logger.info("consumeMessage:" + url);
        String reqJson = JacksonUtils.toJson(serviceValue);
        StringEntity entity = new StringEntity(reqJson);
        try {

            String result;
            if (url.startsWith("https")) {
                result = HttpsClientUtil.doPostByStringEntity(url, entity, Constants.DEFAULT_CHARSET);
            } else {
                result = Request.Post(url).setHeader("Content-type", "application/json")
                        .body(entity).socketTimeout(5000).connectTimeout(5000)
                        .execute().returnContent().asString();
            }
            WebsitResult response = JacksonUtils.toObject(result, WebsitResult.class);
            if (!Constants.SUCCESS.equals(response.getCode())) {
                logger.error("request:" + reqJson + ",result:" + result);
            }
            logger.info("ServiceFeedbackQueueConsumer consumeMessage cost:" + (System.currentTimeMillis() - start));
        } catch (Exception e) {
            logger.error(e);
            ServiceCacheUtil.getInstance().add(serviceValue);
        }

    }

    @Override
    public void stateChanged(CuratorFramework client, ConnectionState newState) {
        logger.error("new state: " + newState);
        if (ConnectionState.RECONNECTED.equals(newState) || ConnectionState.CONNECTED.equals(newState)) {
            QueueSerializer<String> serializer = new QueueItemSerializer();
            try {
                if (!zkHelper.checkExists(queuePath)) {
                    zkHelper.addPersistent(queuePath, "");
                }
                DistributedQueue distributedQueue = zkHelper.getDistributedQueue(this, serializer, queuePath);
            } catch (Exception e) {
                logger.error(e);
            }
        }
    }
}
