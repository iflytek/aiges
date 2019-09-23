package com.iflytek.ccr.polaris.cynosure.domain;

import java.util.Date;

/**
 * 灰度组模型
 *
 * @author sctang2
 * @create 2017-11-15 17:12
 **/
public class GrayGroup {
    //唯一标识
    private String id;

    //灰度组名称
    private String name;

    //灰度组描述
    private String description;

    //推送实例内容
    private String content;

    //版本id
    private String versionId;

    //用户id
    private String userId;

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

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public String getContent() {
        return content;
    }

    public void setContent(String content) {
        this.content = content;
    }

    public String getVersionId() {
        return versionId;
    }

    public void setVersionId(String versionId) {
        this.versionId = versionId;
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

    @Override
    public String toString() {
        return "GrayGroup{" +
                "id='" + id + '\'' +
                ", name='" + name + '\'' +
                ", description='" + description + '\'' +
                ", content='" + content + '\'' +
                ", versionId='" + versionId + '\'' +
                ", userId='" + userId + '\'' +
                ", createTime=" + createTime +
                ", updateTime=" + updateTime +
                '}';
    }
}
