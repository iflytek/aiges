package com.iflytek.ccr.polaris.cynosure.response.track;

import java.io.Serializable;

/**
 * 轨迹发现明细-响应
 *
 * @author sctang2
 * @create 2017-12-14 9:47
 **/
public class TrackDiscoveryDetailResponseBody implements Serializable {
    private static final long serialVersionUID = 8376224965503859607L;

    //反馈id
    private String id;

    //推送id
    private String pushId;

    //项目名称
    private String project;

    //服务组名称
    private String cluster;

    //消费方服务名称
    private String consumerService;

    //消费方版本号
    private String consumerVersion;

    //提供发服务名称
    private String providerService;

    //提供发版本号
    private String providerVersion;

    //地址
    private String addr;

    //更新状态
    private int updateStatus;

    //更新时间
    private long updateTime;

    //加载状态
    private int loadStatus;

    //加载时间
    private long loadTime;

    private String type;

    private String typeName;

    private String apiVersion;

    public String getType() {
        return type;
    }

    public void setType(String type) {
        this.type = type;
    }

    public String getTypeName() {
        return typeName;
    }

    public void setTypeName(String typeName) {
        this.typeName = typeName;
    }

    public String getApiVersion() {
        return apiVersion;
    }

    public void setApiVersion(String apiVersion) {
        this.apiVersion = apiVersion;
    }

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
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

    public String getCluster() {
        return cluster;
    }

    public void setCluster(String cluster) {
        this.cluster = cluster;
    }

    public String getConsumerService() {
        return consumerService;
    }

    public void setConsumerService(String consumerService) {
        this.consumerService = consumerService;
    }

    public String getConsumerVersion() {
        return consumerVersion;
    }

    public void setConsumerVersion(String consumerVersion) {
        this.consumerVersion = consumerVersion;
    }

    public String getProviderService() {
        return providerService;
    }

    public void setProviderService(String providerService) {
        this.providerService = providerService;
    }

    public String getProviderVersion() {
        return providerVersion;
    }

    public void setProviderVersion(String providerVersion) {
        this.providerVersion = providerVersion;
    }

    public String getAddr() {
        return addr;
    }

    public void setAddr(String addr) {
        this.addr = addr;
    }

    public int getUpdateStatus() {
        return updateStatus;
    }

    public void setUpdateStatus(int updateStatus) {
        this.updateStatus = updateStatus;
    }

    public long getUpdateTime() {
        return updateTime;
    }

    public void setUpdateTime(long updateTime) {
        this.updateTime = updateTime;
    }

    public int getLoadStatus() {
        return loadStatus;
    }

    public void setLoadStatus(int loadStatus) {
        this.loadStatus = loadStatus;
    }

    public long getLoadTime() {
        return loadTime;
    }

    public void setLoadTime(long loadTime) {
        this.loadTime = loadTime;
    }
}
