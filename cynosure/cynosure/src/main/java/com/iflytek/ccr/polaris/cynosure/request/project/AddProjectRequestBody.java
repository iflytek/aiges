package com.iflytek.ccr.polaris.cynosure.request.project;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import org.hibernate.validator.constraints.Length;
import org.hibernate.validator.constraints.NotBlank;

import java.io.Serializable;

/**
 * 添加项目-请求
 *
 * @author sctang2
 * @create 2017-11-19 21:07
 **/
public class AddProjectRequestBody implements Serializable {
    private static final long serialVersionUID = -2999333987452302789L;

    //项目名称
    @NotBlank(message = SystemErrCode.ERRMSG_PROJECT_NAME_NOT_NULL)
    @Length(min = 1, max = 100, message = SystemErrCode.ERRMSG_PROJECT_NAME_MAX_LENGTH)
    private String name;

    //项目描述
    @Length(max = 500, message = SystemErrCode.ERRMSG_PROJECT_DESC_MAX_LENGTH)
    private String desc;

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getDesc() {
        return desc;
    }

    public void setDesc(String desc) {
        this.desc = desc;
    }

    @Override
    public String toString() {
        return "AddProjectRequestBody{" +
                "name='" + name + '\'' +
                ", desc='" + desc + '\'' +
                '}';
    }
}
