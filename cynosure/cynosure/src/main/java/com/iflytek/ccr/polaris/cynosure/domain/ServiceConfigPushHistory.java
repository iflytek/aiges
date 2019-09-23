package com.iflytek.ccr.polaris.cynosure.domain;

import java.util.Date;

/**
 * 服务配置推送历史-模型
 *
 * @author sctang2
 * @create 2017-11-24 9:12
 **/
public class ServiceConfigPushHistory {
    //推送id
    private String id;

    //用户id
    private String userId;

    //灰度组id
    private String grayId;

    //项目名称
    private String project;

    //服务组名称
    private String serviceGroup;

    //服务名称
    private String service;

    //服务版本号
    private String version;

    //集群
    private String clusterText;

    //服务配置
    private String serviceConfigText;

    //推送时间
    private Date pushTime;

    public ServiceConfigPushHistory() {
    }

    public ServiceConfigPushHistory(String id, String userId, String grayId, String project, String serviceGroup, String service, String version, String clusterText, String serviceConfigText, Date pushTime) {
        this.id = id;
        this.userId = userId;
        this.grayId = grayId;
        this.project = project;
        this.serviceGroup = serviceGroup;
        this.service = service;
        this.version = version;
        this.clusterText = clusterText;
        this.serviceConfigText = serviceConfigText;
        this.pushTime = pushTime;
    }

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

    public String getVersion() {
        return version;
    }

    public void setVersion(String version) {
        this.version = version;
    }

    public String getClusterText() {
        return clusterText;
    }

    public void setClusterText(String clusterText) {
        this.clusterText = clusterText;
    }

    public String getServiceConfigText() {
        return serviceConfigText;
    }

    public void setServiceConfigText(String serviceConfigText) {
        this.serviceConfigText = serviceConfigText;
    }

    public Date getPushTime() {
        return pushTime;
    }

    public void setPushTime(Date pushTime) {
        this.pushTime = pushTime;
    }

    public String getGrayId() {
        return grayId;
    }

    public void setGrayId(String grayId) {
        this.grayId = grayId;
    }

    @Override
    public String toString() {
        return "ServiceConfigPushHistory{" +
                "id='" + id + '\'' +
                ", userId='" + userId + '\'' +
                ", grayId='" + grayId + '\'' +
                ", project='" + project + '\'' +
                ", serviceGroup='" + serviceGroup + '\'' +
                ", service='" + service + '\'' +
                ", version='" + version + '\'' +
                ", clusterText='" + clusterText + '\'' +
                ", serviceConfigText='" + serviceConfigText + '\'' +
                ", pushTime=" + pushTime +
                '}';
    }
}
