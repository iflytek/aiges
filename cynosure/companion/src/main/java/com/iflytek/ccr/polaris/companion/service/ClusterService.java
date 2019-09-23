package com.iflytek.ccr.polaris.companion.service;

import com.iflytek.ccr.nakedserver.http.HttpBody;
import com.iflytek.ccr.polaris.companion.common.JsonResult;

import java.util.List;
import java.util.Map;

public interface ClusterService {

    /**
     * 配置推送
     *
     * @param map
     * @return
     */
    JsonResult pushConfig(Map<String, List<HttpBody>> map);

    /**
     * 灰度配置推送
     *
     * @param map
     * @return
     */
    JsonResult grayPushConfig(Map<String, List<HttpBody>> map);

    /**
     * 灰度配置推送
     *
     * @param map
     * @return
     */
    JsonResult delGrayGroup(Map<String, List<HttpBody>> map);

    /**
     * 服务配置信息推送(userData\sdkData两部分信息的推送)
     *
     * @param map
     * @return
     */
    JsonResult pushServiceConfig(Map<String, String> map);

    /**
     * 服务实例配置信息推送(userData\sdkData两部分信息的推送)
     *
     * @param map
     * @return
     */
    JsonResult pushServiceInstanceConfig(Map<String, List<HttpBody>> map);

    /**
     * 实例级配置信息保存到zookeeper
     *
     * @param map
     * @return
     */
    JsonResult pushInstanceConfig(Map<String, List<HttpBody>> map);
}
