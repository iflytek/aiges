package com.iflytek.ccr.finder;

import com.iflytek.ccr.finder.cache.CacheUtils;
import com.iflytek.ccr.finder.cache.GlobalCache;
import com.iflytek.ccr.finder.cache.ServiceCache;
import com.iflytek.ccr.finder.constants.Constants;
import com.iflytek.ccr.finder.constants.TaskType;
import com.iflytek.ccr.finder.handler.ConfigChangedHandler;
import com.iflytek.ccr.finder.handler.ServiceHandle;
import com.iflytek.ccr.finder.monitor.MonitorTask;
import com.iflytek.ccr.finder.monitor.PathMonitorTask;
import com.iflytek.ccr.finder.service.ConfigFinder;
import com.iflytek.ccr.finder.service.ServiceFinder;
import com.iflytek.ccr.finder.service.impl.ConfigFinderImpl;
import com.iflytek.ccr.finder.service.impl.ServiceFinderImpl;
import com.iflytek.ccr.finder.utils.*;
import com.iflytek.ccr.finder.value.*;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.List;
import java.util.Map;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

/**
 * 服务发现、配置中心管理对象
 */
public class FinderManager {

    private final Logger logger = LoggerFactory.getLogger(FinderManager.class);

    /**
     * 初始化是否成功
     */
    private boolean initSuccess = false;

    private ConfigFinder configFinder = new ConfigFinderImpl();

    private ServiceFinder serviceFinder = new ServiceFinderImpl();

    public BootConfig getBootConfig() {
        return bootConfig;
    }

    private BootConfig bootConfig;

    private GlobalCache globalCache = new GlobalCache();

    private ServiceCache serviceCache = new ServiceCache();

    public GlobalCache getGlobalCache() {
        return globalCache;
    }

    public ServiceCache getServiceCache() {
        return serviceCache;
    }

    public FinderManager() {
        //启动sdk监控进程(服务订阅等失败后重试)
        Executors.newCachedThreadPool().submit(new MonitorTask(this));
    }

    public ConfigFinder getConfigFinder() {
        return configFinder;
    }

    /**
     * 初始化管理对象,如果初始化异常，可以设置读取缓存
     *
     * @param bootConfig
     */
    public CommonResult init(final BootConfig bootConfig) {
        this.bootConfig = bootConfig;
        //如果没有设置缓存路径，则设置默认值
        if (StringUtils.isNullOrEmpty(bootConfig.getCachePath())) {
            bootConfig.setCachePath(Constants.DEFAULT_CACHEP_PATH);
        }
        FinderFileUtils.createFolder(bootConfig.getCachePath());

        //获取zookeeper集群信息，已经zookeeper上的存储路径等
        CommonResult response = RemoteUtil.queryZkInfo(bootConfig);
        if (null != response && Constants.SUCCESS == response.getRet()) {
            Map map = (Map) response.getData();
            String configPath = (String) map.get(Constants.CONFIG_PATH);
            String servicePath = (String) map.get(Constants.SERVICE_PATH);
            String zkNodePath = (String) map.get(Constants.KEY_ZK_NODE_PATH);
            List zkAddrList = (List) map.get(Constants.ZK_ADDR);
            ConfigManager.getInstance().put(Constants.CONFIG_PATH, configPath);
            ConfigManager.getInstance().put(Constants.SERVICE_PATH, servicePath);
            ConfigManager.getInstance().put(Constants.KEY_ZK_NODE_PATH, zkNodePath);
            ConfigManager.getInstance().put(Constants.ZK_ADDR, StringUtils.join(zkAddrList, ","));
            ZkInstanceUtil.setZkAddr(StringUtils.join(zkAddrList, ","));
            if (bootConfig.getZkSessionTimeout() > 0 && bootConfig.getZkConnectTimeout() > 0) {
                ZkInstanceUtil.init(StringUtils.join(zkAddrList, ","), bootConfig.getZkConnectTimeout(), bootConfig.getZkSessionTimeout());
            } else {
                ZkInstanceUtil.init(StringUtils.join(zkAddrList, ","));
            }
            //启动监控线程
            ExecutorService service = Executors.newFixedThreadPool(1);
            service.submit(new PathMonitorTask(this));
        } else {
            initSuccess = false;
            this.getGlobalCache().taskQueue.add(TaskType.INIT);
        }
        return response;
    }

    /**
     * 获取当前配置信息，并订阅
     *
     * @param configNameList
     * @param configChangedHandler
     * @return
     */
    public CommonResult<List<Config>> useAndSubscribeConfig(List<String> configNameList, ConfigChangedHandler configChangedHandler) {
        //订阅配置，并返回当前的配置信息
        CommonResult<List<Config>> commonResult = null;

        if (StringUtils.isEmpty(configNameList)) {
            commonResult = new CommonResult();
            commonResult.setRet(ErrorCode.PARAM_INVALID);
            commonResult.setMsg("The list of files that are subscribed can not be empty");
            return commonResult;
        }
        if (configChangedHandler == null) {
            commonResult = new CommonResult();
            commonResult.setRet(ErrorCode.PARAM_INVALID);
            commonResult.setMsg("ConfigChangedHandler can not be null");
            return commonResult;
        }
        try {
            commonResult = configFinder.useAndSubscribeConfig(this, configNameList, configChangedHandler);
            if (ErrorCode.SUCCESS != commonResult.getRet() && bootConfig.isConfigCache()) {
                String cacheBasePath = PathUtils.getCacheFilePath(this, "config");
                commonResult = CacheUtils.getCacheConfigResult(configNameList, cacheBasePath);
                if (initSuccess) {
                    this.getGlobalCache().taskQueue.add(TaskType.CONFIG);
                }
            }
        } catch (Exception e) {
            logger.error("", e);
        }

        return commonResult;
    }

    /**
     * 取消订阅
     *
     * @param configName
     * @return
     */
    public CommonResult unSubscribeConfig(String configName) {
        return configFinder.unSubscribeConfig(this, configName);
    }

    /**
     * 注册服务
     *
     * @param apiVersion api版本号
     * @param addr       服务实例地址（xxx:xx,eg:1.2.3.4:8080）
     * @return
     */
    public CommonResult registerService(String addr, String apiVersion) {
        return serviceFinder.registerService(this, addr, apiVersion);
    }

    /**
     * 注册服务(addr       服务实例地址（xxx:xx,eg:1.2.3.4:8080）使用初始化数据里面的ServiceMeteData对象中的地址)
     *
     * @param apiVersion api版本号
     * @return
     */
    public CommonResult registerService(String apiVersion) {
        return serviceFinder.registerService(this, this.bootConfig.getMeteData().getAddress(), apiVersion);
    }

    /**
     * 注销服务
     *
     * @return
     */
    public CommonResult unRegisterService(String apiVersion) {
        return serviceFinder.unRegisterService(this, this.bootConfig.getMeteData().getAddress(), apiVersion);
    }

    /**
     * 注销服务
     *
     * @return
     */
    public CommonResult unRegisterService(String addr, String apiVersion) {
        return serviceFinder.unRegisterService(this, addr, apiVersion);
    }

    public CommonResult<List<Service>> useAndSubscribeService(List<SubscribeRequestValue> requestValueList, ServiceHandle serviceHandle) {

        CommonResult<List<Service>> commonResult = null;

        if (null == serviceHandle) {
            commonResult = new CommonResult();
            commonResult.setRet(ErrorCode.PARAM_INVALID);
            commonResult.setMsg("ServiceHandle can not be empty");
            return commonResult;
        }
        if (null == requestValueList || requestValueList.isEmpty()) {
            commonResult = new CommonResult();
            commonResult.setRet(ErrorCode.PARAM_INVALID);
            commonResult.setMsg("SubscribeRequestValue can not be empty");
            return commonResult;
        }

        for (SubscribeRequestValue value : requestValueList) {
            if (null == value || StringUtils.isNullOrEmpty(value.getApiVersion()) || StringUtils.isNullOrEmpty(value.getServiceName())) {
                commonResult = new CommonResult();
                commonResult.setRet(ErrorCode.PARAM_INVALID);
                StringBuffer buffer = new StringBuffer();
                buffer.append("SubscribeRequestValue is invalid:").append(value);
                commonResult.setMsg(buffer.toString());
                return commonResult;
            }
        }

        try {
            commonResult = serviceFinder.useAndSubscribeService(this, requestValueList, serviceHandle);
            if (ErrorCode.SUCCESS != commonResult.getRet() && bootConfig.isConfigCache()) {
                String cacheBasePath = PathUtils.getCacheFilePath(this, "service");
                commonResult = CacheUtils.getCacheServiceResult(requestValueList, cacheBasePath);
                if (initSuccess) {
                    this.getGlobalCache().taskQueue.add(TaskType.SERVICE);
                }
            }
        } catch (Exception e) {
            logger.error("", e);
        }

        return commonResult;
    }

    /**
     * 取消服务的订阅
     *
     * @param requestValue
     * @return
     */
    public CommonResult unSubscribeService(SubscribeRequestValue requestValue) {
        return serviceFinder.unSubscribeService(this, requestValue);
    }

    /**
     * 批量取消服务的订阅
     *
     * @param requestValueList
     * @return
     */

    public CommonResult unSubscribeService(List<SubscribeRequestValue> requestValueList) {
        return serviceFinder.unSubscribeService(this, requestValueList);
    }
}
