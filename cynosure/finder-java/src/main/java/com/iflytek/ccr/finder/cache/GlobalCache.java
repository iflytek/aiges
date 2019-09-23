package com.iflytek.ccr.finder.cache;

import com.iflytek.ccr.finder.constants.TaskType;
import com.iflytek.ccr.finder.handler.ConfigChangedHandler;
import com.iflytek.ccr.finder.value.GrayConfigValue;
import org.apache.curator.framework.recipes.cache.NodeCache;

import java.util.*;
import java.util.concurrent.ArrayBlockingQueue;

public class GlobalCache {

    /**
     *
     */
    public ArrayBlockingQueue<TaskType> taskQueue = new ArrayBlockingQueue(5);

    /**
     * NodeCache缓存
     */
    public List<NodeCache> nodeCachesList = new ArrayList<>();

    /**
     * Lisener缓存
     */
    public Map<String, Object> configListenerMap = new HashMap<>();

    /**
     * 存放要监听的临时节点路径，这些节点如果不存在了，要自动恢复
     */
    public Vector<String> monitorPathList = new Vector<>();

    /**
     * 存放初始化的时候，已经返回的配置数据
     */
    public Map<String, Boolean> initMap = new HashMap<>();

    /**
     * 配置的路径
     */
    private String basePath;
    /**
     * 配置消费者路径
     */
    private String configConsumerPath;
    /**
     * 灰度配置对象
     */
    private GrayConfigValue grayConfigValue;

    /**
     * 配置文件列表
     */
    private List<String> configNameList;

    public List<String> getConfigNameList() {
        return configNameList;
    }

    public void setConfigNameList(List<String> configNameList) {
        this.configNameList = configNameList;
    }

    /**
     * 配置变更通知接口对象
     */
    private ConfigChangedHandler configChangedHandler;

    public ConfigChangedHandler getConfigChangedHandler() {
        return configChangedHandler;
    }

    public void setConfigChangedHandler(ConfigChangedHandler configChangedHandler) {
        this.configChangedHandler = configChangedHandler;
    }

    public GrayConfigValue getGrayConfigValue() {
        return grayConfigValue;
    }

    public void setGrayConfigValue(GrayConfigValue grayConfigValue) {
        this.grayConfigValue = grayConfigValue;
    }


    public String getConfigConsumerPath() {
        return configConsumerPath;
    }

    public void setConfigConsumerPath(String configConsumerPath) {
        this.configConsumerPath = configConsumerPath;
    }

    public String getBasePath() {
        return basePath;
    }

    public void setBasePath(String basePath) {
        this.basePath = basePath;
    }
}
