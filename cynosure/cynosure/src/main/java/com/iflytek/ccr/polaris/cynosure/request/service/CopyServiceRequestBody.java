package com.iflytek.ccr.polaris.cynosure.request.service;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import org.hibernate.validator.constraints.Length;
import org.hibernate.validator.constraints.NotBlank;

import java.io.Serializable;

/**
 * 复制服务-请求
 *
 * @author sctang2
 * @create 2017-11-17 16:04
 **/
public class CopyServiceRequestBody implements Serializable {
    private static final long serialVersionUID = 2388285789083489433L;

    //复制出新的服务所属的集群id
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_ID_NOT_NULL)
    @Length(max = 50, message = SystemErrCode.ERRMSG_SERVICE_ID_MAX_LENGTH)
    private String clusterId;

    //被复制的服务的id
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_VERSION_ID_NOT_NULL)
    @Length(max = 50, message = SystemErrCode.ERRMSG_SERVICE_VERSION_ID_MAX_LENGTH)
    private String oldServiceId;

    //服务名称
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_VERSION_NOT_NULL)
    @Length(min = 1, max = 20, message = SystemErrCode.ERRMSG_SERVICE_VERSION_MAX_LENGTH)
    private String serviceName;

    //服务描述
    @Length(max = 500, message = SystemErrCode.ERRMSG_SERVICE_VERSION_DESC_MAX_LENGTH)
    private String desc;

    public CopyServiceRequestBody() {
    }

    public CopyServiceRequestBody(String clusterId, String oldServiceId, String serviceName, String desc) {
        this.clusterId = clusterId;
        this.oldServiceId = oldServiceId;
        this.serviceName = serviceName;
        this.desc = desc;
    }

    public String getServiceName() {
        return serviceName;
    }

    public void setServiceName(String serviceName) {
        this.serviceName = serviceName;
    }

    public String getDesc() {
        return desc;
    }

    public void setDesc(String desc) {
        this.desc = desc;
    }

    public String getOldServiceId() {
        return oldServiceId;
    }

    public void setOldServiceId(String oldServiceId) {
        this.oldServiceId = oldServiceId;
    }

    public String getClusterId() {
        return clusterId;
    }

    public void setClusterId(String clusterId) {
        this.clusterId = clusterId;
    }

    @Override
    public String toString() {
        return "CopyServiceRequestBody{" +
                "serviceName='" + serviceName + '\'' +
                ", desc='" + desc + '\'' +
                ", oldServiceId='" + oldServiceId + '\'' +
                ", clusterId='" + clusterId + '\'' +
                '}';
    }
}
