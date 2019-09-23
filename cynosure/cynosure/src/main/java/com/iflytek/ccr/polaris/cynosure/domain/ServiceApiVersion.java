package com.iflytek.ccr.polaris.cynosure.domain;

import java.util.Date;
import java.util.List;

/**
 * 服务版本模型
 *
 * @author sctang2
 * @create 2017-11-17 14:50
 **/
public class ServiceApiVersion {
    //版本id
    private String id;

    //版本号
    private String apiVersion;

    //服务id
    private String serviceId;

    //用户id
    private String userId;

    //版本描述
    private String description;

    //创建时间
    private Date createTime;

    //更新时间
    private Date updateTime;

    //配置列表
    private List<ServiceConfig> serviceConfigList;

    //服务
    private Service service;

    //服务组
    private Cluster serviceGroup;

    //项目
    private Project project;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getApiVersion() {
        return apiVersion;
    }

    public void setApiVersion(String apiVersion) {
        this.apiVersion = apiVersion;
    }

    public String getServiceId() {
        return serviceId;
    }

    public void setServiceId(String serviceId) {
        this.serviceId = serviceId;
    }

    public String getUserId() {
        return userId;
    }

    public void setUserId(String userId) {
        this.userId = userId;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public Date getCreateTime() {
        return createTime;
    }

    public void setCreateTime(Date createTime) {
        this.createTime = createTime;
    }

    public Date getUpdateTime() {
        return updateTime;
    }

    public void setUpdateTime(Date updateTime) {
        this.updateTime = updateTime;
    }

    public List<ServiceConfig> getServiceConfigList() {
        return serviceConfigList;
    }

    public void setServiceConfigList(List<ServiceConfig> serviceConfigList) {
        this.serviceConfigList = serviceConfigList;
    }

    public Service getService() {
        return service;
    }

    public void setService(Service service) {
        this.service = service;
    }

    public Cluster getServiceGroup() {
        return serviceGroup;
    }

    public void setServiceGroup(Cluster serviceGroup) {
        this.serviceGroup = serviceGroup;
    }

    public Project getProject() {
        return project;
    }

    public void setProject(Project project) {
        this.project = project;
    }
}
