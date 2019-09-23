package com.iflytek.ccr.finder.value;

import java.io.Serializable;

/**
 * 服务配置对象
 */
public class ServiceConfig implements Serializable {

//    private List<RouteInfo> configList;

    /**
     * json配置信息，供服务订阅者使用
     */
    private String jsonConfig;

    public String getJsonConfig() {
        return jsonConfig;
    }

    public void setJsonConfig(String jsonConfig) {
        this.jsonConfig = jsonConfig;
    }

//    public List<RouteInfo> getConfigList() {
//        return configList;
//    }
//
//    public void setConfigList(List<RouteInfo> configList) {
//        this.configList = configList;
//    }

    @Override
    public String toString() {
        return "ServiceConfig{" +
                "jsonConfig='" + jsonConfig + '\'' +
                '}';
    }
}
