package com.iflytek.ccr.polaris.cynosure.dbcondition;

import com.iflytek.ccr.polaris.cynosure.domain.ServiceConfigPushFeedback;
import com.iflytek.ccr.polaris.cynosure.request.serviceconfig.ServiceConfigFeedBackRequestBody;

import java.util.HashMap;
import java.util.List;

/**
 * 服务配置推送反馈条件接口
 *
 * @author sctang2
 * @create 2017-12-11 9:05
 **/
public interface IServiceConfigPushFeedbackCondition {
    /**
     * 查询服务配置推送反馈总数
     *
     * @param map
     * @return
     */
    int findTotalCount(HashMap<String, Object> map);

    /**
     * 删除服务配置推送反馈
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
     * 查询服务配置推送反馈列表
     *
     * @param map
     * @return
     */
    List<ServiceConfigPushFeedback> findList(HashMap<String, Object> map);

    /**
     * 新增更新反馈
     *
     * @param body
     * @return
     */
    ServiceConfigPushFeedback add(ServiceConfigFeedBackRequestBody body);
}
