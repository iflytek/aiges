package com.iflytek.ccr.polaris.cynosure.response.serviceversion;

import java.io.Serializable;
import java.util.Date;

/**
 * 服务版本明细-响应
 *
 * @author sctang2
 * @create 2017-11-17 16:13
 **/
public class ServiceVersionDetailResponseBody implements Serializable {
    private static final long serialVersionUID = -5936440506999176691L;

    //版本id
    private String id;

    //版本号
    private String version;

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

    public String getVersion() {
        return version;
    }

    public void setVersion(String version) {
        this.version = version;
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
        return "ServiceVersionDetailResponseBody{" +
                "id='" + id + '\'' +
                ", version='" + version + '\'' +
                ", serviceId='" + serviceId + '\'' +
                ", desc='" + desc + '\'' +
                ", createTime=" + createTime +
                ", updateTime=" + updateTime +
                '}';
    }
}
