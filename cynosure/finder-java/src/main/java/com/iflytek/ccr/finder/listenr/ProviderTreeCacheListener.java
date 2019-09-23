package com.iflytek.ccr.finder.listenr;

import com.iflytek.ccr.finder.FinderManager;
import com.iflytek.ccr.finder.constants.Constants;
import com.iflytek.ccr.finder.handler.ServiceHandle;
import com.iflytek.ccr.finder.service.RouteService;
import com.iflytek.ccr.finder.service.impl.RouteServiceImpl;
import com.iflytek.ccr.finder.utils.*;
import com.iflytek.ccr.finder.value.*;
import com.iflytek.ccr.zkutil.ZkHelper;
import org.apache.curator.framework.CuratorFramework;
import org.apache.curator.framework.recipes.cache.TreeCacheEvent;
import org.apache.curator.framework.recipes.cache.TreeCacheListener;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.List;

/**
 * 服务发现route节点监听
 */
public class ProviderTreeCacheListener implements TreeCacheListener {

    private static final Logger logger = LoggerFactory.getLogger(ProviderTreeCacheListener.class);

    private FinderManager finderManager;
    private String serviceName;
    private ServiceHandle serviceHandle;
    private String apiVersion;

    private RouteService routeService = new RouteServiceImpl();

    public ProviderTreeCacheListener(FinderManager finderManager, SubscribeRequestValue requestValue, ServiceHandle serviceHandle) {
        this.finderManager = finderManager;
        this.serviceName = requestValue.getServiceName();
        this.serviceHandle = serviceHandle;
        this.apiVersion = requestValue.getApiVersion();
    }

    @Override
    public void childEvent(CuratorFramework client, TreeCacheEvent event) throws Exception {
        logger.info("childEvent:" + event);
        if (TreeCacheEvent.Type.NODE_ADDED.equals(event.getType()) || TreeCacheEvent.Type.NODE_REMOVED.equals(event.getType()) || TreeCacheEvent.Type.NODE_UPDATED.equals(event.getType())) {
            try {
                String path = event.getData().getPath();
                String providerPath = PathUtils.getProviderPath(serviceName, apiVersion);
                //由于监听的providerPath（/polaris/service/c94370c924bac3f56e77434935613b23/iatExecutor/1.0/provider）及其子节点path，所以忽略providerPath本身
                if (path.equals(providerPath)) {
                    return;
                }
                String addr = path.substring(path.lastIndexOf("/") + 1);
                String cacheKey = serviceName + "_" + apiVersion;
                ServiceInstance cacheInstance = getCacheServiceInstance(path, cacheKey);


                switch (event.getType()) {
                    //新增服务节点
                    case NODE_ADDED:
                        //已经存在，则忽略，不进行处理
                        if (null != cacheInstance && addr.equals(cacheInstance.getAddr())) {
                            return;
                        }

                        dealInstanceChanged(event, InstanceChangedEvent.Type.ADD);
                        break;
                    //服务节点下线
                    case NODE_REMOVED:
                        dealInstanceChanged(event, InstanceChangedEvent.Type.REMVOE);
                        break;
                    //配置信息发生变化
                    case NODE_UPDATED:
                        boolean isSuccess = false;

                        //是否在禁用列表中
                        boolean isDiscontinuation = finderManager.getServiceCache().discontinuationPathList.contains(path);
//                        String cacheKey = serviceName + "_" + apiVersion;
//                        ServiceInstance cacheInstance = getCacheServiceInstance(path, cacheKey);
                        //如果不在禁用列表中，而且不在当前服务提供者列表中，则忽略
                        if (!isDiscontinuation) {
                            //不在当前服务提供者列表中，忽略
                            if (null == cacheInstance) {
                                return;
                            }
                        }

                        ServiceInstance instance = routeService.parseServiceInstanceData(finderManager, path);
                        //获取的实例不为空，则说明是可用状态
                        if (null != instance) {
                            // 不可用--》 可用
                            if (isDiscontinuation) {
                                finderManager.getServiceCache().discontinuationPathList.remove(path);
                                isSuccess = dealInstanceChanged(event, InstanceChangedEvent.Type.ADD);
                            } else {
                                //可用--》可用，说明实例配置信息可能发生变更
                                //判断是否真的有变化,如果没有真正的变化，则直接忽略
                                if (isIgnore(cacheInstance, instance)) {
                                    return;
                                }
                                cacheInstance.setJsonConfig(instance.getJsonConfig());
                                //刷新缓存文件
                                refreshCacheFile(cacheKey, finderManager.getServiceCache().instanceListMap.get(cacheKey));
                                isSuccess = serviceHandle.onServiceInstanceConfigChanged(serviceName, instance.getAddr(), instance.getJsonConfig());
                            }
                        } else {
                            //获取的实例为空，则说明是不可用状态:检查以前是否是不可用，如果是则忽略，否则进行下线处理
                            if (!isDiscontinuation) {
                                isSuccess = dealUpdateInstanceChanged(event, InstanceChangedEvent.Type.REMVOE);
                            }
                        }

                        //解析pushId
                        ZkHelper zkHelper = ZkInstanceUtil.getInstance();
                        ZkDataValue zkDataValue = ByteUtil.parseZkData(zkHelper.getByteData(path));
                        //如果pushid非法，则不进行反馈
                        if (null == zkDataValue || StringUtils.isNullOrEmpty(zkDataValue.getPushId())) {
                            break;
                        }
                        long updateTime = System.currentTimeMillis();
                        if (isSuccess) {
                            RemoteUtil.pushServiceFeedback(finderManager, zkDataValue.getPushId(), path.substring(path.lastIndexOf("/") + 1), Constants.UPDATE_STATUS_SUCCESS, Constants.LOAD_STATUS_SUCCESS, String.valueOf(updateTime), String.valueOf(System.currentTimeMillis()), apiVersion, Constants.KEY_SERVICE_INSTANCE_CHANGE);
                        } else {
                            RemoteUtil.pushServiceFeedback(finderManager, zkDataValue.getPushId(), path.substring(path.lastIndexOf("/") + 1), Constants.UPDATE_STATUS_SUCCESS, Constants.LOAD_STATUS_FAIL, String.valueOf(updateTime), String.valueOf(System.currentTimeMillis()), apiVersion, Constants.KEY_SERVICE_INSTANCE_CHANGE);
                        }

                        break;
                    case CONNECTION_SUSPENDED:
                    case CONNECTION_RECONNECTED:
                    case CONNECTION_LOST:
                    case INITIALIZED:
                }
            } catch (Exception e) {
                logger.error("", e);
            }
        }
    }

    /**
     * 判断前后是否有变化，如果没有变化，则忽略
     *
     * @param cacheInstance
     * @param instance
     * @return
     */
    private boolean isIgnore(ServiceInstance cacheInstance, ServiceInstance instance) {
        boolean isIgnore = false;
        if (null == cacheInstance.getJsonConfig() && null == instance.getJsonConfig()) {
            isIgnore = true;
        } else if (null != cacheInstance.getJsonConfig() && cacheInstance.getJsonConfig().equals(instance.getJsonConfig())) {
            isIgnore = true;
        }
        return isIgnore;
    }

    /**
     * 从缓存中获取匹配的实例
     *
     * @param path
     * @param cacheKey
     * @return
     */
    private ServiceInstance getCacheServiceInstance(String path, String cacheKey) {
        String changeAddr = path.substring(path.lastIndexOf("/") + 1);
        List<ServiceInstance> instanceList = finderManager.getServiceCache().instanceListMap.get(cacheKey);
        for (ServiceInstance instance : instanceList) {
            if (changeAddr.equals(instance.getAddr())) {
                return instance;
            }
        }
        return null;
    }


    /**
     * 处理实例变化
     *
     * @param event
     * @param type
     */
    private boolean dealInstanceChanged(TreeCacheEvent event, InstanceChangedEvent.Type type) {
        //获取路由规则
        String routePath = PathUtils.getRoutePath(serviceName, apiVersion);
        List<ServiceRouteValue> routeValueList = routeService.parseRouteData(routePath);

        ServiceInstance instance = routeService.parseServiceInstanceData(finderManager, event.getData().getPath());
        instance = routeService.filterServiceInstanceByRoute(instance, routeValueList, finderManager.getBootConfig().getMeteData().getAddress());

        //如果属于可以增加的实例，则通知订阅者
        if (null != instance) {
            //获取正在使用的提供者列表
            String cacheKey = serviceName + "_" + apiVersion;
            List<ServiceInstance> instanceList = finderManager.getServiceCache().instanceListMap.get(cacheKey);
            instanceList.add(instance);
            //刷新缓存文件
            if (finderManager.getBootConfig().isServiceCache()) {
                try {
                    Service service = (Service) FinderFileUtils.readObjectFromFile(PathUtils.getCacheFilePath(finderManager, "service") + cacheKey);
                    service.getServerList().add(instance);
                    FinderFileUtils.writeObjectToFile(PathUtils.getCacheFilePath(finderManager, "service") + cacheKey, service);
                } catch (Exception e) {
                    logger.error("", e);
                }
            }

            List<ServiceInstance> serviceInstanceList = new ArrayList<>();
            serviceInstanceList.add(instance);
            InstanceChangedEvent instanceChangedEvent = new InstanceChangedEvent(type, serviceInstanceList);
            List<InstanceChangedEvent> eventList = new ArrayList<>();
            eventList.add(instanceChangedEvent);
            return serviceHandle.onServiceInstanceChanged(serviceName, eventList);

        }
        return false;
    }

    /**
     * 处理实例变化
     *
     * @param event
     * @param type
     */
    private boolean dealUpdateInstanceChanged(TreeCacheEvent event, InstanceChangedEvent.Type type) {
        //获取当前变更实例地址
        String addr = event.getData().getPath().substring(event.getData().getPath().lastIndexOf("/") + 1);
        //获取正在使用的提供者列表
        String cacheKey = serviceName + "_" + apiVersion;
        List<ServiceInstance> instanceList = finderManager.getServiceCache().instanceListMap.get(cacheKey);
        ServiceInstance instance = null;
        for (ServiceInstance temp : instanceList) {
            if (temp.getAddr().equals(addr)) {
                instance = temp;
                break;
            }
        }

        //如果属于可以变更的实例，则通知订阅者
        if (null != instance) {
            //更新当前可用实例列表缓存
            instanceList.remove(instance);
            //刷新缓存文件
            refreshCacheFile(cacheKey, instanceList);

            List<ServiceInstance> serviceInstanceList = new ArrayList<>();
            serviceInstanceList.add(instance);
            InstanceChangedEvent instanceChangedEvent = new InstanceChangedEvent(type, serviceInstanceList);
            List<InstanceChangedEvent> eventList = new ArrayList<>();
            eventList.add(instanceChangedEvent);
            return serviceHandle.onServiceInstanceChanged(serviceName, eventList);
        }
        return false;
    }

    /**
     * 刷新缓存文件
     *
     * @param cacheKey
     * @param instanceList
     */
    private void refreshCacheFile(String cacheKey, List<ServiceInstance> instanceList) {
        if (finderManager.getBootConfig().isServiceCache()) {
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
