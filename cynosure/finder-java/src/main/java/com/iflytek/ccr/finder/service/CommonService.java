package com.iflytek.ccr.finder.service;

import com.iflytek.ccr.finder.value.CommonResult;

public interface CommonService {

    /**
     * 注册服务发现/配置中心的消费节点
     *
     * @param path
     * @return
     */
    CommonResult registerConsumer(String path);

    /**
     * 取消注册
     *
     * @param path
     * @return
     */
    CommonResult unRegisterConsumer(String path);
}
