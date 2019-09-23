package com.iflytek.ccr.polaris.cynosure.dbcondition;

import com.iflytek.ccr.polaris.cynosure.companionservice.domain.PushResult;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceDiscoveryPushHistory;

import java.util.HashMap;
import java.util.List;

/**
 * 服务发现推送历史条件接口
 *
 * @author sctang2
 * @create 2017-12-12 17:25
 **/
public interface IServiceDiscoveryPushHistoryCondition {
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
    List<ServiceDiscoveryPushHistory> findList(HashMap<String, Object> map);

    /**
     * 新增服务发现推送历史
     *
     * @param project
     * @param cluster
     * @param service
     * @param cacheCenterPushResult
     * @return
     */
    ServiceDiscoveryPushHistory add(String project, String cluster, String service, String apiVersion, PushResult cacheCenterPushResult);

    /**
     * 通过id查询服务推送历史
     *
     * @param id
     * @return
     */
    ServiceDiscoveryPushHistory findById(String id);

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
