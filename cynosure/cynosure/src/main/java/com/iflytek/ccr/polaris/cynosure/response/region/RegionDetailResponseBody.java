package com.iflytek.ccr.polaris.cynosure.response.region;

import java.io.Serializable;
import java.util.Date;

/**
 * 区域明细-响应
 *
 * @author sctang2
 * @create 2017-11-14 20:48
 **/
public class RegionDetailResponseBody implements Serializable {
    private static final long serialVersionUID = -8813622248841127526L;

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

    @Override
    public String toString() {
        return "RegionDetailResponseBody{" +
                "id='" + id + '\'' +
                ", name='" + name + '\'' +
                ", pushUrl='" + pushUrl + '\'' +
                ", createTime=" + createTime +
                ", updateTime=" + updateTime +
                '}';
    }
}
