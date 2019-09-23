package com.iflytek.ccr.polaris.cynosure.request.servicediscovery;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.request.BaseRequestBody;
import org.hibernate.validator.constraints.NotBlank;

import java.io.Serializable;

/**
 * 查询服务发现明细-请求
 *
 * @author sctang2
 * @create 2017-12-05 18:58
 **/
public class QueryServiceDiscoveryDetailRequestBody extends BaseRequestBody implements Serializable {
    private static final long serialVersionUID = -4143507557040652713L;

    //项目名称
    @NotBlank(message = SystemErrCode.ERRMSG_PROJECT_NAME_NOT_NULL)
    private String project;

    //集群名称
    @NotBlank(message = SystemErrCode.ERRMSG_CLUSTER_NAME_NOT_NULL)
    private String cluster;

    //服务名
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_NAME_NOT_NULL)
    private String service;

    //服务版本
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_VERSION_NOT_NULL)
    private String apiVersion;

    //区域名称
    @NotBlank(message = SystemErrCode.ERRMSG_REGION_NAME_NOT_NULL)
    private String region;

    public String getProject() {
        return project;
    }

    public void setProject(String project) {
        this.project = project;
    }

    public String getCluster() {
        return cluster;
    }

    public void setCluster(String cluster) {
        this.cluster = cluster;
    }

    public String getService() {
        return service;
    }

    public void setService(String service) {
        this.service = service;
    }

    public String getRegion() {
        return region;
    }

    public void setRegion(String region) {
        this.region = region;
    }

    public String getApiVersion() {
        return apiVersion;
    }

    public void setApiVersion(String apiVersion) {
        this.apiVersion = apiVersion;
    }
}
