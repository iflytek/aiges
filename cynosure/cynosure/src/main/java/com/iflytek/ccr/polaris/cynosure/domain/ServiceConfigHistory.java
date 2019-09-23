package com.iflytek.ccr.polaris.cynosure.domain;

import java.util.Date;

/**
 * 服务配置历史-模型
 *
 * @author sctang2
 * @create 2017-11-22 10:07
 **/
public class ServiceConfigHistory {
    //配置id
    private String id;

    //用户id
    private String userId;

    //配置id
    private String configId;

    //配置描述
    private String description;

    //推送版本号
    private String pushVersion;

    //创建时间
    private Date createTime;

    //配置内容
    private byte[] content;

    //配置内容md5
    private String md5;

    //服务配置
    private ServiceConfig serviceConfig;

    public ServiceConfigHistory() {
    }

    public ServiceConfigHistory(String id, String userId, String configId, String description, String pushVersion, Date createTime, byte[] content, String md5) {
        this.id = id;
        this.userId = userId;
        this.configId = configId;
        this.description = description;
        this.pushVersion = pushVersion;
        this.createTime = createTime;
        this.content = content;
        this.md5 = md5;
    }

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

    public String getConfigId() {
        return configId;
    }

    public void setConfigId(String configId) {
        this.configId = configId;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
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

    public byte[] getContent() {
        return content;
    }

    public void setContent(byte[] content) {
        this.content = content;
    }

    public String getMd5() {
        return md5;
    }

    public void setMd5(String md5) {
        this.md5 = md5;
    }

    public ServiceConfig getServiceConfig() {
        return serviceConfig;
    }

    public void setServiceConfig(ServiceConfig serviceConfig) {
        this.serviceConfig = serviceConfig;
    }
}
