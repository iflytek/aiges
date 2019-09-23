package com.iflytek.ccr.polaris.cynosure.request.project;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import org.hibernate.validator.constraints.Length;
import org.hibernate.validator.constraints.NotBlank;

import java.io.Serializable;

/**
 * 编辑项目-请求
 *
 * @author sctang2
 * @create 2017-11-19 21:07
 **/
public class EditProjectRequestBody implements Serializable {
    private static final long serialVersionUID = 6695885100507321754L;

    //项目id
    @NotBlank(message = SystemErrCode.ERRMSG_PROJECT_ID_NOT_NULL)
    private String id;

    //项目描述
    @Length(max = 500, message = SystemErrCode.ERRMSG_PROJECT_DESC_MAX_LENGTH)
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
        return "EditProjectRequestBody{" +
                "id='" + id + '\'' +
                ", desc='" + desc + '\'' +
                '}';
    }
}
