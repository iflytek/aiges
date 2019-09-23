package com.iflytek.ccr.finder.value;

/**
 * 服务发现订阅请求对象
 */
public class SubscribeRequestValue {
    /**
     * 服务名称
     */
    private String serviceName;

    /**
     * api版本号
     */
    private String apiVersion;

    public String getServiceName() {
        return serviceName;
    }

    public void setServiceName(String serviceName) {
        this.serviceName = serviceName;
    }

    public String getApiVersion() {
        return apiVersion;
    }

    public void setApiVersion(String apiVersion) {
        this.apiVersion = apiVersion;
    }

    public String getCacheKey() {
        return this.getServiceName() + "_" + this.getApiVersion();
    }

    @Override
    public String toString() {
        return "SubscribeRequestValue{" +
                "serviceName='" + serviceName + '\'' +
                ", apiVersion='" + apiVersion + '\'' +
                '}';
    }
}
