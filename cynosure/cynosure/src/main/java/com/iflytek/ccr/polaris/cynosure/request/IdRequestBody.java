package com.iflytek.ccr.polaris.cynosure.request;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;

import io.swagger.annotations.ApiModelProperty;

import org.hibernate.validator.constraints.NotBlank;

import java.io.Serializable;

/**
 * id-请求
 *
 * @author sctang2
 * @create 2017-11-16 9:18
 **/
public class IdRequestBody implements Serializable {
    private static final long serialVersionUID = 7866227766357729385L;

    //唯一标识
    @NotBlank(message = SystemErrCode.ERRMSG_ID_NOT_NULL)
    @ApiModelProperty("唯一标识")
    private String id;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    @Override
    public String toString() {
        return "IdRequestBody{" +
                "id='" + id + '\'' +
                '}';
    }
}
