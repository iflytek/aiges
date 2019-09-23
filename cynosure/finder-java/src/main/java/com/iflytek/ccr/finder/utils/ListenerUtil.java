package com.iflytek.ccr.finder.utils;

import org.apache.curator.framework.recipes.cache.NodeCache;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.util.List;

public class ListenerUtil {

    private final static Logger logger = LoggerFactory.getLogger(ListenerUtil.class);

    /**
     * 关闭nodeCache监听
     *
     * @param nodeCachesList
     */
    public static void closeNodeCache(List<NodeCache> nodeCachesList) {
        if (nodeCachesList != null && !nodeCachesList.isEmpty()) {
            for (NodeCache nodeCache : nodeCachesList) {
                try {
                    nodeCache.close();
                } catch (IOException e) {
                    logger.error(e.getMessage());
                }
            }
        }
    }
}
