package com.iflytek.ccr.polaris.cynosure.request.graygroup;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import org.hibernate.validator.constraints.Length;
import org.hibernate.validator.constraints.NotBlank;

import java.io.Serializable;

/**
 * 新增灰度组-请求
 *
 * @author sctang2
 * @create 2017-11-21 14:21
 **/
public class QueryGrayGroupListRequestBody implements Serializable {
    private static final long serialVersionUID = 4229119834922797764L;

    //版本id
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_VERSION_ID_NOT_NULL)
    @Length(max = 50, message = SystemErrCode.ERRMSG_SERVICE_VERSION_ID_MAX_LENGTH)
    private String versionId;

    public String getVersionId() {
        return versionId;
    }

    public void setVersionId(String versionId) {
        this.versionId = versionId;
    }

    @Override
    public String toString() {
        return "QueryGrayGroupListRequestBody{" +
                "versionId='" + versionId + '\'' +
                '}';
    }
}
