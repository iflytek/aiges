package com.iflytek.ccr.polaris.cynosure.domain;

import java.util.Date;

/**
 * 最近搜索模型
 *
 * @author sctang2
 * @create 2018-01-26 10:12
 **/
public class LastestSearch {
    //唯一标识
    private String id;

    //用户id
    private String userId;

    //请求地址
    private String url;

    //前置条件
    private String preCondition;

    //创建时间
    private Date createTime;

    //更新时间
    private Date updateTime;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id == null ? null : id.trim();
    }

    public String getUserId() {
        return userId;
    }

    public void setUserId(String userId) {
        this.userId = userId == null ? null : userId.trim();
    }

    public String getUrl() {
        return url;
    }

    public void setUrl(String url) {
        this.url = url == null ? null : url.trim();
    }

    public String getPreCondition() {
        return preCondition;
    }

    public void setPreCondition(String preCondition) {
        this.preCondition = preCondition;
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
}
