package com.iflytek.ccr.polaris.cynosure.domain;

import java.util.Date;
import java.util.List;

/**
 * 项目模型
 *
 * @author sctang2
 * @create 2017-11-20 9:03
 **/
public class Project {
    //项目id
    private String id;

    //项目名称
    private String name;

    //项目描述
    private String description;

    //用户id
    private String userId;

    //创建时间
    private Date createTime;

    //更新时间
    private Date updateTime;

    //集群列表
    private List<Cluster> clusterList;

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
        this.description = description == null ? null : description.trim();
    }

    public String getUserId() {
        return userId;
    }

    public void setUserId(String userId) {
        this.userId = userId == null ? null : userId.trim();
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

    public List<Cluster> getClusterList() {
        return clusterList;
    }

    public void setClusterList(List<Cluster> clusterList) {
        this.clusterList = clusterList;
    }
}
