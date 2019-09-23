package com.iflytek.ccr.polaris.companion.common;

/**
 * 配置反馈对象
 */
public class ConfigFeedBackValue {

    private String pushId;
    private String grayGroupId;
    private String project;
    private String group;
    private String service;
    private String version;
    private String config;
    private String addr;
    private String updateStatus;
    private String updateTime;
    private String loadStatus;
    private String loadTime;

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

    public String getService() {
        return service;
    }

    public void setService(String service) {
        this.service = service;
    }

    public String getVersion() {
        return version;
    }

    public void setVersion(String version) {
        this.version = version;
    }

    public String getConfig() {
        return config;
    }

    public void setConfig(String config) {
        this.config = config;
    }

    public String getAddr() {
        return addr;
    }

    public void setAddr(String addr) {
        this.addr = addr;
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

    public String getGrayGroupId() {
        return grayGroupId;
    }

    public void setGrayGroupId(String grayGroupId) {
        this.grayGroupId = grayGroupId;
    }

    @Override
    public String toString() {
        return "ConfigFeedBackValue{" +
                "pushId='" + pushId + '\'' +
                ", grayGroupId='" + grayGroupId + '\'' +
                ", project='" + project + '\'' +
                ", group='" + group + '\'' +
                ", service='" + service + '\'' +
                ", version='" + version + '\'' +
                ", config='" + config + '\'' +
                ", addr='" + addr + '\'' +
                ", updateStatus='" + updateStatus + '\'' +
                ", updateTime='" + updateTime + '\'' +
                ", loadStatus='" + loadStatus + '\'' +
                ", loadTime='" + loadTime + '\'' +
                '}';
    }
}

