package com.iflytek.ccr.finder.utils;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.HashMap;
import java.util.Map;

public class ConfigManager {
    public final static String KEY_WEBSITE_URL = "websiteUrl";
    private final Logger logger = LoggerFactory.getLogger(ConfigManager.class);
    private Map<String, Object> configMap = new HashMap<String, Object>();

    private ConfigManager() {
    }

    public static final ConfigManager getInstance() {
        return ConfigManagerHolder.INSTANCE;
    }

    public String getStringConfigByKey(String key) {
        return (String) configMap.get(key);
    }

    public Object getObjectConfigByKey(String key) {
        return configMap.get(key);
    }

    public int getIntConfigByKey(String key) {
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
