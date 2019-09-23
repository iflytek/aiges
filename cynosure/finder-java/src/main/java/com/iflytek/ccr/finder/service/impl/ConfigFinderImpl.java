package com.iflytek.ccr.finder.service.impl;

import com.iflytek.ccr.finder.FinderManager;
import com.iflytek.ccr.finder.constants.Constants;
import com.iflytek.ccr.finder.handler.ConfigChangedHandler;
import com.iflytek.ccr.finder.listenr.GrayNodeCacheListener;
import com.iflytek.ccr.finder.service.CommonService;
import com.iflytek.ccr.finder.service.ConfigFinder;
import com.iflytek.ccr.finder.service.GrayConfigService;
import com.iflytek.ccr.finder.utils.*;
import com.iflytek.ccr.finder.value.*;
import com.iflytek.ccr.zkutil.ZkHelper;
import org.apache.curator.framework.recipes.cache.NodeCache;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.ArrayList;
import java.util.List;

public class ConfigFinderImpl implements ConfigFinder {

    private static final Logger logger = LoggerFactory.getLogger(ConfigFinderImpl.class);

    CommonService commonService = new CommonServiceImpl();

    GrayConfigService grayConfigService = new GrayConfigServiceImpl();


    @Override
    public CommonResult<List<Config>> useAndSubscribeConfig(FinderManager finderManager, List<String> configNameList, ConfigChangedHandler configChangedHandler) {
        logger.info(String.format("configNameList:%s", configNameList));
        return useAndSubscribeConfig(finderManager, configNameList, configChangedHandler, false);
    }

    @Override
    public CommonResult<List<Config>> useAndSubscribeConfig(FinderManager finderManager, List<String> configNameList, ConfigChangedHandler configChangedHandler, boolean isRecover) {
        logger.info(String.format("configNameList:%s", configNameList));
        for (String name : configNameList) {
            finderManager.getGlobalCache().configListenerMap.put(name, null);
            if (!isRecover) {
                finderManager.getGlobalCache().initMap.put(name, true);
            }
        }
        finderManager.getGlobalCache().setConfigChangedHandler(configChangedHandler);
        finderManager.getGlobalCache().setConfigNameList(configNameList);
        CommonResult commonResult = null;
        try {
            String rootConfigPath = ConfigManager.getInstance().getStringConfigByKey(Constants.CONFIG_PATH);
            //获取gray节点的路径
            String grayConfigPath = rootConfigPath + Constants.GRAY_NODE_PATH;

            ZkHelper zkHelper = ZkInstanceUtil.getInstance();
            if (!zkHelper.checkExists(grayConfigPath)) {
                zkHelper.addPersistent(grayConfigPath, "");
            }

            //监听gray节点的变化
            GrayNodeCacheListener grayNodeCacheListener = new GrayNodeCacheListener(finderManager, configNameList);

            zkHelper.addListener(grayNodeCacheListener, grayConfigPath, false);

            //获取当前组件实例的配置根目录（灰度or正常）
            GrayConfigValue grayConfigValue = grayConfigService.getGrayServer(finderManager, grayConfigService.parseGrayData(grayConfigPath));
            String basePath = null;
            if (null == grayConfigValue) {
                basePath = rootConfigPath;
            } else {
                basePath = grayConfigPath + "/" + grayConfigValue.getGroupId();
            }

            commonResult = getCurrentConfig(finderManager, basePath, configNameList);
        } catch (Exception e) {
            logger.error(String.format("useAndSubscribeConfigSupportGray error:%s", e.getMessage()), e);
            commonResult = new CommonResult();
            commonResult.setRet(ErrorCode.INTERNAL_EXCEPTION);
            commonResult.setMsg(e.getMessage());
        }
        return commonResult;
    }

    @Override
    public CommonResult unSubscribeConfig(FinderManager finderManager, String configName) {
        CommonResult commonResult = new CommonResult();
        try {
            Object o = finderManager.getGlobalCache().configListenerMap.get(configName);
            if (o instanceof NodeCache) {
                ((NodeCache) o).close();
            }
            finderManager.getGlobalCache().configListenerMap.remove(configName);

            if (finderManager.getGlobalCache().configListenerMap.isEmpty()) {
                finderManager.getGlobalCache().monitorPathList.remove(finderManager.getGlobalCache().getConfigConsumerPath());
                commonService.unRegisterConsumer(finderManager.getGlobalCache().getConfigConsumerPath());
            }
            commonResult.setRet(ErrorCode.SUCCESS);
        } catch (Exception e) {
            logger.error(String.format("unSubscribeConfig error:%s", e.getMessage()), e);
            commonResult.setRet(ErrorCode.INTERNAL_EXCEPTION);
            commonResult.setMsg(e.getMessage());
        }

        return commonResult;
    }

    @Override
    public CommonResult<List<Config>> getCurrentConfig(FinderManager finderManager, String basePath, List<String> configNameList) {
        logger.info(String.format("configNameList:%s", configNameList));
        CommonResult commonResult = new CommonResult();
        commonResult.setRet(ErrorCode.SUCCESS);
        try {
            ZkHelper zkHelper = ZkInstanceUtil.getInstance();
            List<Config> configs = new ArrayList<>();
            for (String name : configNameList) {
                String path = basePath + "/" + name;
                Config config = new Config();
                if (zkHelper.checkExists(path)) {
                    byte[] data = zkHelper.getByteData(path);
                    ZkDataValue zkDataValue = ByteUtil.parseZkData(data);
                    if (ErrorCode.SUCCESS == zkDataValue.getRet()) {
                        config.setFile(zkDataValue.getRealData());
                        FinderFileUtils.writeByteArrayToFile(PathUtils.getCacheFilePath(finderManager, "config") + name, config.getFile());
                    } else {
                        commonResult.setRet(zkDataValue.getRet());
                        commonResult.setMsg(zkDataValue.getDesc());
                        logger.error(String.format("data parse error:%s", zkDataValue.toString()));
                    }
                } else {
                    logger.error("path:" + path + " does not exists");
                    commonResult.setRet(ErrorCode.CONFIG_MISS_FILE);
                    commonResult.setMsg("file:" + path + " does not exists");
                    return commonResult;
                }
                config.setName(name);
                if (config.getName().endsWith(".toml") || config.getName().endsWith(".TOML")) {
                    config.setConfigMap(FinderFileUtils.parseTomlFile(config.getFile()));
                }
                configs.add(config);
            }
            commonResult.setData(configs);
            logger.info(String.format("useConfig:%s", JacksonUtils.toJson(commonResult)));
        } catch (Exception e) {
            logger.error("", e);
            commonResult.setRet(ErrorCode.INTERNAL_EXCEPTION);
            commonResult.setMsg(e.getMessage());
        }
        return commonResult;
    }
}
