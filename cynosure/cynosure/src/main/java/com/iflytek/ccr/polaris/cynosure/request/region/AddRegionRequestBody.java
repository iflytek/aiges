package com.iflytek.ccr.polaris.cynosure.request.region;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import org.hibernate.validator.constraints.Length;
import org.hibernate.validator.constraints.NotBlank;
import org.hibernate.validator.constraints.URL;

import java.io.Serializable;

/**
 * 新增区域-请求
 *
 * @author sctang2
 * @create 2017-11-14 20:28
 **/
public class AddRegionRequestBody implements Serializable {
    private static final long serialVersionUID = -8830365021004310539L;

    //区域名称
    @NotBlank(message = SystemErrCode.ERRMSG_REGION_NAME_NOT_NULL)
    @Length(min = 1, max = 100, message = SystemErrCode.ERRMSG_REGION_NAME_MAX_LENGTH)
    private String name;

    //推送地址
    @NotBlank(message = SystemErrCode.ERRMSG_REGION_PUSH_URL_NOT_NULL)
    @Length(min = 1, max = 500, message = SystemErrCode.ERRMSG_REGION_PUSH_URL_MAX_LENGTH)
    @URL(message = SystemErrCode.ERRMSG_REGION_PUSH_URL_INVALID)
    private String pushUrl;

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getPushUrl() {
        return pushUrl;
    }

    public void setPushUrl(String pushUrl) {
        this.pushUrl = pushUrl;
    }

    @Override
    public String toString() {
        return "AddRegionRequestBody{" +
                "name='" + name + '\'' +
                ", pushUrl='" + pushUrl + '\'' +
                '}';
    }
}
