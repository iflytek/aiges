package com.iflytek.ccr.finder.cache;

import com.iflytek.ccr.finder.handler.ServiceHandle;
import com.iflytek.ccr.finder.value.*;

import java.util.*;

public class ServiceCache {

    /**
     * 停用的服务实例路径列表
     */
    public Vector<String> discontinuationPathList = new Vector<String>();

    /**
     * 存放初始化的时候，已经返回的配置数据
     */
    public Set<String> initSet = new HashSet<>();

    public Map<String, MonitorValue> serviceMonitorMap = new HashMap<>();

    public Map<String, List<ServiceInstance>> instanceListMap = new HashMap<>();

    public Map<String, ServiceConfig> confCacheMap = new HashMap<>();

    public Map<String, List<ServiceRouteValue>> routeCacheMap = new HashMap<>();

    private List<SubscribeRequestValue> requestValueList;
    private ServiceHandle serviceHandle;

    public List<SubscribeRequestValue> getRequestValueList() {
        return requestValueList;
    }

    public void setRequestValueList(List<SubscribeRequestValue> requestValueList) {
        this.requestValueList = requestValueList;
    }

    public ServiceHandle getServiceHandle() {
        return serviceHandle;
    }

    public void setServiceHandle(ServiceHandle serviceHandle) {
        this.serviceHandle = serviceHandle;
    }
}
