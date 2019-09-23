package com.iflytek.ccr.finder.service.impl;

import com.iflytek.ccr.finder.FinderManager;
import com.iflytek.ccr.finder.constants.Constants;
import com.iflytek.ccr.finder.handler.ServiceHandle;
import com.iflytek.ccr.finder.listenr.ConfNodeCacheListener;
import com.iflytek.ccr.finder.listenr.ProviderTreeCacheListener;
import com.iflytek.ccr.finder.listenr.RouteNodeCacheListener;
import com.iflytek.ccr.finder.service.CommonService;
import com.iflytek.ccr.finder.service.RouteService;
import com.iflytek.ccr.finder.service.ServiceFinder;
import com.iflytek.ccr.finder.utils.*;
import com.iflytek.ccr.finder.value.*;
import com.iflytek.ccr.zkutil.ZkHelper;
import org.apache.curator.framework.recipes.cache.NodeCache;
import org.apache.curator.framework.recipes.cache.TreeCache;
import org.codehaus.jackson.JsonNode;
import org.codehaus.jackson.map.ObjectMapper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

/**
 * 服务发现实现类
 */
public class ServiceFinderImpl implements ServiceFinder {

    private Logger logger = LoggerFactory.getLogger(ServiceFinderImpl.class);

    CommonService commonService = new CommonServiceImpl();
    private RouteService routeService = new RouteServiceImpl();

    @Override
    public CommonResult registerService(FinderManager finderManager, String addr, String apiVersion) {
        String serviceName = finderManager.getBootConfig().getMeteData().getService();
        logger.info(String.format("name:%s,addr:%s", serviceName, addr));
        CommonResult commonResult = new CommonResult();
        try {
            //注册服务信息
            ZkHelper zkHelper = ZkInstanceUtil.getInstance();
            String path = PathUtils.getProviderPath(serviceName, apiVersion) + "/" + addr;
            logger.info(String.format("ProviderPath:%s", path));

            //如果存在节点，则要先删除，后增加
            if (zkHelper.checkExists(path)) {
                zkHelper.remove(path);
            }
            zkHelper.addEphemeral(path, "");
            createBaseNode(finderManager, apiVersion);

            commonResult.setRet(ErrorCode.SUCCESS);
            logger.info(commonResult.toString());

            //添加到监控列表
            finderManager.getGlobalCache().monitorPathList.add(path);

            //同步注册服务信息到网站上
            String result = RemoteUtil.registerServiceInfo(finderManager.getBootConfig(), apiVersion);
            logger.info(String.format("registerServiceInfo result:%s", result));
        } catch (Exception e) {
            logger.error(String.format("registerService error:%s", e.getMessage()), e);
            commonResult.setRet(ErrorCode.INTERNAL_EXCEPTION);
            commonResult.setMsg(e.getMessage());
        }
        return commonResult;
    }

    /**
     * 创建基础路径
     *
     * @param finderManager
     * @param apiVersion
     */
    private void createBaseNode(FinderManager finderManager, String apiVersion) {
        ZkHelper zkHelper = ZkInstanceUtil.getInstance();
        String serviceName = finderManager.getBootConfig().getMeteData().getService();
        String confPath = PathUtils.getConfPath(serviceName, apiVersion);
        if (!zkHelper.checkExists(confPath)) {
            zkHelper.addPersistent(confPath, "");
        }
        String routePath = PathUtils.getRoutePath(serviceName, apiVersion);
        if (!zkHelper.checkExists(routePath)) {
            zkHelper.addPersistent(routePath, "");
        }


    }

    @Override
    public CommonResult unRegisterService(FinderManager finderManager, String apiVersion) {
        return unRegisterService(finderManager, finderManager.getBootConfig().getMeteData().getAddress(), apiVersion);
    }

    @Override
    public CommonResult unRegisterService(FinderManager finderManager, String addr, String apiVersion) {
        String serviceName = finderManager.getBootConfig().getMeteData().getService();
        logger.info(String.format("name:%s,addr:%s", serviceName, addr));
        CommonResult commonResult = new CommonResult();
        try {
            ZkHelper zkHelper = ZkInstanceUtil.getInstance();
            String path = PathUtils.getProviderPath(serviceName, apiVersion) + "/" + addr;
            logger.info(String.format("unRegisterService path:%s", path));
            if (zkHelper.checkExists(path) || finderManager.getGlobalCache().monitorPathList.contains(path)) {
                //从监控列表中移除
                finderManager.getGlobalCache().monitorPathList.remove(path);
                //从zk中移除
                zkHelper.remove(path);
                commonResult.setRet(ErrorCode.SUCCESS);
            } else {
                commonResult.setRet(ErrorCode.PARAM_INVALID);
                commonResult.setMsg("path does not exists:" + path);
            }
        } catch (Exception e) {
            logger.error(String.format("unRegisterService error:%s", e.getMessage()), e);
            commonResult.setRet(ErrorCode.INTERNAL_EXCEPTION);
            commonResult.setMsg(e.getMessage());
        }
        return commonResult;
    }

    @Override
    public CommonResult<List<Service>> useAndSubscribeService(FinderManager finderManager, List<SubscribeRequestValue> requestValueList, ServiceHandle serviceHandle) {
        CommonResult commonResult = new CommonResult();
        try {
            //保存请求参数到缓存
            finderManager.getServiceCache().setRequestValueList(requestValueList);
            finderManager.getServiceCache().setServiceHandle(serviceHandle);

            List<Service> serviceList = new ArrayList<>();
            for (SubscribeRequestValue requestValue : requestValueList) {

                String confPath = PathUtils.getConfPath(requestValue.getServiceName(), requestValue.getApiVersion());
                String routePath = PathUtils.getRoutePath(requestValue.getServiceName(), requestValue.getApiVersion());
                String providerPath = PathUtils.getProviderPath(requestValue.getServiceName(), requestValue.getApiVersion());

                //初始化服务信息
                Service service = getService(finderManager, requestValue, providerPath, confPath, routePath);
                serviceList.add(service);

                //添加对conf节点的监听
                ConfNodeCacheListener confNodeCacheListener = new ConfNodeCacheListener(finderManager, requestValue.getServiceName(), serviceHandle, requestValue.getApiVersion());
                ZkHelper zkHelper = ZkInstanceUtil.getInstance();
                NodeCache confNodeCache = zkHelper.addListener(confNodeCacheListener, confPath, false);

                //添加对route节点的监听
                RouteNodeCacheListener routeNodeCacheListener = new RouteNodeCacheListener(finderManager, requestValue.getServiceName(), serviceHandle, requestValue.getApiVersion());
                NodeCache routeNodeCache = zkHelper.addListener(routeNodeCacheListener, routePath, false);

                //监听provider节点
                ProviderTreeCacheListener providerListener = new ProviderTreeCacheListener(finderManager, requestValue, serviceHandle);
                TreeCache providerTreeCache = zkHelper.addListener(providerListener, providerPath);

                String consumerPath = PathUtils.getServiceConsumerPath(requestValue, finderManager.getBootConfig().getMeteData().getAddress());
                commonService.registerConsumer(consumerPath);
                //增加对消费者路径的监控
                finderManager.getGlobalCache().monitorPathList.add(consumerPath);

                //保存节点监控信息
                MonitorValue monitorValue = new MonitorValue();
                monitorValue.setConfPath(confPath);
                monitorValue.setConfNodeCache(confNodeCache);
                monitorValue.setRoutePath(routePath);
                monitorValue.setRouteNodeCache(routeNodeCache);
                monitorValue.setProviderPath(providerPath);
                monitorValue.setProviderTreeCache(providerTreeCache);
                //保存到缓存中
                finderManager.getServiceCache().serviceMonitorMap.put(requestValue.getCacheKey(), monitorValue);

            }


            commonResult.setRet(ErrorCode.SUCCESS);
            commonResult.setData(serviceList);
        } catch (Exception e) {
            logger.error(String.format("useAndSubscribeService error:%s", e.getMessage()), e);
            commonResult.setRet(ErrorCode.INTERNAL_EXCEPTION);
            commonResult.setMsg(e.getMessage());
        }

        return commonResult;
    }

    /**
     * 获取service 对象
     *
     * @param finderManager
     * @param requestValue
     * @param providerPath
     * @param confPath
     * @param routePath
     * @return
     */
    private Service getService(FinderManager finderManager, SubscribeRequestValue requestValue, String providerPath, String confPath, String routePath) {
        List<ServiceInstance> instanceList = routeService.parseServiceInstanceByPath(finderManager, providerPath);
        //初始化init时加载的实例列表
//        for (ServiceInstance instance : instanceList) {
//            finderManager.getServiceCache().initSet.add(instance.getAddr());
//        }
        //获取当前service信息
        Service service = new Service();
        service.setName(requestValue.getServiceName());
        service.setApiVersion(requestValue.getApiVersion());
        ServiceConfig config = new ServiceConfig();
        config.setJsonConfig(routeService.parseConfData(confPath));
        service.setServiceConfig(config);
        service.setServerList(routeService.filterServiceInstanceByRoute(instanceList, routeService.parseRouteData(routePath), finderManager.getBootConfig().getMeteData().getAddress()));

        //更新缓存
        finderManager.getServiceCache().instanceListMap.put(requestValue.getCacheKey(), service.getServerList());
        finderManager.getServiceCache().confCacheMap.put(requestValue.getCacheKey(), service.getServiceConfig());
        FinderFileUtils.writeObjectToFile(PathUtils.getCacheFilePath(finderManager, "service") + requestValue.getCacheKey(), service);

        return service;
    }


    /**
     * 获取服务配置
     *
     * @param finderManager
     * @param serviceName
     * @return
     */
    public CommonResult<ServiceConfig> getServiceConfig(FinderManager finderManager, String serviceName, String apiVersion) {
        CommonResult<ServiceConfig> configCommonResult = new CommonResult<>();
        ServiceConfig serviceConfig = new ServiceConfig();
        configCommonResult.setData(serviceConfig);
        //获取conf路径
        String path = PathUtils.getConfPath(serviceName, apiVersion);
        ZkHelper zkHelper = ZkInstanceUtil.getInstance();

        boolean isExists = zkHelper.checkExists(path);
        if (isExists) {
            byte[] data = zkHelper.getByteData(path);
            ZkDataValue zkDataValue = ByteUtil.parseZkData(data);
            if (ErrorCode.SUCCESS == zkDataValue.getRet()) {
                configCommonResult.setMsg(zkDataValue.getPushId());
                ObjectMapper mapper = new ObjectMapper();
                try {
                    JsonNode rootNode = mapper.readTree(new String(zkDataValue.getRealData(), Constants.DEFAULT_CHARSET));
                    JsonNode jsonConfig = rootNode.path("data");//.path("consumer");
                    if (null != jsonConfig) {
                        serviceConfig.setJsonConfig(jsonConfig.getTextValue());
                    }
                } catch (IOException e) {
                    logger.error("", e);
                }
            } else {
                logger.error(String.format("data parse error:%s", zkDataValue.toString()));
            }
        } else {
            logger.error(String.format("path:%s does not exists", path));
        }
        return configCommonResult;
    }


    @Override
    public CommonResult unSubscribeService(FinderManager finderManager, SubscribeRequestValue requestValue) {
        CommonResult commonResult = new CommonResult();
        try {
            MonitorValue monitorValue = finderManager.getServiceCache().serviceMonitorMap.get(requestValue.getCacheKey());
            commonResult.setRet(ErrorCode.SUCCESS);
            String consumerPath = PathUtils.getServiceConsumerPath(requestValue, finderManager.getBootConfig().getMeteData().getAddress());
            commonService.unRegisterConsumer(consumerPath);
            finderManager.getGlobalCache().monitorPathList.remove(consumerPath);
            if (null != monitorValue) {
                try {
                    monitorValue.getConfNodeCache().close();
                } catch (IOException e) {
                    logger.error("", e);
                }
                try {
                    monitorValue.getRouteNodeCache().close();
                } catch (IOException e) {
                    logger.error("", e);
                }
                monitorValue.getProviderTreeCache().close();
            } else {
                commonResult.setRet(ErrorCode.PARAM_INVALID);
                commonResult.setMsg("Can not find Subscribe infomation");
            }
        } catch (Exception e) {
            logger.error(String.format("unSubscribeService error:%s", e.getMessage()), e);
            commonResult.setRet(ErrorCode.INTERNAL_EXCEPTION);
            commonResult.setMsg(e.getMessage());
        }

        return commonResult;
    }

    @Override
    public CommonResult unSubscribeService(FinderManager finderManager, List<SubscribeRequestValue> requestValueList) {
        for (SubscribeRequestValue value : requestValueList) {
            CommonResult commonResult = unSubscribeService(finderManager, value);
            if (ErrorCode.SUCCESS != commonResult.getRet()) {
                return commonResult;
            }
        }
        CommonResult commonResult = new CommonResult();
        commonResult.setRet(ErrorCode.SUCCESS);
        return commonResult;
    }
}
