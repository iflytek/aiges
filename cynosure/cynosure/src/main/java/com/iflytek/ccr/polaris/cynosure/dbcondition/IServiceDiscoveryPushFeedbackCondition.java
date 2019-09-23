package com.iflytek.ccr.polaris.cynosure.dbcondition;

import com.iflytek.ccr.polaris.cynosure.domain.ServiceDiscoveryPushFeedback;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceDiscoveryPushFeedbackDetail;
import com.iflytek.ccr.polaris.cynosure.request.servicediscovery.ServiceDiscoveryFeedBackRequestBody;

import java.util.HashMap;
import java.util.List;

/**
 * 服务发现推送反馈条件接口
 *
 * @author sctang2
 * @create 2017-12-12 17:25
 **/
public interface IServiceDiscoveryPushFeedbackCondition {
    /**
     * 查询服务推送反馈总数
     *
     * @param map
     * @return
     */
    int findTotalCount(HashMap<String, Object> map);

    /**
     * 删除服务发现推送反馈
     *
     * @param pushId
     * @return
     */
    int delete(String pushId);

    /**
     * 通过pushIds删除服务配置推送反馈
     *
     * @param pushIds
     * @return
     */
    int deleteByPushIds(List<String> pushIds);

    /**
     * 查询服务推送反馈列表
     *
     * @param map
     * @return
     */
    List<ServiceDiscoveryPushFeedbackDetail> findList(HashMap<String, Object> map);

    /**
     * 新增服务发现推送反馈
     *
     * @param body
     * @return
     */
    ServiceDiscoveryPushFeedback add(ServiceDiscoveryFeedBackRequestBody body);
}
