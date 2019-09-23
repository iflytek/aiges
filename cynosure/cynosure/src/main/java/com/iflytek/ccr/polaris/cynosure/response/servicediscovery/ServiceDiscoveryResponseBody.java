package com.iflytek.ccr.polaris.cynosure.response.servicediscovery;

import com.iflytek.ccr.polaris.cynosure.request.servicediscovery.RouteRule;

import java.io.Serializable;
import java.util.List;

/**
 * 编辑服务发现-响应
 *
 * @author sctang2
 * @create 2017-12-07 18:56
 **/
public class ServiceDiscoveryResponseBody implements Serializable {
    private static final long serialVersionUID = -162640371432002294L;

    private String apiVersion;

    private String region;

    //负载均衡名称
    private String loadbalance;

    //自定义参数
    private String params;

    //路由规则设置
    private List<RouteRule> routeRules;

    public String getApiVersion() {
        return apiVersion;
    }

    public void setApiVersion(String apiVersion) {
        this.apiVersion = apiVersion;
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

    public String getParams() {
        return params;
    }

    public void setParams(String params) {
        this.params = params;
    }

    public List<RouteRule> getRouteRules() {
        return routeRules;
    }

    public void setRouteRules(List<RouteRule> routeRules) {
        this.routeRules = routeRules;
    }

    @Override
    public String toString() {
        return "ServiceDiscoveryResponseBody{" +
                "apiVersion='" + apiVersion + '\'' +
                ", region='" + region + '\'' +
                ", loadbalance='" + loadbalance + '\'' +
                ", params='" + params + '\'' +
                ", routeRules=" + routeRules +
                '}';
    }
}
