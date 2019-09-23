package com.iflytek.ccr.polaris.companion.utils;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;

import java.util.HashMap;
import java.util.Map;

public class ConfigManager {
    public final static String KEY_HOST = "host";
    public final static String KEY_PORT = "port";
    public final static String KEY_ZKSTR = "zkStr";
    public final static String KEY_WEBSITE_URL = "websiteUrl";
    public final static String KEY_ZK_NODE_PATH = "zkNodePath";
    private static final EasyLogger logger = EasyLoggerFactory.getInstance(ConfigManager.class);
    private static Map<String, Object> configMap = new HashMap<String, Object>();

    private ConfigManager() {
    }

    public static final ConfigManager getInstance() {
        return ConfigManager.ConfigManagerHolder.INSTANCE;
    }

    public static String getStringConfigByKey(String key) {
        return (String) configMap.get(key);
    }

    public static int getIntConfigByKey(String key) {
        Object o = configMap.get(key);
        if (null == o) {
            return 0;
        }
        return (int) o;
    }

    public void put(String key, Object value) {
        configMap.put(key, value);
    }

    private static class ConfigManagerHolder {
        private static final ConfigManager INSTANCE = new ConfigManager();
    }
}
