package com.iflytek.ccr.finder.value;

import java.util.Arrays;
import java.util.Map;

/**
 * 配置信息对象
 */
public class Config {

    /**
     * 配置文件名称
     */
    private String name;

    /**
     * 配置文件字节内容
     */
    private byte[] file;

    /**
     * toml类型的配置文件解析后的集合
     */
    private Map<String, Object> configMap;

    public Map<String, Object> getConfigMap() {
        return configMap;
    }

    public void setConfigMap(Map<String, Object> configMap) {
        this.configMap = configMap;
    }


    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public byte[] getFile() {
        return file;
    }

    public void setFile(byte[] file) {
        this.file = file;
    }

    @Override
    public String toString() {
        return "Config{" +
                "name='" + name + '\'' +
                ", file=" + Arrays.toString(file) +
                ", configMap=" + configMap +
                '}';
    }
}
