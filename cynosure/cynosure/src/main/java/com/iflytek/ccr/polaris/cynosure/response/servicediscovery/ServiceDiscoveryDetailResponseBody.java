package com.iflytek.ccr.polaris.cynosure.response.servicediscovery;

import java.io.Serializable;
import java.util.List;

/**
 * 服务发现明细-响应
 *
 * @author sctang2
 * @create 2017-12-08 17:13
 **/
public class ServiceDiscoveryDetailResponseBody implements Serializable {
    private static final long serialVersionUID = 8705609676406826317L;

    //服务名称
    private String service;

    //负载均衡列表
    private List<LoadBalanceDetail> loadbalance;

    public String getService() {
        return service;
    }

    public void setService(String service) {
        this.service = service;
    }

    public List<LoadBalanceDetail> getLoadbalance() {
        return loadbalance;
    }

    public void setLoadbalance(List<LoadBalanceDetail> loadbalance) {
        this.loadbalance = loadbalance;
    }

    @Override
    public String toString() {
        return "ServiceDiscoveryDetailResponseBody{" +
                "service='" + service + '\'' +
                ", loadbalance=" + loadbalance +
                '}';
    }
}
