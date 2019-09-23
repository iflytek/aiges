package com.iflytek.ccr.polaris.cynosure.response.serviceconfig;

import com.iflytek.ccr.polaris.cynosure.domain.ServiceConfig;

import java.io.Serializable;
import java.util.List;

/**
 * 服务配置明细-响应
 *
 * @author sctang2
 * @create 2017-11-21 11:54
 **/
public class GrayConfigListDetailResponseBody implements Serializable {
    private static final long serialVersionUID = 3465537468615599850L;

    private List<ServiceConfig> serviceConfigList;

    public List<ServiceConfig> getServiceConfigList() {
        return serviceConfigList;
    }

    public void setServiceConfigList(List<ServiceConfig> serviceConfigList) {
        this.serviceConfigList = serviceConfigList;
    }

    @Override
    public String toString() {
        return "GrayConfigListDetailResponseBody{" +
                "serviceConfigList=" + serviceConfigList +
                '}';
    }
}
