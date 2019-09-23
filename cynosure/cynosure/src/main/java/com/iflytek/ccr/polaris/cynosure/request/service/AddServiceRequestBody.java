package com.iflytek.ccr.polaris.cynosure.request.service;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import org.hibernate.validator.constraints.Length;
import org.hibernate.validator.constraints.NotBlank;

import java.io.Serializable;

/**
 * 新增服务-请求
 *
 * @author sctang2
 * @create 2017-11-16 19:45
 **/
public class AddServiceRequestBody implements Serializable {
    private static final long serialVersionUID = -183073552397679807L;

    //服务名称
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_NAME_NOT_NULL)
    @Length(min = 1, max = 100, message = SystemErrCode.ERRMSG_SERVICE_NAME_MAX_LENGTH)
    private String name;

    //服务描述
    @Length(max = 500, message = SystemErrCode.ERRMSG_SERVICE_DESC_MAX_LENGTH)
    private String desc;

    //集群id
    @NotBlank(message = SystemErrCode.ERRMSG_CLUSTER_ID_NOT_NULL)
    @Length(max = 50, message = SystemErrCode.ERRMSG_CLUSTER_ID_MAX_LENGTH)
    private String clusterId;

    public AddServiceRequestBody() {
    }

    public AddServiceRequestBody(String name, String desc, String clusterId) {
        this.name = name;
        this.desc = desc;
        this.clusterId = clusterId;
    }

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

    public String getClusterId() {
        return clusterId;
    }

    public void setClusterId(String clusterId) {
        this.clusterId = clusterId;
    }

    @Override
    public String toString() {
        return "AddServiceRequestBody{" +
                "name='" + name + '\'' +
                ", desc='" + desc + '\'' +
                ", clusterId='" + clusterId + '\'' +
                '}';
    }
}
