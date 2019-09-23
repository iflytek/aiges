package com.iflytek.ccr.polaris.cynosure.mapper;

import com.iflytek.ccr.polaris.cynosure.domain.ServiceConfigPushFeedback;

import java.util.HashMap;
import java.util.List;

/**
 * 服务配置推送反馈持久层接口
 *
 * @author sctang2
 * @create 2017-11-25 11:16
 **/
public interface ServiceConfigPushFeedbackMapper {
    /**
     * 保存服务配置推送反馈
     *
     * @param serviceConfigPushFeedback
     * @return
     */
    int insert(ServiceConfigPushFeedback serviceConfigPushFeedback);

    /**
     * 删除服务配置推送反馈
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
     * 查询服务配置推送反馈总数
     *
     * @param map
     * @return
     */
    int findTotalCount(HashMap<String, Object> map);

    /**
     * 查询服务配置推送反馈列表
     *
     * @param map
     * @return
     */
    List<ServiceConfigPushFeedback> findList(HashMap<String, Object> map);
}
