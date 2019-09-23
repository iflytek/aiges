package com.iflytek.ccr.polaris.cynosure.mapper;

import com.iflytek.ccr.polaris.cynosure.domain.ServiceDiscoveryPushFeedback;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceDiscoveryPushFeedbackDetail;

import java.util.HashMap;
import java.util.List;

/**
 * 服务发现推送反馈持久层接口
 *
 * @author sctang2
 * @create 2017-12-12 19:18
 **/
public interface ServiceDiscoveryPushFeedbackMapper {
    /**
     * 新增服务发现推送反馈
     *
     * @param serviceDiscoveryPushFeedback
     * @return
     */
    int insert(ServiceDiscoveryPushFeedback serviceDiscoveryPushFeedback);

    /**
     * 删除服务发现推送反馈
     *
     * @param pushId
     * @return
     */
    int deleteByPushId(String pushId);

    /**
     * 通过pushIds删除服务配置推送反馈
     *
     * @param pushIds
     * @return
     */
    int deleteByPushIds(List<String> pushIds);

    /**
     * 查询服务推送反馈总数
     *
     * @param map
     * @return
     */
    int findTotalCount(HashMap<String, Object> map);

    /**
     * 查询服务推送反馈列表
     *
     * @param map
     * @return
     */
    List<ServiceDiscoveryPushFeedbackDetail> findList(HashMap<String, Object> map);
}
