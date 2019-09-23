package com.iflytek.ccr.finder.listenr;

import com.iflytek.ccr.finder.FinderManager;
import com.iflytek.ccr.finder.constants.Constants;
import com.iflytek.ccr.finder.handler.ServiceHandle;
import com.iflytek.ccr.finder.service.RouteService;
import com.iflytek.ccr.finder.service.impl.RouteServiceImpl;
import com.iflytek.ccr.finder.utils.*;
import com.iflytek.ccr.finder.value.*;
import com.iflytek.ccr.zkutil.ZkHelper;
import org.apache.curator.framework.recipes.cache.NodeCacheListener;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.List;

/**
 * 服务发现route节点监听
 */
public class RouteNodeCacheListener implements NodeCacheListener {

    private static final Logger logger = LoggerFactory.getLogger(RouteNodeCacheListener.class);

    private FinderManager finderManager;
    private String serviceName;
    private ServiceHandle serviceHandle;
    private String apiVersion;

    private RouteService routeService = new RouteServiceImpl();

    public RouteNodeCacheListener(FinderManager finderManager, String serviceName, ServiceHandle serviceHandle, String apiVersion) {
        this.finderManager = finderManager;
        this.serviceName = serviceName;
        this.serviceHandle = serviceHandle;
        this.apiVersion = apiVersion;
    }

    @Override
    public void nodeChanged() throws Exception {
        try {
            String routePath = PathUtils.getRoutePath(serviceName, apiVersion);
            //获取最新的路由规则信息
            List<ServiceRouteValue> routeValueList = routeService.parseRouteData(routePath);
            String cacheKey = serviceName + "_" + apiVersion;
            //获取缓存的路由信息
            List<ServiceRouteValue> routeCacheList = finderManager.getServiceCache().routeCacheMap.get(cacheKey);
            if (null == routeCacheList) {
                finderManager.getServiceCache().routeCacheMap.put(cacheKey, routeValueList);
                return;
            } else {
                //比较最新的路由规则与缓存的路由规则是否一致，如果一致，则忽略，不进行处理
                if (isSame(routeCacheList, routeValueList)) {
                    return;
                }else {
                    //更新缓存
                    finderManager.getServiceCache().routeCacheMap.put(cacheKey, routeValueList);
                }
            }

            //获取正在使用的提供者列表
            List<ServiceInstance> instanceList = finderManager.getServiceCache().instanceListMap.get(serviceName + "_" + apiVersion);

            //获取当前的提供者实例
            String providerPath = PathUtils.getProviderPath(serviceName, apiVersion);
            List<ServiceInstance> newInstanceList = routeService.parseServiceInstanceByPath(finderManager, providerPath);
            newInstanceList = routeService.filterServiceInstanceByRoute(newInstanceList, routeValueList, finderManager.getBootConfig().getMeteData().getAddress());

            //获取变化的实例
            List<InstanceChangedEvent> notifyList = routeService.compareServiceInstanceList(instanceList, newInstanceList);


            //通知订阅者
            if (null != notifyList && !notifyList.isEmpty()) {
                updateInstanceListMap(apiVersion, notifyList);
                boolean isSuccess = serviceHandle.onServiceInstanceChanged(serviceName, notifyList);

                //解析pushId
                ZkHelper zkHelper = ZkInstanceUtil.getInstance();
                ZkDataValue zkDataValue = ByteUtil.parseZkData(zkHelper.getByteData(routePath));
                //如果pushid非法，则不进行反馈
                if (null == zkDataValue || StringUtils.isNullOrEmpty(zkDataValue.getPushId())) {
                    return;
                }
                long updateTime = System.currentTimeMillis();
                if (isSuccess) {
                    RemoteUtil.pushServiceFeedback(finderManager, zkDataValue.getPushId(), "", Constants.UPDATE_STATUS_SUCCESS, Constants.LOAD_STATUS_SUCCESS, String.valueOf(updateTime), String.valueOf(System.currentTimeMillis()), apiVersion, Constants.KEY_SERVICE_ROUTE_CHANGE);
                } else {
                    RemoteUtil.pushServiceFeedback(finderManager, zkDataValue.getPushId(), "", Constants.UPDATE_STATUS_SUCCESS, Constants.LOAD_STATUS_FAIL, String.valueOf(updateTime), String.valueOf(System.currentTimeMillis()), apiVersion, Constants.KEY_SERVICE_ROUTE_CHANGE);
                }
            }
        } catch (Exception e) {
            logger.error(String.format("RouteNodeCacheListener error:%s", e.getMessage()), e);
        }
    }

    /**
     * 判断两个集合中的路由规则是否一致
     *
     * @param routeCacheList
     * @param routeValueList
     * @return
     */
    private boolean isSame(List<ServiceRouteValue> routeCacheList, List<ServiceRouteValue> routeValueList) {
        boolean isSame = true;
        if (routeCacheList.size() != routeValueList.size()) {
            isSame = false;
        } else {
            for (ServiceRouteValue value : routeCacheList) {
                boolean find = false;
                for(ServiceRouteValue temp : routeValueList){
                    if( value.toString().equals(temp.toString())){
                        find = true;
                        break;
                    }
                }
                if(!find){
                    isSame = false;
                    break;
                }
            }
        }
        return isSame;
    }


    /**
     * 更新当前缓存的实例列表
     *
     * @param apiVersion
     * @param notifyList
     */
    private void updateInstanceListMap(String apiVersion, List<InstanceChangedEvent> notifyList) {
        //获取正在使用的提供者列表
        List<ServiceInstance> instanceList = finderManager.getServiceCache().instanceListMap.get(serviceName + "_" + apiVersion);
        for (InstanceChangedEvent event : notifyList) {
            if (InstanceChangedEvent.Type.REMVOE.equals(event.getType())) {
                List<ServiceInstance> removeList = new ArrayList<>();
                for (ServiceInstance instance : event.getServiceInstanceList()) {

                    for (ServiceInstance temp : instanceList) {
                        if (instance.getAddr().equals(temp.getAddr())) {
                            removeList.add(temp);
                            break;
                        }
                    }
                }
                if (!removeList.isEmpty()) {
                    instanceList.removeAll(removeList);
                }
            } else {
                instanceList.addAll(event.getServiceInstanceList());
            }

        }

        //刷新缓存文件
        if (finderManager.getBootConfig().isServiceCache()) {
            String cacheKey = serviceName + "_" + apiVersion;
            try {
                Service service = (Service) FinderFileUtils.readObjectFromFile(PathUtils.getCacheFilePath(finderManager, "service") + cacheKey);
                service.setServerList(instanceList);
                FinderFileUtils.writeObjectToFile(PathUtils.getCacheFilePath(finderManager, "service") + cacheKey, service);
            } catch (Exception e) {
                logger.error("", e);
            }
        }
    }
}
