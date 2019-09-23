package com.iflytek.ccr.polaris.companion.common;

/**
 * 反馈对象
 */
public class ServiceFeedBackValue {
    private String pushId ;
    private String project;
    private String group ;
    private String consumer;
    private String provider;
    private String addr;
    private String type;
    private String consumerVersion ;
    private String providerVersion;
    private String apiVersion;
    private String updateStatus;
    private String updateTime;
    private String loadStatus;
    private String loadTime ;

    public String getApiVersion() {
        return apiVersion;
    }

    public void setApiVersion(String apiVersion) {
        this.apiVersion = apiVersion;
    }

    public String getType() {
        return type;
    }

    public void setType(String type) {
        this.type = type;
    }

    public String getProviderVersion() {
        return providerVersion;
    }

    public void setProviderVersion(String providerVersion) {
        this.providerVersion = providerVersion;
    }

    public String getPushId() {
        return pushId;
    }

    public void setPushId(String pushId) {
        this.pushId = pushId;
    }

    public String getProject() {
        return project;
    }

    public void setProject(String project) {
        this.project = project;
    }

    public String getGroup() {
        return group;
    }

    public void setGroup(String group) {
        this.group = group;
    }

    public String getConsumer() {
        return consumer;
    }

    public void setConsumer(String consumer) {
        this.consumer = consumer;
    }

    public String getProvider() {
        return provider;
    }

    public void setProvider(String provider) {
        this.provider = provider;
    }

    public String getAddr() {
        return addr;
    }

    public void setAddr(String addr) {
        this.addr = addr;
    }

    public String getConsumerVersion() {
        return consumerVersion;
    }

    public void setConsumerVersion(String consumerVersion) {
        this.consumerVersion = consumerVersion;
    }


    public String getUpdateStatus() {
        return updateStatus;
    }

    public void setUpdateStatus(String updateStatus) {
        this.updateStatus = updateStatus;
    }

    public String getUpdateTime() {
        return updateTime;
    }

    public void setUpdateTime(String updateTime) {
        this.updateTime = updateTime;
    }

    public String getLoadStatus() {
        return loadStatus;
    }

    public void setLoadStatus(String loadStatus) {
        this.loadStatus = loadStatus;
    }

    public String getLoadTime() {
        return loadTime;
    }

    public void setLoadTime(String loadTime) {
        this.loadTime = loadTime;
    }

    @Override
    public String toString() {
        return "ServiceFeedBackValue{" +
                "pushId='" + pushId + '\'' +
                ", project='" + project + '\'' +
                ", group='" + group + '\'' +
                ", consumer='" + consumer + '\'' +
                ", provider='" + provider + '\'' +
                ", addr='" + addr + '\'' +
                ", type='" + type + '\'' +
                ", consumerVersion='" + consumerVersion + '\'' +
                ", providerVersion='" + providerVersion + '\'' +
                ", apiVersion='" + apiVersion + '\'' +
                ", updateStatus='" + updateStatus + '\'' +
                ", updateTime='" + updateTime + '\'' +
                ", loadStatus='" + loadStatus + '\'' +
                ", loadTime='" + loadTime + '\'' +
                '}';
    }
}

