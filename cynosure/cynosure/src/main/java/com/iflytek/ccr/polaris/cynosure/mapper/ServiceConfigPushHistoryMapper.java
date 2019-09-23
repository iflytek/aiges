package com.iflytek.ccr.polaris.cynosure.mapper;

import com.iflytek.ccr.polaris.cynosure.domain.ServiceConfigPushHistory;

import java.util.HashMap;
import java.util.List;

/**
 * 服务推送持久层接口
 *
 * @author sctang2
 * @create 2017-11-24 9:16
 **/
public interface ServiceConfigPushHistoryMapper {
    /**
     * 保存服务推送历史
     *
     * @param servicePushHistory
     * @return
     */
    int insert(ServiceConfigPushHistory servicePushHistory);

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
    ServiceConfigPushHistory findById(String id);

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
}
