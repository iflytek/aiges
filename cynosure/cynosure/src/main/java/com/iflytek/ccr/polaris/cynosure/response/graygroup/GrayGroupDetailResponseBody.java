package com.iflytek.ccr.polaris.cynosure.response.graygroup;

import java.io.Serializable;
import java.util.Date;

/**
 * Created by DELL-5490 on 2018/7/2.
 */
public class GrayGroupDetailResponseBody implements Serializable {
    private static final long serialVersionUID = 3465537467615599850L;

    //灰度组id
    private String id;

    //版本id
    private String versionId;

    //用户id
    private String userId;

    //灰度组名称
    private String name;

    //推送实例内容
    private String content;

    //灰度组描述
    private String desc;

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

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getContent() {
        return content;
    }

    public void setContent(String content) {
        this.content = content;
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

    @Override
    public String toString() {
        return "GrayGroupDetailResponseBody{" +
                "id='" + id + '\'' +
                ", versionId='" + versionId + '\'' +
                ", userId='" + userId + '\'' +
                ", name='" + name + '\'' +
                ", content='" + content + '\'' +
                ", desc='" + desc + '\'' +
                ", createTime=" + createTime +
                ", updateTime=" + updateTime +
                '}';
    }
}
