package com.iflytek.ccr.polaris.cynosure.response.service;

import java.io.Serializable;
import java.util.Date;

/**
 * 服务明细-响应
 *
 * @author sctang2
 * @create 2017-11-16 19:23
 **/
public class ServiceDetailResponseBody implements Serializable {
    private static final long serialVersionUID = 5109620798832409922L;

    //服务id
    private String id;

    //服务名称
    private String name;

    //服务描述
    private String desc;

    //集群id
    private String clusterId;

    //创建时间
    private Date createTime;

    //更新时间
    private Date updateTime;

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

    public String getClusterId() {
        return clusterId;
    }

    public void setClusterId(String clusterId) {
        this.clusterId = clusterId;
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

    @Override
    public String toString() {
        return "ServiceDetailResponseBody{" +
                "id='" + id + '\'' +
                ", name='" + name + '\'' +
                ", desc='" + desc + '\'' +
                ", clusterId='" + clusterId + '\'' +
                ", createTime=" + createTime +
                ", updateTime=" + updateTime +
                '}';
    }
}
