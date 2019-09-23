package com.iflytek.ccr.polaris.cynosure.response.servicediscovery;

import java.io.Serializable;

/**
 * 查询服务发现-响应
 *
 * @author sctang2
 * @create 2018-02-01 14:48
 **/
public class QueryServiceDiscoveryResponseBody implements Serializable {
    private static final long serialVersionUID = -6737149774466048184L;

    //服务名称
    private String apiVersion;

    //区域名称
    private String region;

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

    @Override
    public String toString() {
        return "QueryServiceDiscoveryResponseBody{" +
                "apiVersion='" + apiVersion + '\'' +
                ", region='" + region + '\'' +
                '}';
    }
}
