package com.iflytek.ccr.polaris.cynosure.dbcondition;

import com.iflytek.ccr.polaris.cynosure.companionservice.domain.PushResult;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceConfig;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceConfigPushHistory;

import java.util.HashMap;
import java.util.List;

/**
 * 服务推送历史条件接口
 *
 * @author sctang2
 * @create 2017-12-10 22:23
 **/
public interface IServiceConfigPushHistoryCondition {
    /**
     * 查询服务推送历史总数
     *
     * @param map
     * @return
     */
    int findTotalCount(HashMap<String, Object> map);

    /**
     * 查询服务推送历史列表
     *
     * @param map
     * @return
     */
    List<ServiceConfigPushHistory> findList(HashMap<String, Object> map);

    /**
     * 新增
     *
     * @param serviceConfig
     * @param cacheCenterPushResult
     * @return
     */
    ServiceConfigPushHistory add(ServiceConfig serviceConfig, PushResult cacheCenterPushResult);

    /**
     * 批量新增
     *
     * @param serviceConfigList
     * @param cacheCenterPushResult
     * @return
     */
    ServiceConfigPushHistory add(List<ServiceConfig> serviceConfigList, PushResult cacheCenterPushResult);

    /**
     * 通过id查询服务推送历史
     *
     * @param id
     * @return
     */
    ServiceConfigPushHistory findById(String id);

    /**
     * 通过id删除服务推送历史
     *
     * @param id
     * @return
     */
    int deleteById(String id);

    /**
     * 通过ids删除服务推送历史
     *
     * @param ids
     * @return
     */
    int deleteByIds(List<String> ids);
}
