package com.iflytek.ccr.finder.service;

import com.iflytek.ccr.finder.FinderManager;
import com.iflytek.ccr.finder.value.InstanceChangedEvent;
import com.iflytek.ccr.finder.value.Service;
import com.iflytek.ccr.finder.value.ServiceInstance;
import com.iflytek.ccr.finder.value.ServiceRouteValue;

import java.util.List;

/**
 * 路由服务
 */
public interface RouteService {

    /**
     * 解析路由配置数据
     *
     * @param routePath
     * @return
     */
    List<ServiceRouteValue> parseRouteData(String routePath);


    /**
     * 解析conf配置数据
     *
     * @param confPath
     * @return
     */
    String parseConfData(String confPath);

    Service parseServiceData(FinderManager finderManager, String basePath, String serviceName, String apiVersion);

    /**
     * 解析服务实例数据
     *
     * @param path
     * @return
     */
    ServiceInstance parseServiceInstanceData(FinderManager finderManager, String path);

    /**
     * 按照路径解析服务实例列表
     *
     * @param finderManager
     * @param path
     * @return
     */
    List<ServiceInstance> parseServiceInstanceByPath(FinderManager finderManager, String path);


    /**
     * 按照路由来过滤服务实例列表
     *
     * @param list
     * @param routeValueList
     * @param addr
     * @return
     */
    List<ServiceInstance>  filterServiceInstanceByRoute(List<ServiceInstance> list, List<ServiceRouteValue> routeValueList, String addr);

    /**
     * 按照路由来过滤服务实例列表
     *
     * @param instance
     * @param routeValueList
     * @param addr
     * @return
     */
    ServiceInstance filterServiceInstanceByRoute(ServiceInstance instance, List<ServiceRouteValue> routeValueList, String addr);


    /**
     * 对比实例变化
     *
     * @param before
     * @param after
     * @return
     */
    List<InstanceChangedEvent> compareServiceInstanceList(List<ServiceInstance> before, List<ServiceInstance> after);

}
