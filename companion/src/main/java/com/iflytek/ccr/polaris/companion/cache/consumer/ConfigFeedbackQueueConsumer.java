package com.iflytek.ccr.polaris.companion.cache.consumer;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.companion.cache.FeedbackCacheUtil;
import com.iflytek.ccr.polaris.companion.common.ConfigFeedBackValue;
import com.iflytek.ccr.polaris.companion.common.Constants;
import com.iflytek.ccr.polaris.companion.common.QueueItemSerializer;
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

public class ConfigFeedbackQueueConsumer implements QueueConsumer<String> {


    private static final EasyLogger logger = EasyLoggerFactory.getInstance(ConfigFeedbackQueueConsumer.class);

    private static ZkHelper zkHelper = ZkInstanceUtil.getInstance();

    private String queuePath;

    public ConfigFeedbackQueueConsumer(String queuePath) {
        this.queuePath = queuePath;
    }

    @Override
    public void consumeMessage(String message) throws Exception {
        long start = System.currentTimeMillis();
        ConfigFeedBackValue feedBackValue = JacksonUtils.toObject(message, ConfigFeedBackValue.class);
        String url = Program.CONFIG_VALUE.getWebsiteUrl() + Constants.PUSH_CONFIG_FEEDBACK_SITE_URI;
        logger.info("consumeMessage:" + url);
        String reqJson = JacksonUtils.toJson(feedBackValue);
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
            logger.info("ConfigFeedbackQueueConsumer consumeMessage cost:" + (System.currentTimeMillis() - start));
        } catch (Exception e) {
            logger.error(e);
            FeedbackCacheUtil.getInstance().add(feedBackValue);
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
