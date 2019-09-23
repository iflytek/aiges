package com.iflytek.ccr.finder.service;

import com.iflytek.ccr.finder.FinderManager;
import com.iflytek.ccr.finder.handler.ServiceHandle;
import com.iflytek.ccr.finder.value.CommonResult;
import com.iflytek.ccr.finder.value.Service;
import com.iflytek.ccr.finder.value.ServiceConfig;
import com.iflytek.ccr.finder.value.SubscribeRequestValue;

import java.util.List;

/**
 * 服务发现接口
 */
public interface ServiceFinder {

    /**
     * 服务注册
     *
     * @return
     */
    CommonResult registerService(FinderManager finderManager, String apiVersion, String addr);

    /**
     * 取消服务注册
     *
     * @return
     */
    CommonResult unRegisterService(FinderManager finderManager, String apiVersion);

    /**
     * 取消服务注册
     *
     * @return
     */
    CommonResult unRegisterService(FinderManager finderManager, String addr, String apiVersion);

    /**
     * 订阅服务
     *
     * @param finderManager
     * @param requestValueList
     * @param serviceHandle
     * @return
     */
    CommonResult<List<Service>> useAndSubscribeService(FinderManager finderManager, List<SubscribeRequestValue> requestValueList, ServiceHandle serviceHandle);

    /**
     * 取消订阅
     *
     * @param finderManager
     * @param requestValue
     * @return
     */
    CommonResult unSubscribeService(FinderManager finderManager, SubscribeRequestValue requestValue);

    /**
     * 批量取消订阅
     *
     * @param finderManager
     * @param requestValueList
     * @return
     */
    CommonResult unSubscribeService(FinderManager finderManager, List<SubscribeRequestValue> requestValueList);


    CommonResult<ServiceConfig> getServiceConfig(FinderManager finderManager, String serviceName, String apiVersion);
}
