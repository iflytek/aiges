package com.iflytek.ccr.polaris.cynosure.request.graygroup;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import org.hibernate.validator.constraints.Length;
import org.hibernate.validator.constraints.NotBlank;

import java.io.Serializable;

/**
 * 复制灰度组-请求
 *
 * @author sctang2
 * @create 2017-11-21 14:21
 **/
public class CopyGrayGroupRequestBody implements Serializable {
    private static final long serialVersionUID = 4229119734922797764L;

    //版本id
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_VERSION_ID_NOT_NULL)
    @Length(max = 50, message = SystemErrCode.ERRMSG_SERVICE_VERSION_ID_MAX_LENGTH)
    private String versionId;

    @Override
    public String toString() {
        return "CopyGrayGroupRequestBody{" +
                "versionId='" + versionId + '\'' +
                '}';
    }
}
