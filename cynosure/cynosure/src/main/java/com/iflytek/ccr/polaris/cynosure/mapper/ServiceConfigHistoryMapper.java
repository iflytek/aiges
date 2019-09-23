package com.iflytek.ccr.polaris.cynosure.mapper;

import com.iflytek.ccr.polaris.cynosure.domain.ServiceConfigHistory;

import java.util.HashMap;
import java.util.List;

/**
 * 服务配置历史持久层接口
 *
 * @author sctang2
 * @create 2017-11-22 10:24
 **/
public interface ServiceConfigHistoryMapper {
    /**
     * 保存服务配置历史
     *
     * @param serviceConfigHistory
     * @return
     */
    int insert(ServiceConfigHistory serviceConfigHistory);

    /**
     * 批量新增服务配置历史
     *
     * @param serviceConfigHistoryList
     * @return
     */
    int batchInsert(List<ServiceConfigHistory> serviceConfigHistoryList);

    /**
     * 删除服务配置历史
     *
     * @param id
     * @return
     */
    int deleteById(String id);

    /**
     * 通过configId删除服务配置历史
     *
     * @param configId
     * @return
     */
    int deleteByConfigId(String configId);

    /**
     * 通过configIds删除服务配置历史
     *
     * @param configIds
     * @return
     */
    int deleteByConfigIds(List<String> configIds);

    /**
     * 更新服务配置历史
     *
     * @param serviceConfigHistory
     * @return
     */
    int updateById(ServiceConfigHistory serviceConfigHistory);

    /**
     * 通过ids查询服务配置历史
     *
     * @param ids
     * @return
     */
    List<ServiceConfigHistory> findByIds(List<String> ids);

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
