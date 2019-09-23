package com.iflytek.ccr.polaris.companion.cache;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.companion.common.ConfigFeedBackValue;
import com.iflytek.ccr.polaris.companion.common.ServiceFeedBackValue;

import java.util.concurrent.ArrayBlockingQueue;

public class LocalCacheThread implements Runnable {

    private static final EasyLogger logger = EasyLoggerFactory.getInstance(LocalCacheThread.class);

    private ArrayBlockingQueue<Object> feedBackQueue;

    public LocalCacheThread(ArrayBlockingQueue<Object> feedBackQueue) {
        this.feedBackQueue = feedBackQueue;
    }

    @Override
    public void run() {
        while (true) {
            try {
                Object t = feedBackQueue.take();
                if (t instanceof ConfigFeedBackValue) {
                    boolean flag = FeedbackCacheUtil.getInstance().add((ConfigFeedBackValue) t);
                    if (!flag) {
                        feedBackQueue.add(t);
                    }
                }
                if (t instanceof ServiceFeedBackValue) {
                    boolean flag = ServiceCacheUtil.getInstance().add((ServiceFeedBackValue) t);
                    if (!flag) {
                        feedBackQueue.add(t);
                    }
                }
                //防止zk出问题，持续循环无效处理
                Thread.sleep(5000);

            } catch (Exception e) {
                logger.error(e);
            }
        }
    }
}
