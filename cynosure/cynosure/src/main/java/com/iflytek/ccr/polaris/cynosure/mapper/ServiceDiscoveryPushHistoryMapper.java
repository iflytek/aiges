package com.iflytek.ccr.polaris.cynosure.mapper;

import com.iflytek.ccr.polaris.cynosure.domain.ServiceDiscoveryPushHistory;

import java.util.HashMap;
import java.util.List;

/**
 * 服务发现推送历史持久层接口
 *
 * @author sctang2
 * @create 2017-12-12 18:41
 **/
public interface ServiceDiscoveryPushHistoryMapper {
    /**
     * 新增服务发现推送历史
     *
     * @param serviceDiscoveryPushHistory
     * @return
     */
    int insert(ServiceDiscoveryPushHistory serviceDiscoveryPushHistory);

    /**
     * 删除服务推送历史
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

    /**
     * 通过id查询服务推送历史
     *
     * @param id
     * @return
     */
    ServiceDiscoveryPushHistory findById(String id);

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
}
