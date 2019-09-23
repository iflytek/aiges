package com.iflytek.ccr.polaris.cynosure.domain;

import java.util.Date;

/**
 * 服务发现推送历史模型
 *
 * @author sctang2
 * @create 2017-12-12 18:38
 **/
public class ServiceDiscoveryPushHistory {
    //推送id
    private String id;

    //用户id
    private String userId;

    //项目名称
    private String project;

    //服务组名称
    private String serviceGroup;

    //服务名称
    private String service;

    //服务版本
    private String version;

    //推送时间
    private Date pushTime;

    //集群
    private String clusterText;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getUserId() {
        return userId;
    }

    public void setUserId(String userId) {
        this.userId = userId;
    }

    public String getProject() {
        return project;
    }

    public void setProject(String project) {
        this.project = project;
    }

    public String getServiceGroup() {
        return serviceGroup;
    }

    public void setServiceGroup(String serviceGroup) {
        this.serviceGroup = serviceGroup;
    }

    public String getService() {
        return service;
    }

    public void setService(String service) {
        this.service = service;
    }

    public Date getPushTime() {
        return pushTime;
    }

    public void setPushTime(Date pushTime) {
        this.pushTime = pushTime;
    }

    public String getClusterText() {
        return clusterText;
    }

    public void setClusterText(String clusterText) {
        this.clusterText = clusterText;
    }

    public String getVersion() {
        return version;
    }

    public void setVersion(String version) {
        this.version = version;
    }
}
