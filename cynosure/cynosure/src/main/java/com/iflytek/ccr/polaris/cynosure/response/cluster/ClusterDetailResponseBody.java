package com.iflytek.ccr.polaris.cynosure.response.cluster;

import java.io.Serializable;
import java.util.Date;

/**
 * 集群明细-响应
 *
 * @author sctang2
 * @create 2017-11-16 13:49
 **/
public class ClusterDetailResponseBody implements Serializable {
    private static final long serialVersionUID = 168875798149643457L;

    //唯一标识
    private String id;

    //集群名称
    private String name;

    //集群描述
    private String desc;

    //创建时间
    private Date createTime;

    //更新时间
    private Date updateTime;

    //项目id
    private String projectId;

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

    public String getDesc() {
        return desc;
    }

    public void setDesc(String desc) {
        this.desc = desc;
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

    @Override
    public String toString() {
        return "ClusterDetailResponseBody{" +
                "id='" + id + '\'' +
                ", name='" + name + '\'' +
                ", desc='" + desc + '\'' +
                ", createTime=" + createTime +
                ", updateTime=" + updateTime +
                ", projectId='" + projectId + '\'' +
                '}';
    }
}
