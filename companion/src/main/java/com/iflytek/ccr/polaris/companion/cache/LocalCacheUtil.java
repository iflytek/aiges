package com.iflytek.ccr.polaris.companion.cache;

import java.util.concurrent.ArrayBlockingQueue;
import java.util.concurrent.Executors;

public class LocalCacheUtil {

    public static ArrayBlockingQueue<Object> feedBackQueue = new ArrayBlockingQueue<Object>(10000);
    private static boolean isInit = false;

    public static void init() {
        if (!isInit) {
            Executors.newSingleThreadExecutor().submit(new LocalCacheThread(feedBackQueue));
        } else {
            isInit = true;
        }

    }


}
