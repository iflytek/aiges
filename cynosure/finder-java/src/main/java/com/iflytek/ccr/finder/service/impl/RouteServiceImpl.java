package com.iflytek.ccr.finder.service.impl;

import com.iflytek.ccr.finder.FinderManager;
import com.iflytek.ccr.finder.constants.Constants;
import com.iflytek.ccr.finder.service.RouteService;
import com.iflytek.ccr.finder.utils.ByteUtil;
import com.iflytek.ccr.finder.utils.JacksonUtils;
import com.iflytek.ccr.finder.utils.ZkInstanceUtil;
import com.iflytek.ccr.finder.value.*;
import com.iflytek.ccr.zkutil.ZkHelper;
import org.codehaus.jackson.JsonNode;
import org.codehaus.jackson.map.ObjectMapper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

public class RouteServiceImpl implements RouteService {

    /**
     * logger
     */
    private static final Logger logger = LoggerFactory.getLogger(RouteServiceImpl.class);

    @Override
    public List<ServiceRouteValue> parseRouteData(String routePath) {
        List<ServiceRouteValue> routeValueList = new ArrayList<>();
        try {
            ZkHelper zkHelper = ZkInstanceUtil.getInstance();
            byte[] data = zkHelper.getByteData(routePath);
            ZkDataValue zkDataValue = ByteUtil.parseZkData(data);
            if (ErrorCode.SUCCESS == zkDataValue.getRet()) {
                String routeJson = new String(zkDataValue.getRealData(), Constants.DEFAULT_CHARSET);
                List grayList = JacksonUtils.toObject(routeJson, List.class);

                for (Object v : grayList) {
                    Map<String, Object> vMap = (Map<String, Object>) v;
                    ServiceRouteValue routeValue = new ServiceRouteValue();
                    routeValue.setConsumer((ArrayList) vMap.get(Constants.KEY_ROUTE_CONSUMER));
                    routeValue.setProvider((ArrayList) vMap.get(Constants.KEY_ROUTE_PROVIDER));
                    routeValue.setOnly(String.valueOf(vMap.get(Constants.KEY_ROUTE_ONLY)));
                    routeValue.setId(String.valueOf(vMap.get(Constants.KEY_ROUTE_ROUTE_RULE_ID)));
                    routeValueList.add(routeValue);
                }
            }
        } catch (Exception e) {
            logger.error(String.format("parseRouteData error:%s", e.getMessage()), e);
        }
        return routeValueList;
    }

    @Override
    public String parseConfData(String confPath) {
        try {
            ZkHelper zkHelper = ZkInstanceUtil.getInstance();
            byte[] data = zkHelper.getByteData(confPath);
            ZkDataValue zkDataValue = ByteUtil.parseZkData(data);
            if (ErrorCode.SUCCESS == zkDataValue.getRet()) {
                String json = new String(zkDataValue.getRealData(), Constants.DEFAULT_CHARSET);
                return json;
            } else {
                logger.error(String.format("parseZkData error,confPath:%s", confPath));
            }
        } catch (Exception e) {
            logger.error(String.format("ConfNodeCacheListener error:%s", e.getMessage()), e);
        }
        return null;
    }

    @Override
    public Service parseServiceData(FinderManager finderManager, String basePath, String serviceName, String apiVersion) {
        Service service = new Service();
        service.setApiVersion(apiVersion);
        ServiceConfig config = new ServiceConfig();
        config.setJsonConfig(parseConfData(basePath + "/conf"));
        service.setServiceConfig(config);
        service.setName(serviceName);
        service.setServerList(parseServiceInstanceByPath(finderManager, basePath + "/provider"));
        return service;
    }

    @Override
    public ServiceInstance parseServiceInstanceData(FinderManager finderManager, String path) {
        ZkHelper zkHelper = ZkInstanceUtil.getInstance();
        byte[] data = zkHelper.getByteData(path);
        ServiceInstance serviceInstance = new ServiceInstance();
        serviceInstance.setAddr(path.substring(path.lastIndexOf("/") + 1));
        //如果没有配置实例级配置信息，则直接返回实例
        if (null == data || data.length == 0) {
            return serviceInstance;
        } else {
            ZkDataValue zkDataValue = ByteUtil.parseZkData(data);
            if (ErrorCode.SUCCESS == zkDataValue.getRet()) {
                ObjectMapper mapper = new ObjectMapper();
                try {
                    JsonNode rootNode = mapper.readTree(new String(zkDataValue.getRealData(), Constants.DEFAULT_CHARSET));
                    //解析服务订阅者需要的信息
                    JsonNode jsonConfig = rootNode.path("user");
                    if (null != jsonConfig) {
                        serviceInstance.setJsonConfig(jsonConfig.toString());
                    }
                    //解析sdk使用的实例数据
                    ServiceInstanceConfig config = parseSdkConfig(serviceInstance, rootNode);
                    if (config.isValid()) {
                        return serviceInstance;
                    } else {
                        if (!finderManager.getServiceCache().discontinuationPathList.contains(path)) {
                            finderManager.getServiceCache().discontinuationPathList.add(path);
                        }
                        return null;
                    }
                } catch (IOException e) {
                    logger.error("", e);
                }
            }
        }

        return null;
    }

    @Override
    public List<ServiceInstance> parseServiceInstanceByPath(FinderManager finderManager, String path) {
        ZkHelper zkHelper = ZkInstanceUtil.getInstance();
        List<String> childList = zkHelper.getChildren(path);
        List<ServiceInstance> list = new ArrayList<>();
        if (childList != null && !childList.isEmpty()) {
            for (String child : childList) {
                ServiceInstance instance = parseServiceInstanceData(finderManager, path + "/" + child);
                if (null != instance) {
                    list.add(instance);
                }
            }
        }
        return list;
    }

    @Override
    public ServiceInstance filterServiceInstanceByRoute(ServiceInstance instance, List<ServiceRouteValue> routeValueList, String addr) {
        if (null != instance) {
            for (ServiceRouteValue routeValue : routeValueList) {
                //消费者实例在规则中，则对比规则中的provider与当前提供者是否匹配，如果匹配，则返回，否则不返回
                if (routeValue.getConsumer().contains(addr)) {
                    List<String> providerList = routeValue.getProvider();
                    if (providerList.contains(instance.getAddr())) {
                        return instance;
                    } else if (providerList.contains(instance.getAddr()) && Constants.KEY_ROUTE_ONLY_Y.equals(routeValue.getOnly())) {
                        return null;
                    }
                    //消费者实例不在规则中，但是提供者实例在规则中，而且是only=Y，则忽略该提供者
                } else if (Constants.KEY_ROUTE_ONLY_Y.equals(routeValue.getOnly())) {
                    List<String> providerList = routeValue.getProvider();
                    if (providerList.contains(instance.getAddr())) {
                        return null;
                    }
                }
            }
        }
        return instance;
    }

    @Override
    public List<ServiceInstance> filterServiceInstanceByRoute(List<ServiceInstance> list, List<ServiceRouteValue> routeValueList, String addr) {
        List<ServiceInstance> resultList = new ArrayList<>();
        List<String> removeList = new ArrayList<>();
        for (ServiceRouteValue routeValue : routeValueList) {
            //如果规则匹配，则直接返回该条路由规则对应的提供者
            if (routeValue.getConsumer().contains(addr)) {
                List<String> providerList = routeValue.getProvider();
                for (ServiceInstance instance : list) {
                    if (providerList.contains(instance.getAddr())) {
                        resultList.add(instance);
                    }
                }
                return resultList;
            } else if (Constants.KEY_ROUTE_ONLY_Y.equals(routeValue.getOnly())) {
                removeList.addAll(routeValue.getProvider());
            }
        }

        for (ServiceInstance instance : list) {
            if (!removeList.contains(instance.getAddr())) {
                resultList.add(instance);
            }
        }
        return resultList;
    }


    @Override
    public List<InstanceChangedEvent> compareServiceInstanceList(List<ServiceInstance> before, List<ServiceInstance> after) {
        List<InstanceChangedEvent> changedEvents = new ArrayList<>();
        //之前有，现在没有：则为下线
        List<ServiceInstance> removeInstanceList = mergeServiceInstanceList(after, before);
        if (null != removeInstanceList && !removeInstanceList.isEmpty()) {
            InstanceChangedEvent removeEvent = new InstanceChangedEvent(InstanceChangedEvent.Type.REMVOE, removeInstanceList);
            changedEvents.add(removeEvent);
        }
        List<ServiceInstance> addInstanceList = mergeServiceInstanceList(before, after);
        if (null != addInstanceList && !addInstanceList.isEmpty()) {
            InstanceChangedEvent addEvent = new InstanceChangedEvent(InstanceChangedEvent.Type.ADD, addInstanceList);
            changedEvents.add(addEvent);
        }
        return changedEvents;
    }

    /**
     * 求集合A、B的差集：即B中有，A有无
     *
     * @param a
     * @param b
     * @return
     */
    private List<ServiceInstance> mergeServiceInstanceList(List<ServiceInstance> a, List<ServiceInstance> b) {

        if (null == b || b.isEmpty()) {
            return a;
        }
        if (null == a || a.isEmpty()) {
            return a;
        }
        List<ServiceInstance> addInstanceList = new ArrayList<>();
        for (ServiceInstance instance : b) {
            boolean isAdd = true;
            for (ServiceInstance temp : a) {
                if (instance.getAddr().equals(temp.getAddr())) {
                    isAdd = false;
                    break;
                }
            }
            if (isAdd) {
                addInstanceList.add(instance);
            }
        }

        return addInstanceList;
    }


    /**
     * 解析sdk需要的服务实例配置
     *
     * @param serviceInstance
     * @paroam rootNode
     */
    private ServiceInstanceConfig parseSdkConfig(ServiceInstance serviceInstance, JsonNode rootNode) {
        try {
            ServiceInstanceConfig config = new ServiceInstanceConfig();
            //设置默认值为true
            config.setValid(true);
            //解析sdk需要的信息
            JsonNode sdkJsonConfig = rootNode.path("sdk").path("is_valid");
            if (null != sdkJsonConfig) {
                if ("N".equals(sdkJsonConfig.getTextValue()) || "false".equals(sdkJsonConfig.toString())) {
                    config.setValid(false);
                }
            }
            return config;
        } catch (Exception e) {
            logger.error("parseSdkConfig error", e);
        }
        return null;
    }


}
