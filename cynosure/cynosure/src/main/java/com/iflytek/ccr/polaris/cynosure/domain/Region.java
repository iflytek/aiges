package com.iflytek.ccr.polaris.cynosure.domain;

import java.util.Date;

/**
 * 区域模型
 *
 * @author sctang2
 * @create 2017-11-14 19:55
 **/
public class Region {
    //唯一标识
    private String id;

    //区域名称
    private String name;

    //推送地址
    private String pushUrl;

    //创建时间
    private Date createTime;

    //更新时间
    private Date updateTime;

    //用户id
    private String userId;

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

    public String getPushUrl() {
        return pushUrl;
    }

    public void setPushUrl(String pushUrl) {
        this.pushUrl = pushUrl;
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

    public String getUserId() {
        return userId;
    }

    public void setUserId(String userId) {
        this.userId = userId;
    }

    @Override
    public String toString() {
        return "Region{" +
                "id='" + id + '\'' +
                ", name='" + name + '\'' +
                ", pushUrl='" + pushUrl + '\'' +
                ", createTime=" + createTime +
                ", updateTime=" + updateTime +
                ", userId='" + userId + '\'' +
                '}';
    }
}
