package com.iflytek.ccr.polaris.cynosure.request.servicediscovery;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import org.hibernate.validator.constraints.NotBlank;

import java.io.Serializable;
import java.util.List;

/**
 * 编辑服务发现-请求
 *
 * @author sctang2
 * @create 2017-12-07 18:56
 **/
public class EditServiceDiscoveryRequestBody implements Serializable {
    private static final long serialVersionUID = -162640371432002294L;

    //项目名称
    @NotBlank(message = SystemErrCode.ERRMSG_PROJECT_NAME_NOT_NULL)
    private String project;

    //集群名称
    @NotBlank(message = SystemErrCode.ERRMSG_CLUSTER_NAME_NOT_NULL)
    private String cluster;

    //服务名称
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_NAME_NOT_NULL)
    private String service;

    //服务版本
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_VERSION_NOT_NULL)
    private String apiVersion;

    //区域名称
    @NotBlank(message = SystemErrCode.ERRMSG_REGION_NAME_NOT_NULL)
    private String region;

    //负载均衡名称
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_DISCOVERY_LOADBALANCE_NOT_NULL)
    private String loadbalance;

    //自定义参数
    private List<ServiceParam> params;

    //路由规则设置
    private List<RouteRule> routeRules;

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

    public String getLoadbalance() {
        return loadbalance;
    }

    public void setLoadbalance(String loadbalance) {
        this.loadbalance = loadbalance;
    }

    public List<ServiceParam> getParams() {
        return params;
    }

    public void setParams(List<ServiceParam> params) {
        this.params = params;
    }

    public List<RouteRule> getRouteRules() {
        return routeRules;
    }

    public void setRouteRules(List<RouteRule> routeRules) {
        this.routeRules = routeRules;
    }

    public String getApiVersion() {
        return apiVersion;
    }

    public void setApiVersion(String apiVersion) {
        this.apiVersion = apiVersion;
    }
}