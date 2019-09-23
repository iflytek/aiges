package com.iflytek.ccr.polaris.cynosure.request.servicediscovery;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import org.hibernate.validator.constraints.Length;
import org.hibernate.validator.constraints.NotBlank;

import java.io.Serializable;

/**
 * 新增服务版本-请求
 *
 * @author sctang2
 * @create 2017-11-17 16:04
 **/
public class AddServiceApiVersionRequestBody implements Serializable {
    private static final long serialVersionUID = 2388185789083489433L;

    //版本号
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICEAPI_VERSION_NOT_NULL)
    @Length(min = 1, max = 20, message = SystemErrCode.ERRMSG_SERVICE_VERSION_MAX_LENGTH)
    private String apiVersion;

    //版本描述
    @Length(max = 500, message = SystemErrCode.ERRMSG_SERVICE_VERSION_DESC_MAX_LENGTH)
    private String desc;

    //服务id
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_ID_NOT_NULL)
    @Length(max = 50, message = SystemErrCode.ERRMSG_SERVICE_ID_MAX_LENGTH)
    private String serviceId;

    public String getApiVersion() {
        return apiVersion;
    }

    public void setApiVersion(String apiVersion) {
        this.apiVersion = apiVersion;
    }

    public String getDesc() {
        return desc;
    }

    public void setDesc(String desc) {
        this.desc = desc;
    }

    public String getServiceId() {
        return serviceId;
    }

    public void setServiceId(String serviceId) {
        this.serviceId = serviceId;
    }

    @Override
    public String toString() {
        return "AddServiceApiVersionRequestBody{" +
                "apiVersion='" + apiVersion + '\'' +
                ", desc='" + desc + '\'' +
                ", serviceId='" + serviceId + '\'' +
                '}';
    }
}
