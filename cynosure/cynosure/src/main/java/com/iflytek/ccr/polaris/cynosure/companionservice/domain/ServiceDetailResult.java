package com.iflytek.ccr.polaris.cynosure.companionservice.domain;

import java.io.Serializable;

/**
 * 服务明细结果
 *
 * @author sctang2
 * @create 2018-02-06 10:40
 **/
public class ServiceDetailResult implements Serializable {
    private static final long serialVersionUID = 5150271535369434190L;

    //集群名称
    private String name;

    //服务名
    private String apiVersion;

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getApiVersion() {
        return apiVersion;
    }

    public void setApiVersion(String apiVersion) {
        this.apiVersion = apiVersion;
    }

    @Override
    public String toString() {
        return "ServiceDetailResult{" +
                "name='" + name + '\'' +
                ", apiVersion=" + apiVersion +
                '}';
    }
}
