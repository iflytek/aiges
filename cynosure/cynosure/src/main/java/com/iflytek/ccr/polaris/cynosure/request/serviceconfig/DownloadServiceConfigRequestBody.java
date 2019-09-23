package com.iflytek.ccr.polaris.cynosure.request.serviceconfig;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import org.hibernate.validator.constraints.NotBlank;

import java.io.Serializable;

/**
 * 下载配置文件请求类
 * @create by ygli3 2018.08.15
 */
public class DownloadServiceConfigRequestBody implements Serializable {
    private static final long serialVersionUID = 891356528565431542L;

    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_CONFIG_ID_NOT_NULL)
    private String id;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

}
