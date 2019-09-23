package com.iflytek.ccr.finder.value;

import java.io.Serializable;
import java.util.List;

/**
 * 服务对象
 */
public class Service implements Serializable {

    /**
     * api版本号
     */
    private String apiVersion;

    /**
     * 服务名称
     */
    private String name;

    /**
     * 服务实例对象
     */
    private List<ServiceInstance> serverList;

    /**
     * 服务配置对象
     */
    private ServiceConfig serviceConfig;

    public String getApiVersion() {
        return apiVersion;
    }

    public void setApiVersion(String apiVersion) {
        this.apiVersion = apiVersion;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public List<ServiceInstance> getServerList() {
        return serverList;
    }

    public void setServerList(List<ServiceInstance> serverList) {
        this.serverList = serverList;
    }

    public ServiceConfig getServiceConfig() {
        return serviceConfig;
    }

    public void setServiceConfig(ServiceConfig serviceConfig) {
        this.serviceConfig = serviceConfig;
    }

    @Override
    public String toString() {
        return "Service{" +
                "apiVersion='" + apiVersion + '\'' +
                ", name='" + name + '\'' +
                ", serverList=" + serverList +
                ", serviceConfig=" + serviceConfig +
                '}';
    }
}
