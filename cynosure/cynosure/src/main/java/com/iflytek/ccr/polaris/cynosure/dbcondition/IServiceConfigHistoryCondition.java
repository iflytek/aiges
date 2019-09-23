package com.iflytek.ccr.polaris.cynosure.dbcondition;

import com.iflytek.ccr.polaris.cynosure.domain.ServiceConfig;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceConfigHistory;

import java.util.HashMap;
import java.util.List;

/**
 * 服务配置历史条件接口
 *
 * @author sctang2
 * @create 2017-12-10 22:21
 **/
public interface IServiceConfigHistoryCondition {
    /**
     * 查询服务配置总数
     *
     * @param map
     * @return
     */
    int findTotalCount(HashMap<String, Object> map);

    /**
     * 查询服务配置历史
     *
     * @param map
     * @return
     */
    List<ServiceConfigHistory> findList(HashMap<String, Object> map);

    /**
     * 新增服务配置历史
     *
     * @param serviceConfig
     * @return
     */
    ServiceConfigHistory add(ServiceConfig serviceConfig);

    /**
     * 批量新增服务配置历史
     *
     * @param serviceConfigList
     */
    List<ServiceConfigHistory> batchAdd(List<ServiceConfig> serviceConfigList);

    /**
     * 通过id查询服务配置历史
     *
     * @param ids
     * @return
     */
    List<ServiceConfigHistory> findByIds(List<String> ids);

    /**
     * 通过configId删除配置历史
     *
     * @param configId
     * @return
     */
    int deleteByConfigId(String configId);

    /**
     * 通过configIds删除配置历史
     *
     * @param configIds
     * @return
     */
    int deleteByConfigIds(List<String> configIds);

    /**
     * 查询灰度配置总数
     *
     * @param map
     * @return
     */
    int findGrayTotalCount(HashMap<String, Object> map);

    /**
     * 查询灰度配置历史
     *
     * @param map
     * @return
     */
    List<ServiceConfigHistory> findGrayList(HashMap<String, Object> map);
}
