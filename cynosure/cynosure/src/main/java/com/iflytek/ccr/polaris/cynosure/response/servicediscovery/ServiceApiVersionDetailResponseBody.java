package com.iflytek.ccr.polaris.cynosure.response.servicediscovery;

import java.io.Serializable;
import java.util.Date;

/**
 * 服务版本明细-响应
 *
 * @author sctang2
 * @create 2017-11-17 16:13
 **/
public class ServiceApiVersionDetailResponseBody implements Serializable {
    private static final long serialVersionUID = -5936440506999176691L;

    //版本id
    private String id;

    //版本号
    private String apiVersion;

    //服务id
    private String serviceId;

    //版本描述
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

    public String getApiVersion() {
        return apiVersion;
    }

    public void setApiVersion(String apiVersion) {
        this.apiVersion = apiVersion;
    }

    public String getServiceId() {
        return serviceId;
    }

    public void setServiceId(String serviceId) {
        this.serviceId = serviceId;
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
        return "ServiceApiVersionDetailResponseBody{" +
                "id='" + id + '\'' +
                ", apiVersion='" + apiVersion + '\'' +
                ", serviceId='" + serviceId + '\'' +
                ", desc='" + desc + '\'' +
                ", createTime=" + createTime +
                ", updateTime=" + updateTime +
                '}';
    }
}
