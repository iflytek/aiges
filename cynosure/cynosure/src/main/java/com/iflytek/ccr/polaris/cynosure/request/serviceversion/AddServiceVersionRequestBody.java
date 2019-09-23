package com.iflytek.ccr.polaris.cynosure.request.serviceversion;

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
public class AddServiceVersionRequestBody implements Serializable {
    private static final long serialVersionUID = 2388185789083489433L;

    //版本号
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_VERSION_NOT_NULL)
    @Length(min = 1, max = 20, message = SystemErrCode.ERRMSG_SERVICE_VERSION_MAX_LENGTH)
    private String version;

    //版本描述
    @Length(max = 500, message = SystemErrCode.ERRMSG_SERVICE_VERSION_DESC_MAX_LENGTH)
    private String desc;

    //服务id
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_ID_NOT_NULL)
    @Length(max = 50, message = SystemErrCode.ERRMSG_SERVICE_ID_MAX_LENGTH)
    private String serviceId;

    public String getVersion() {
        return version;
    }

    public void setVersion(String version) {
        this.version = version;
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

    public AddServiceVersionRequestBody() {
    }

    public AddServiceVersionRequestBody(String version, String desc, String serviceId) {
        this.version = version;
        this.desc = desc;
        this.serviceId = serviceId;
    }

    @Override
    public String toString() {
        return "AddServiceVersionRequestBody{" +
                "version='" + version + '\'' +
                ", desc='" + desc + '\'' +
                ", serviceId='" + serviceId + '\'' +
                '}';
    }
}
