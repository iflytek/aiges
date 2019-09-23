package com.iflytek.ccr.finder.cache;

import com.iflytek.ccr.finder.utils.FinderFileUtils;
import com.iflytek.ccr.finder.value.*;

import java.util.ArrayList;
import java.util.List;

/**
 * 缓存工具类
 */
public class CacheUtils {

    public static CommonResult<List<Service>> getServiceCacheResult(String cacheBasePath) {
        CommonResult<List<Service>> commonResult = (CommonResult<List<Service>>) FinderFileUtils.readObjectFromFile(cacheBasePath);
        if (null == commonResult) {
            commonResult = new CommonResult<>();
            commonResult.setMsg("read sevice file fail");
            commonResult.setRet(ErrorCode.READ_FILE_FAIL);
        }
        return commonResult;
    }

    /**
     * 从缓存中读取配置文件
     *
     * @param configNameList
     * @param cacheBasePath
     * @return
     */
    public static CommonResult<List<Config>> getCacheConfigResult(List<String> configNameList, String cacheBasePath) {
        CommonResult<List<Config>> commonResult = null;
        List<Config> configs = new ArrayList<>();
        for (String configName : configNameList) {
            String fileName = cacheBasePath + configName;
            byte[] fileBtye = FinderFileUtils.readFileToByteArray(fileName);
            if (null == fileBtye) {
                continue;
            }
            Config config = new Config();
            config.setName(configName);
            config.setFile(fileBtye);
            if (config.getName().endsWith(".toml") || config.getName().endsWith(".TOML")) {
                config.setConfigMap(FinderFileUtils.parseTomlFile(config.getFile()));
            }
            configs.add(config);
        }
        if (!configs.isEmpty()) {
            commonResult = new CommonResult<>();
            commonResult.setData(configs);
            commonResult.setRet(ErrorCode.READ_CACHE_SUCCESS);
        } else {
            commonResult = new CommonResult<>();
            commonResult.setMsg("read cache config file fail");
            commonResult.setRet(ErrorCode.READ_FILE_FAIL);
        }
        return commonResult;
    }

    /**
     * 从缓存中读取配置文件
     *
     * @param requestValueList
     * @param cacheBasePath
     * @return
     */
    public static CommonResult<List<Service>> getCacheServiceResult(List<SubscribeRequestValue> requestValueList, String cacheBasePath) {
        CommonResult<List<Service>> commonResult = null;
        List<Service> serviceList = new ArrayList<>();
        for (SubscribeRequestValue value : requestValueList) {
            String fileName = cacheBasePath + value.getCacheKey();
            Service service = (Service) FinderFileUtils.readObjectFromFile(fileName);
            serviceList.add(service);
        }
        if (!serviceList.isEmpty()) {
            commonResult = new CommonResult<>();
            commonResult.setRet(ErrorCode.READ_CACHE_SUCCESS);
            commonResult.setData(serviceList);
        } else {
            commonResult = new CommonResult<>();
            commonResult.setMsg("read cache service file fail");
            commonResult.setRet(ErrorCode.READ_FILE_FAIL);
        }

        return commonResult;
    }
}
