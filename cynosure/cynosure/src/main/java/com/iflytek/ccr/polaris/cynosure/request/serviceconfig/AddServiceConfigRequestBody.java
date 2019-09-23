package com.iflytek.ccr.polaris.cynosure.request.serviceconfig;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import org.hibernate.validator.constraints.Length;
import org.hibernate.validator.constraints.NotBlank;

import java.io.Serializable;

/**
 * 新增服务配置-请求
 *
 * @author sctang2
 * @create 2017-11-21 14:21
 **/
public class AddServiceConfigRequestBody implements Serializable {
    private static final long serialVersionUID = 4219119734922767763L;

    //描述
    @Length(max = 500, message = SystemErrCode.ERRMSG_SERVICE_CONFIG_DESC_MAX_LENGTH)
    private String desc;

    //版本id
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_VERSION_ID_NOT_NULL)
    @Length(max = 50, message = SystemErrCode.ERRMSG_SERVICE_VERSION_ID_MAX_LENGTH)
    private String versionId;

    public String getDesc() {
        return desc;
    }

    public void setDesc(String desc) {
        this.desc = desc;
    }

    public String getVersionId() {
        return versionId;
    }

    public void setVersionId(String versionId) {
        this.versionId = versionId;
    }

    @Override
    public String toString() {
        return "AddServiceConfigRequestBody{" +
                "desc='" + desc + '\'' +
                ", versionId='" + versionId + '\'' +
                '}';
    }
}
