package com.iflytek.ccr.polaris.cynosure.request.InstanceManageRequestBody;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import org.hibernate.validator.constraints.Length;
import org.hibernate.validator.constraints.NotBlank;

import java.io.Serializable;

/**
 * 编辑推送实例-请求
 *
 * @author sctang2
 * @create 2017-11-21 14:21
 **/
public class EditInstanceRequestBody implements Serializable {
    private static final long serialVersionUID = 4219119634922797764L;

    //灰度组id
    @NotBlank(message = SystemErrCode.ERRMSG_GRAY_GROUP_ID_NOT_NULL)
    @Length(max = 50, message = SystemErrCode.ERRMSG_GRAY_GROUP_ID_MAX_LENGTH)
    private String grayId;

    //版本id
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_VERSION_ID_NOT_NULL)
    @Length(max = 50, message = SystemErrCode.ERRMSG_SERVICE_VERSION_ID_MAX_LENGTH)
    private String versionId;

    //推送实例列表
    private String content;

    public String getGrayId() {
        return grayId;
    }

    public void setGrayId(String grayId) {
        this.grayId = grayId;
    }

    public String getVersionId() {
        return versionId;
    }

    public void setVersionId(String versionId) {
        this.versionId = versionId;
    }

    public String getContent() {
        return content;
    }

    public void setContent(String content) {
        this.content = content;
    }

    @Override
    public String toString() {
        return "EditInstanceRequestBody{" +
                "grayId='" + grayId + '\'' +
                ", versionId='" + versionId + '\'' +
                ", content='" + content + '\'' +
                '}';
    }
}
