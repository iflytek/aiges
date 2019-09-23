package com.iflytek.ccr.polaris.cynosure.request.region;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import org.hibernate.validator.constraints.Length;
import org.hibernate.validator.constraints.NotBlank;
import org.hibernate.validator.constraints.URL;

import java.io.Serializable;

/**
 * 编辑区域-请求
 *
 * @author sctang2
 * @create 2017-11-14 20:44
 **/
public class EditRegionRequestBody implements Serializable {
    private static final long serialVersionUID = 7266645362328469103L;

    //区域id
    @NotBlank(message = SystemErrCode.ERRMSG_REGION_ID_NOT_NULL)
    private String id;

    //推送地址
    @NotBlank(message = SystemErrCode.ERRMSG_REGION_PUSH_URL_NOT_NULL)
    @Length(min = 1, max = 500, message = SystemErrCode.ERRMSG_REGION_PUSH_URL_MAX_LENGTH)
    @URL(message = SystemErrCode.ERRMSG_REGION_PUSH_URL_INVALID)
    private String pushUrl;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getPushUrl() {
        return pushUrl;
    }

    public void setPushUrl(String pushUrl) {
        this.pushUrl = pushUrl;
    }

    @Override
    public String toString() {
        return "EditRegionRequestBody{" +
                "id='" + id + '\'' +
                ", pushUrl='" + pushUrl + '\'' +
                '}';
    }
}
