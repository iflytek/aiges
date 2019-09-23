package com.iflytek.ccr.polaris.cynosure.domain;

import java.util.Date;
import java.util.List;

/**
 * 服务模型
 *
 * @author sctang2
 * @create 2017-11-16 17:45
 **/
public class Service {
    //服务id
    private String id;

    //服务名称
    private String name;

    //服务描述
    private String description;

    //用户id
    private String userId;

    //服务组id
    private String groupId;

    //创建时间
    private Date createTime;

    //更新时间
    private Date updateTime;

    //项目
    private Project project;

    //服务版本列表
    private List<ServiceVersion> serviceVersionList;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
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

    public String getGroupId() {
        return groupId;
    }

    public void setGroupId(String groupId) {
        this.groupId = groupId;
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

    public List<ServiceVersion> getServiceVersionList() {
        return serviceVersionList;
    }

    public void setServiceVersionList(List<ServiceVersion> serviceVersionList) {
        this.serviceVersionList = serviceVersionList;
    }

    public Project getProject() {
        return project;
    }

    public void setProject(Project project) {
        this.project = project;
    }

    @Override
    public String toString() {
        return "Service{" +
                "id='" + id + '\'' +
                ", name='" + name + '\'' +
                ", description='" + description + '\'' +
                ", userId='" + userId + '\'' +
                ", groupId='" + groupId + '\'' +
                ", createTime=" + createTime +
                ", updateTime=" + updateTime +
                ", serviceVersionList=" + serviceVersionList +
                '}';
    }
}
