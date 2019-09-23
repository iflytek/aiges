package com.iflytek.ccr.polaris.cynosure.request.servicediscovery;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import org.hibernate.validator.constraints.Length;
import org.hibernate.validator.constraints.NotBlank;

import java.io.Serializable;
import java.util.Date;

/**
 * 服务发现更新反馈-请求
 *
 * @author sctang2
 * @create 2017-12-12 19:13
 **/
public class ServiceDiscoveryFeedBackRequestBody implements Serializable {
    private static final long serialVersionUID = 9075997934616762554L;

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

    //消费端名称
    @NotBlank(message = SystemErrCode.ERRMSG_CLUSTER_NAME_NOT_NULL)
    @Length(max = 100, message = SystemErrCode.ERRMSG_SERVICE_NAME_MAX_LENGTH)
    private String consumer;

    //提供端名称
//    @NotBlank(message = SystemErrCode.ERRMSG_CLUSTER_NAME_NOT_NULL)
    @Length(max = 100, message = SystemErrCode.ERRMSG_SERVICE_NAME_MAX_LENGTH)
    private String provider;

    //地址
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_CONFIG_ADDR_NOT_NULL)
    @Length(max = 100, message = SystemErrCode.ERRMSG_SERVICE_CONFIG_ADDR_MAX_LENGTH)
    private String addr;

    //变化类型
    private String type;

    //消费端版本号
    @Length(max = 20, message = SystemErrCode.ERRMSG_SERVICE_VERSION_MAX_LENGTH)
    private String consumerVersion;

    //提供端版本号
    @Length(max = 20, message = SystemErrCode.ERRMSG_SERVICE_VERSION_MAX_LENGTH)
    private String providerVersion;

    private String apiVersion;

    //更新状态
    private int updateStatus;

    //更新时间
    private Date updateTime;

    //加载状态
    private int loadStatus;

    //加载时间
    private Date loadTime;

    public String getType() {
        return type;
    }

    public void setType(String type) {
        this.type = type;
    }

    public String getApiVersion() {
        return apiVersion;
    }

    public void setApiVersion(String apiVersion) {
        this.apiVersion = apiVersion;
    }

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

    public String getConsumer() {
        return consumer;
    }

    public void setConsumer(String consumer) {
        this.consumer = consumer;
    }

    public String getProvider() {
        return provider;
    }

    public void setProvider(String provider) {
        this.provider = provider;
    }

    public String getAddr() {
        return addr;
    }

    public void setAddr(String addr) {
        this.addr = addr;
    }

    public String getConsumerVersion() {
        return consumerVersion;
    }

    public void setConsumerVersion(String consumerVersion) {
        this.consumerVersion = consumerVersion;
    }

    public String getProviderVersion() {
        return providerVersion;
    }

    public void setProviderVersion(String providerVersion) {
        this.providerVersion = providerVersion;
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

    @Override
    public String toString() {
        return "ServiceDiscoveryFeedBackRequestBody{" +
                "pushId='" + pushId + '\'' +
                ", project='" + project + '\'' +
                ", group='" + group + '\'' +
                ", consumer='" + consumer + '\'' +
                ", provider='" + provider + '\'' +
                ", addr='" + addr + '\'' +
                ", consumerVersion='" + consumerVersion + '\'' +
                ", providerVersion='" + providerVersion + '\'' +
                ", updateStatus=" + updateStatus +
                ", updateTime=" + updateTime +
                ", loadStatus=" + loadStatus +
                ", loadTime=" + loadTime +
                '}';
    }
}
