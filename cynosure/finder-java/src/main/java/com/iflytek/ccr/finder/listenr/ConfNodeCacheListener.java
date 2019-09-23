package com.iflytek.ccr.finder.listenr;

import com.iflytek.ccr.finder.FinderManager;
import com.iflytek.ccr.finder.constants.Constants;
import com.iflytek.ccr.finder.handler.ServiceHandle;
import com.iflytek.ccr.finder.utils.*;
import com.iflytek.ccr.finder.value.ErrorCode;
import com.iflytek.ccr.finder.value.Service;
import com.iflytek.ccr.finder.value.ServiceConfig;
import com.iflytek.ccr.finder.value.ZkDataValue;
import com.iflytek.ccr.zkutil.ZkHelper;
import org.apache.curator.framework.recipes.cache.NodeCacheListener;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * 服务发现conf节点监听
 */
public class ConfNodeCacheListener implements NodeCacheListener {

    private static final Logger logger = LoggerFactory.getLogger(ConfNodeCacheListener.class);

    private FinderManager finderManager;
    private String serviceName;
    private ServiceHandle serviceHandle;
    private String apiVersion;

    public ConfNodeCacheListener(FinderManager finderManager, String serviceName, ServiceHandle serviceHandle, String apiVersion) {
        this.finderManager = finderManager;
        this.serviceName = serviceName;
        this.serviceHandle = serviceHandle;
        this.apiVersion = apiVersion;
    }

    @Override
    public void nodeChanged() throws Exception {
        logger.info("nodeChanged");
        try {
            String confPath = PathUtils.getConfPath(serviceName, apiVersion);
            ZkHelper zkHelper = ZkInstanceUtil.getInstance();
            byte[] data = zkHelper.getByteData(confPath);
            ZkDataValue zkDataValue = ByteUtil.parseZkData(data);
            if (ErrorCode.SUCCESS == zkDataValue.getRet()) {
                String json = new String(zkDataValue.getRealData(), Constants.DEFAULT_CHARSET);
                String cacheKey = serviceName + "_" + apiVersion;
                ServiceConfig config = finderManager.getServiceCache().confCacheMap.get(cacheKey);
                if (isIgnore(json, config)) {
                    //没有发生实际变更，忽略
                    return;
                } else {
                    config.setJsonConfig(json);
                    //刷新缓存文件
                    refreshCacheFile(cacheKey, json);
                }

                boolean isSuccess = serviceHandle.onServiceConfigChanged(serviceName, json);
                //如果pushid非法，则不进行反馈
                if (StringUtils.isNullOrEmpty(zkDataValue.getPushId())) {
                    return;
                }
                long updateTime = System.currentTimeMillis();
                if (isSuccess) {
                    RemoteUtil.pushServiceFeedback(finderManager, zkDataValue.getPushId(), "", Constants.UPDATE_STATUS_SUCCESS, Constants.LOAD_STATUS_SUCCESS, String.valueOf(updateTime), String.valueOf(System.currentTimeMillis()), apiVersion, Constants.KEY_SERVICE_CONFIG_CHANGE);
                } else {
                    RemoteUtil.pushServiceFeedback(finderManager, zkDataValue.getPushId(), "", Constants.UPDATE_STATUS_FAIL, Constants.LOAD_STATUS_FAIL, String.valueOf(updateTime), String.valueOf(System.currentTimeMillis()), apiVersion, Constants.KEY_SERVICE_CONFIG_CHANGE);
                }

            } else {
                logger.warn(String.format("parseZkData error,confPath:%s", confPath));
            }
        } catch (Exception e) {
            logger.error(String.format("ConfNodeCacheListener error:%s", e.getMessage()), e);
        }


    }

    /**
     * 判断是否需要忽略（如果变化前后，没有实质性的变化，则忽略）
     *
     * @param json
     * @param config
     * @return
     */
    private boolean isIgnore(String json, ServiceConfig config) {
        boolean flag = false;
        if (null == json && (null != config && null == config.getJsonConfig())) {
            flag = true;
        } else {
            flag = null != config && null != config.getJsonConfig() && config.getJsonConfig().equals(json);
        }
        return flag;
    }

    /**
     * 刷新缓存文件
     *
     * @param cacheKey
     * @param jsonConfig
     */
    private void refreshCacheFile(String cacheKey, String jsonConfig) {
        if (finderManager.getBootConfig().isServiceCache()) {
            try {
                Service service = (Service) FinderFileUtils.readObjectFromFile(PathUtils.getCacheFilePath(finderManager, "service") + cacheKey);
                ServiceConfig serviceConfig = new ServiceConfig();
                serviceConfig.setJsonConfig(jsonConfig);
                service.setServiceConfig(serviceConfig);
                FinderFileUtils.writeObjectToFile(PathUtils.getCacheFilePath(finderManager, "service") + cacheKey, service);
            } catch (Exception e) {
                logger.error("", e);
            }
        }
    }

}
