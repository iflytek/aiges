package com.iflytek.ccr.polaris.cynosure.request.serviceversion;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import org.hibernate.validator.constraints.Length;
import org.hibernate.validator.constraints.NotBlank;

import java.io.Serializable;

/**
 * 复制服务版本-请求
 *
 * @author sctang2
 * @create 2017-11-17 16:04
 **/
public class CopyServiceVersionRequestBody implements Serializable {
    private static final long serialVersionUID = 2388185789083489433L;

    //服务id
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_ID_NOT_NULL)
    @Length(max = 50, message = SystemErrCode.ERRMSG_SERVICE_ID_MAX_LENGTH)
    private String serviceId;

    //被复制的版本的id
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_VERSION_ID_NOT_NULL)
    @Length(max = 50, message = SystemErrCode.ERRMSG_SERVICE_VERSION_ID_MAX_LENGTH)
    private String oldVersionId;

    //新的版本号
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_VERSION_NOT_NULL)
    @Length(min = 1, max = 20, message = SystemErrCode.ERRMSG_SERVICE_VERSION_MAX_LENGTH)
    private String version;

    //新的版本描述
    @Length(max = 500, message = SystemErrCode.ERRMSG_SERVICE_VERSION_DESC_MAX_LENGTH)
    private String desc;

    public CopyServiceVersionRequestBody() {
    }

    public CopyServiceVersionRequestBody(String serviceId, String oldVersionId, String version, String desc) {
        this.serviceId = serviceId;
        this.oldVersionId = oldVersionId;
        this.version = version;
        this.desc = desc;
    }

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

    public String getOldVersionId() {
        return oldVersionId;
    }

    public void setOldVersionId(String oldVersionId) {
        this.oldVersionId = oldVersionId;
    }

    @Override
    public String toString() {
        return "CopyServiceVersionRequestBody{" +
                "version='" + version + '\'' +
                ", desc='" + desc + '\'' +
                ", oldVersionId='" + oldVersionId + '\'' +
                ", serviceId='" + serviceId + '\'' +
                '}';
    }
}
