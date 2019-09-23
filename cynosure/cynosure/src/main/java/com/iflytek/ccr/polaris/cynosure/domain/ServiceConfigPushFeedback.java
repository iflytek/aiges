package com.iflytek.ccr.polaris.cynosure.domain;

import java.util.Date;

/**
 * 服务配置推送反馈模型
 *
 * @author sctang2
 * @create 2017-11-25 11:14
 **/
public class ServiceConfigPushFeedback {
    //推送反馈id
    private String id;

    //推送id
    private String pushId;

    //项目名称
    private String project;

    //服务组名称
    private String serviceGroup;

    //服务名称
    private String service;

    //服务版本号
    private String version;

    //配置名称
    private String config;

    //地址
    private String addr;

    //更新状态
    private Byte updateStatus;

    //加载状态
    private Byte loadStatus;

    //更新时间
    private Date updateTime;

    //加载时间
    private Date loadTime;

    //创建时间
    private Date createTime;

    //灰度组标识
    private String grayGroupId;

    //灰度组名字
    private String grayGroupName;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id == null ? null : id.trim();
    }

    public String getPushId() {
        return pushId;
    }

    public void setPushId(String pushId) {
        this.pushId = pushId == null ? null : pushId.trim();
    }

    public String getProject() {
        return project;
    }

    public void setProject(String project) {
        this.project = project == null ? null : project.trim();
    }

    public String getServiceGroup() {
        return serviceGroup;
    }

    public void setServiceGroup(String serviceGroup) {
        this.serviceGroup = serviceGroup == null ? null : serviceGroup.trim();
    }

    public String getService() {
        return service;
    }

    public void setService(String service) {
        this.service = service == null ? null : service.trim();
    }

    public String getVersion() {
        return version;
    }

    public void setVersion(String version) {
        this.version = version == null ? null : version.trim();
    }

    public String getConfig() {
        return config;
    }

    public void setConfig(String config) {
        this.config = config == null ? null : config.trim();
    }

    public String getAddr() {
        return addr;
    }

    public void setAddr(String addr) {
        this.addr = addr == null ? null : addr.trim();
    }

    public Byte getUpdateStatus() {
        return updateStatus;
    }

    public void setUpdateStatus(Byte updateStatus) {
        this.updateStatus = updateStatus;
    }

    public Byte getLoadStatus() {
        return loadStatus;
    }

    public void setLoadStatus(Byte loadStatus) {
        this.loadStatus = loadStatus;
    }

    public Date getUpdateTime() {
        return updateTime;
    }

    public void setUpdateTime(Date updateTime) {
        this.updateTime = updateTime;
    }

    public Date getLoadTime() {
        return loadTime;
    }

    public void setLoadTime(Date loadTime) {
        this.loadTime = loadTime;
    }

    public Date getCreateTime() {
        return createTime;
    }

    public void setCreateTime(Date createTime) {
        this.createTime = createTime;
    }

    public String getGrayGroupId() {
        return grayGroupId;
    }

    public void setGrayGroupId(String grayGroupId) {
        this.grayGroupId = grayGroupId;
    }

    public String getGrayGroupName() {
        return grayGroupName;
    }

    public void setGrayGroupName(String grayGroupName) {
        this.grayGroupName = grayGroupName;
    }

    @Override
    public String toString() {
        return "ServiceConfigPushFeedback{" +
                "id='" + id + '\'' +
                ", pushId='" + pushId + '\'' +
                ", project='" + project + '\'' +
                ", serviceGroup='" + serviceGroup + '\'' +
                ", service='" + service + '\'' +
                ", version='" + version + '\'' +
                ", config='" + config + '\'' +
                ", addr='" + addr + '\'' +
                ", updateStatus=" + updateStatus +
                ", loadStatus=" + loadStatus +
                ", updateTime=" + updateTime +
                ", loadTime=" + loadTime +
                ", createTime=" + createTime +
                ", grayGroupId='" + grayGroupId + '\'' +
                ", grayGroupName='" + grayGroupName + '\'' +
                '}';
    }
}
