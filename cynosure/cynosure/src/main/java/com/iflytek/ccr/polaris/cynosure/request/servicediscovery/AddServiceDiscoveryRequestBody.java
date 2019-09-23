package com.iflytek.ccr.polaris.cynosure.request.servicediscovery;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import org.hibernate.validator.constraints.Length;
import org.hibernate.validator.constraints.NotBlank;

import java.io.Serializable;

/**
 * 新增服务发现-请求
 *
 * @author sctang2
 * @create 2018-02-07 13:37
 **/
public class AddServiceDiscoveryRequestBody implements Serializable {
    private static final long serialVersionUID = 1769085590367739510L;

    //项目名称
    @NotBlank(message = SystemErrCode.ERRMSG_PROJECT_NAME_NOT_NULL)
    @Length(min = 1, max = 100, message = SystemErrCode.ERRMSG_PROJECT_NAME_MAX_LENGTH)
    private String project;

    //集群名称
    @NotBlank(message = SystemErrCode.ERRMSG_CLUSTER_NAME_NOT_NULL)
    @Length(min = 1, max = 100, message = SystemErrCode.ERRMSG_CLUSTER_NAME_MAX_LENGTH)
    private String group;

    //服务名称
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_NAME_NOT_NULL)
    @Length(min = 1, max = 100, message = SystemErrCode.ERRMSG_SERVICE_NAME_MAX_LENGTH)
    private String service;

    //版本号
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICEAPI_VERSION_NOT_NULL)
    @Length(min = 1, max = 20, message = SystemErrCode.ERRMSG_SERVICE_VERSION_MAX_LENGTH)
    private String apiVersion;

    public String getProject() {
        return project;
    }

    public void setProject(String project) {
        this.project = project;
    }

    public String getGroup() {
        return group;
    }

    public void setGroup(String group) {
        this.group = group;
    }

    public String getService() {
        return service;
    }

    public void setService(String service) {
        this.service = service;
    }

    public String getApiVersion() {
        return apiVersion;
    }

    public void setApiVersion(String apiVersion) {
        this.apiVersion = apiVersion;
    }

    @Override
    public String toString() {
        return "AddServiceDiscoveryRequestBody{" +
                "project='" + project + '\'' +
                ", group='" + group + '\'' +
                ", service='" + service + '\'' +
                ", apiVersion='" + apiVersion + '\'' +
                '}';
    }
}
