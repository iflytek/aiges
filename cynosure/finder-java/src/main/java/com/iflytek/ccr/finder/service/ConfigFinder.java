package com.iflytek.ccr.finder.service;

import com.iflytek.ccr.finder.FinderManager;
import com.iflytek.ccr.finder.handler.ConfigChangedHandler;
import com.iflytek.ccr.finder.value.CommonResult;
import com.iflytek.ccr.finder.value.Config;

import java.util.List;

public interface ConfigFinder {

    /**
     * 订阅配置并返回当前配置内容
     *
     * @param finderManager
     * @param configNameList
     * @param configChangedHandler
     * @return
     */
    CommonResult<List<Config>> useAndSubscribeConfig(FinderManager finderManager, List<String> configNameList, ConfigChangedHandler configChangedHandler);

    /**
     * 订阅配置并返回当前配置内容
     *
     * @param finderManager
     * @param configNameList
     * @param configChangedHandler
     * @return
     */
    CommonResult<List<Config>> useAndSubscribeConfig(FinderManager finderManager, List<String> configNameList, ConfigChangedHandler configChangedHandler,boolean isRecover);

    /**
     * 取消配置的订阅
     *
     * @param finderManager
     * @param configName
     * @return
     */
    CommonResult unSubscribeConfig(FinderManager finderManager, String configName);

    /**
     * 获取当前的配置对象
     *
     * @param finderManager
     * @param basePath
     * @param configNameList
     * @return
     */
    CommonResult<List<Config>> getCurrentConfig(FinderManager finderManager, String basePath, List<String> configNameList);

}
