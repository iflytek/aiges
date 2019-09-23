package com.iflytek.ccr.polaris.cynosure.request.serviceconfig;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import org.hibernate.validator.constraints.Length;
import org.hibernate.validator.constraints.NotBlank;

import java.io.Serializable;
import java.util.Date;

/**
 * 服务配置更新反馈-请求
 *
 * @author sctang2
 * @create 2017-11-25 10:20
 **/
public class ServiceConfigFeedBackRequestBody implements Serializable {
    private static final long serialVersionUID = 4557792459406383764L;

    //推送id
    @NotBlank(message = SystemErrCode.ERRMSG_TRACK_ID_NOT_NULL)
    @Length(max = 50, message = SystemErrCode.ERRMSG_TRACK_ID_MAX_LENGTH)
    private String pushId;

    //项目名称
    @NotBlank(message = SystemErrCode.ERRMSG_PROJECT_NAME_NOT_NULL)
    @Length(max = 100, message = SystemErrCode.ERRMSG_PROJECT_NAME_MAX_LENGTH)
    private String project;

    //服务组名称
    @NotBlank(message = SystemErrCode.ERRMSG_CLUSTER_NAME_NOT_NULL)
    @Length(max = 100, message = SystemErrCode.ERRMSG_CLUSTER_NAME_MAX_LENGTH)
    private String group;

    //服务名称
    @NotBlank(message = SystemErrCode.ERRMSG_CLUSTER_NAME_NOT_NULL)
    @Length(max = 100, message = SystemErrCode.ERRMSG_SERVICE_NAME_MAX_LENGTH)
    private String service;

    //版本号
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_VERSION_NOT_NULL)
    @Length(max = 20, message = SystemErrCode.ERRMSG_SERVICE_VERSION_MAX_LENGTH)
    private String version;

    //配置名称
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_CONFIG_NAME_NOT_NULL)
    @Length(max = 100, message = SystemErrCode.ERRMSG_SERVICE_CONFIG_NAME_MAX_LENGTH)
    private String config;

    //地址
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_CONFIG_ADDR_NOT_NULL)
    @Length(max = 100, message = SystemErrCode.ERRMSG_SERVICE_CONFIG_ADDR_MAX_LENGTH)
    private String addr;

    //更新状态
    private int updateStatus;

    //更新时间
    private Date updateTime;

    //加载状态
    private int loadStatus;

    //加载时间
    private Date loadTime;

    //灰度组标识
    private String grayGroupId;

    public String getPushId() {
        return pushId;
    }

    public void setPushId(String pushId) {
        this.pushId = pushId;
    }

    public String getProject() {
        return project;
    }

    public void setProject(String project) {
        this.project = project;
    }

    public String getGroup() {
        return group;
    }

    public void setGroup(String group) {
        this.group = group;
    }

    public String getService() {
        return service;
    }

    public void setService(String service) {
        this.service = service;
    }

    public String getVersion() {
        return version;
    }

    public void setVersion(String version) {
        this.version = version;
    }

    public String getConfig() {
        return config;
    }

    public void setConfig(String config) {
        this.config = config;
    }

    public String getAddr() {
        return addr;
    }

    public void setAddr(String addr) {
        this.addr = addr;
    }

    public int getUpdateStatus() {
        return updateStatus;
    }

    public void setUpdateStatus(int updateStatus) {
        this.updateStatus = updateStatus;
    }

    public Date getUpdateTime() {
        return updateTime;
    }

    public void setUpdateTime(Date updateTime) {
        this.updateTime = updateTime;
    }

    public int getLoadStatus() {
        return loadStatus;
    }

    public void setLoadStatus(int loadStatus) {
        this.loadStatus = loadStatus;
    }

    public Date getLoadTime() {
        return loadTime;
    }

    public void setLoadTime(Date loadTime) {
        this.loadTime = loadTime;
    }

    public String getGrayGroupId() {
        return grayGroupId;
    }

    public void setGrayGroupId(String grayGroupId) {
        this.grayGroupId = grayGroupId;
    }

    @Override
    public String toString() {
        return "ServiceConfigFeedBackRequestBody{" +
                "pushId='" + pushId + '\'' +
                ", project='" + project + '\'' +
                ", group='" + group + '\'' +
                ", service='" + service + '\'' +
                ", version='" + version + '\'' +
                ", config='" + config + '\'' +
                ", addr='" + addr + '\'' +
                ", updateStatus=" + updateStatus +
                ", updateTime=" + updateTime +
                ", loadStatus=" + loadStatus +
                ", loadTime=" + loadTime +
                '}';
    }
}
