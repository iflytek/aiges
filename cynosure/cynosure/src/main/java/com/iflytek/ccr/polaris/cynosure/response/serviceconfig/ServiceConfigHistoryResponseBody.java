package com.iflytek.ccr.polaris.cynosure.response.serviceconfig;

import java.io.Serializable;
import java.util.Date;

/**
 * 服务配置历史-响应
 *
 * @author sctang2
 * @create 2017-11-22 11:58
 **/
public class ServiceConfigHistoryResponseBody implements Serializable {
    private static final long serialVersionUID = -179589381259520998L;

    //配置历史id
    private String id;

    //配置id
    private String configId;

    //配置名称
    private String name;

    //配置内容
    private String content;

    //配置描述
    private String desc;

    //推送版本号
    private String pushVersion;

    //创建时间
    private Date createTime;

    public ServiceConfigHistoryResponseBody() {
    }

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getConfigId() {
        return configId;
    }

    public void setConfigId(String configId) {
        this.configId = configId;
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

    public String getPushVersion() {
        return pushVersion;
    }

    public void setPushVersion(String pushVersion) {
        this.pushVersion = pushVersion;
    }

    public Date getCreateTime() {
        return createTime;
    }

    public void setCreateTime(Date createTime) {
        this.createTime = createTime;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    @Override
    public String toString() {
        return "ServiceConfigHistoryResponseBody{" +
                "id='" + id + '\'' +
                ", configId='" + configId + '\'' +
                ", name='" + name + '\'' +
                ", content='" + content + '\'' +
                ", desc='" + desc + '\'' +
                ", pushVersion='" + pushVersion + '\'' +
                ", createTime=" + createTime +
                '}';
    }
}
