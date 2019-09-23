package com.iflytek.ccr.polaris.cynosure.domain;

import java.util.Date;
import java.util.List;

/**
 * 集群模型
 *
 * @author sctang2
 * @create 2017-11-15 17:12
 **/
public class Cluster {
    //唯一标识
    private String id;

    //集群名称
    private String name;

    //集群描述
    private String description;

    //用户id
    private String userId;

    //创建时间
    private Date createTime;

    //更新时间
    private Date updateTime;

    //项目id
    private String projectId;

    //服务列表
    private List<Service> serviceList;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id == null ? null : id.trim();
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name == null ? null : name.trim();
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public String getUserId() {
        return userId;
    }

    public void setUserId(String userId) {
        this.userId = userId;
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

    public String getProjectId() {
        return projectId;
    }

    public void setProjectId(String projectId) {
        this.projectId = projectId;
    }

    public List<Service> getServiceList() {
        return serviceList;
    }

    public void setServiceList(List<Service> serviceList) {
        this.serviceList = serviceList;
    }

    @Override
    public String toString() {
        return "Cluster{" +
                "id='" + id + '\'' +
                ", name='" + name + '\'' +
                ", description='" + description + '\'' +
                ", userId='" + userId + '\'' +
                ", createTime=" + createTime +
                ", updateTime=" + updateTime +
                ", projectId='" + projectId + '\'' +
                ", serviceList=" + serviceList +
                '}';
    }
}
