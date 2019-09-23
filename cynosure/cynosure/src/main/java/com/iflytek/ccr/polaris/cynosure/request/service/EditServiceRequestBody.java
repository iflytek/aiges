package com.iflytek.ccr.polaris.cynosure.request.service;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import org.hibernate.validator.constraints.Length;
import org.hibernate.validator.constraints.NotBlank;

import java.io.Serializable;

/**
 * 编辑服务-请求
 *
 * @author sctang2
 * @create 2017-11-16 19:46
 **/
public class EditServiceRequestBody implements Serializable {
    private static final long serialVersionUID = -8223387895350456410L;

    //服务id
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_ID_NOT_NULL)
    private String id;

    //服务描述
    @Length(max = 500, message = SystemErrCode.ERRMSG_SERVICE_DESC_MAX_LENGTH)
    private String desc;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getDesc() {
        return desc;
    }

    public void setDesc(String desc) {
        this.desc = desc;
    }

    @Override
    public String toString() {
        return "EditServiceRequestBody{" +
                "id='" + id + '\'' +
                ", desc='" + desc + '\'' +
                '}';
    }
}
