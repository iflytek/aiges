package com.iflytek.ccr.polaris.cynosure.domain;

import java.util.Date;

/**
 * 项目成员模型
 *
 * @author sctang2
 * @create 2018-01-15 20:03
 **/
public class ProjectMember {
    //唯一标识
    private String id;

    //用户id
    private String userId;

    //项目id
    private String projectId;

    //创建时间
    private Date createTime;

    //是否为创建者
    private Byte creator;

    //用户实体
    private User user;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getUserId() {
        return userId;
    }

    public void setUserId(String userId) {
        this.userId = userId;
    }

    public String getProjectId() {
        return projectId;
    }

    public void setProjectId(String projectId) {
        this.projectId = projectId;
    }

    public Date getCreateTime() {
        return createTime;
    }

    public void setCreateTime(Date createTime) {
        this.createTime = createTime;
    }

    public Byte getCreator() {
        return creator;
    }

    public void setCreator(Byte creator) {
        this.creator = creator;
    }

    public User getUser() {
        return user;
    }

    public void setUser(User user) {
        this.user = user;
    }
}
