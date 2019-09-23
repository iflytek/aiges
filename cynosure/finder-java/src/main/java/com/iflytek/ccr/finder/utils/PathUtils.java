package com.iflytek.ccr.finder.utils;

import com.iflytek.ccr.finder.FinderManager;
import com.iflytek.ccr.finder.constants.Constants;
import com.iflytek.ccr.finder.value.SubscribeRequestValue;
import org.codehaus.jackson.map.ObjectMapper;
import org.codehaus.jackson.node.ObjectNode;

import java.io.File;

public class PathUtils {

    /**
     * 获取服务提供者路径
     *
     * @param name
     * @return
     */
    public static String getProviderPath(String name, String apiVersion) {
        StringBuffer pathBuffer = new StringBuffer();
        pathBuffer.append(ConfigManager.getInstance().getStringConfigByKey(Constants.SERVICE_PATH))
                .append("/").append(name).append("/").append(apiVersion).append("/provider");
        return pathBuffer.toString();
    }

    /**
     * 获取conf路径
     *
     * @param name
     * @return
     */
    public static String getConfPath(String name, String apiVersion) {
        StringBuffer pathBuffer = new StringBuffer();
        pathBuffer.append(ConfigManager.getInstance().getStringConfigByKey(Constants.SERVICE_PATH))
                .append("/").append(name).append("/").append(apiVersion).append("/conf");
        return pathBuffer.toString();
    }

    /**
     * 获取conf路径
     *
     * @param name
     * @return
     */
    public static String getRoutePath(String name, String apiVersion) {
        StringBuffer pathBuffer = new StringBuffer();
        pathBuffer.append(ConfigManager.getInstance().getStringConfigByKey(Constants.SERVICE_PATH))
                .append("/").append(name).append("/").append(apiVersion).append("/route");
        return pathBuffer.toString();
    }

    /**
     * 获取服务订阅者路径
     * @param requestValue
     * @param addr
     * @return
     */
    public static String getServiceConsumerPath(SubscribeRequestValue requestValue, String addr) {
        StringBuffer pathBuffer = new StringBuffer();
        pathBuffer.append(ConfigManager.getInstance().getStringConfigByKey(Constants.SERVICE_PATH))
                .append("/").append(requestValue.getServiceName()).append("/").append(requestValue.getApiVersion()).append("/consumer").append("/").append(addr);
        return pathBuffer.toString();
    }

    /**
     * 获取服务配置文件路径
     *
     * @param name
     * @return
     */
    public static String getConfigPath(ConfigManager configManager, String name) {
        StringBuffer pathBuffer = new StringBuffer();
        pathBuffer.append(configManager.getStringConfigByKey(Constants.CONFIG_PATH)).append("/").append(name);
        return pathBuffer.toString();
    }


    public static String getCacheFilePath(FinderManager finderManager, String subPath) {
        return finderManager.getBootConfig().getCachePath() + File.separator + subPath + File.separator;
    }

    public static ObjectNode getServiceDefaultNodes() {
        ObjectMapper mapper = new ObjectMapper();
        ObjectNode nodeData = mapper.createObjectNode();
        nodeData.put(Constants.WEIGHT, Constants.DEFAULT_WEIGHT);
        nodeData.put(Constants.IS_VALID, Constants.DEFAULT_VALID);
        return nodeData;
    }
}
