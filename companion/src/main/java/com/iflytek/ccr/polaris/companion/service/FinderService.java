package com.iflytek.ccr.polaris.companion.service;

import com.iflytek.ccr.polaris.companion.common.ConfigFeedBackValue;
import com.iflytek.ccr.polaris.companion.common.JsonResult;
import com.iflytek.ccr.polaris.companion.common.ServiceFeedBackValue;
import com.iflytek.ccr.polaris.companion.common.ServiceValue;

import java.util.List;

public interface FinderService {

    /**
     * 查询zk临时路径：存放当前zk集群的节点信息
     */
    JsonResult queryZkPath();


    /**
     * 推送反馈配置到缓存
     *
     * @return
     */
    boolean pushConfigFeedback(ConfigFeedBackValue feedBackValue);

    /**
     * 推送服务反馈信息到缓存
     * @param feedBackValue
     * @return
     */
    boolean pushServiceFeedback(ServiceFeedBackValue feedBackValue);

    /**
     * serviceDiscovery
     *
     * @param serviceValue
     * @return
     */
    JsonResult serviceDiscovery(ServiceValue serviceValue);

    /**
     * 推送反馈信息到网站
     *
     * @return
     */
    JsonResult queryServicePath(String path);

    /**
     * 获取zkr集群地址，数组形式
     *
     * @param str
     * @return
     */
    List<String> getZkrAddrs(String str);

}
