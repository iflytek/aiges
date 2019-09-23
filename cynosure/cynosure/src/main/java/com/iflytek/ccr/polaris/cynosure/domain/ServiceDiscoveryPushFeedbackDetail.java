package com.iflytek.ccr.polaris.cynosure.domain;

import java.util.Date;

/**
 * Created by DELL-5490 on 2018/8/10.
 */
public class ServiceDiscoveryPushFeedbackDetail {
    //推送反馈id
    private String id;

    //推送id
    private String pushId;

    //项目名称
    private String project;

    //服务组名称
    private String serviceGroup;

    //消费端服务名称
    private String consumerService;

    //消费端版本
    private String consumerVersion;

    //提供端服务名称
    private String providerService;

    //提供端版本
    private String providerVersion;

    //地址
    private String addr;

    //更新状态
    private Byte updateStatus;

    //加载状态
    private Byte loadStatus;

    //更新时间
    private Date updateTime;

    //加载时间
    private Date loadTime;

    //创建时间
    private Date createTime;

    //服务改变的类型
    private String type;

    private String typeName;

    private String apiVersion;

    public String getTypeName() {
        return typeName;
    }

    public void setTypeName(String typeName) {
        this.typeName = typeName;
    }

    public String getType() {
        return type;
    }

    public void setType(String type) {
        this.type = type;
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
        this.id = id == null ? null : id.trim();
    }

    public String getPushId() {
        return pushId;
    }

    public void setPushId(String pushId) {
        this.pushId = pushId == null ? null : pushId.trim();
    }

    public String getProject() {
        return project;
    }

    public void setProject(String project) {
        this.project = project == null ? null : project.trim();
    }

    public String getServiceGroup() {
        return serviceGroup;
    }

    public void setServiceGroup(String serviceGroup) {
        this.serviceGroup = serviceGroup == null ? null : serviceGroup.trim();
    }

    public String getConsumerService() {
        return consumerService;
    }

    public void setConsumerService(String consumerService) {
        this.consumerService = consumerService == null ? null : consumerService.trim();
    }

    public String getConsumerVersion() {
        return consumerVersion;
    }

    public void setConsumerVersion(String consumerVersion) {
        this.consumerVersion = consumerVersion == null ? null : consumerVersion.trim();
    }

    public String getProviderService() {
        return providerService;
    }

    public void setProviderService(String providerService) {
        this.providerService = providerService == null ? null : providerService.trim();
    }

    public String getProviderVersion() {
        return providerVersion;
    }

    public void setProviderVersion(String providerVersion) {
        this.providerVersion = providerVersion == null ? null : providerVersion.trim();
    }

    public String getAddr() {
        return addr;
    }

    public void setAddr(String addr) {
        this.addr = addr == null ? null : addr.trim();
    }

    public Byte getUpdateStatus() {
        return updateStatus;
    }

    public void setUpdateStatus(Byte updateStatus) {
        this.updateStatus = updateStatus;
    }

    public Byte getLoadStatus() {
        return loadStatus;
    }

    public void setLoadStatus(Byte loadStatus) {
        this.loadStatus = loadStatus;
    }

    public Date getUpdateTime() {
        return updateTime;
    }

    public void setUpdateTime(Date updateTime) {
        this.updateTime = updateTime;
    }

    public Date getLoadTime() {
        return loadTime;
    }

    public void setLoadTime(Date loadTime) {
        this.loadTime = loadTime;
    }

    public Date getCreateTime() {
        return createTime;
    }

    public void setCreateTime(Date createTime) {
        this.createTime = createTime;
    }
}
