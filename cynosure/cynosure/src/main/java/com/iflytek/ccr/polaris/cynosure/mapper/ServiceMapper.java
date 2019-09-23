package com.iflytek.ccr.polaris.cynosure.mapper;

import com.iflytek.ccr.polaris.cynosure.domain.Service;

import java.util.HashMap;
import java.util.List;

/**
 * 服务持久层接口
 *
 * @author sctang2
 * @create 2017-11-16 17:47
 **/
public interface ServiceMapper {
    /**
     * 保存服务
     *
     * @param service
     * @return
     */
    int insert(Service service);

    /**
     * 删除服务
     *
     * @param id
     * @return
     */
    int deleteById(String id);

    /**
     * 更新服务
     *
     * @param service
     * @return
     */
    int updateById(Service service);

    /**
     * 通过id查询服务
     *
     * @param id
     * @return
     */
    Service findById(String id);

    /**
     * 查询服务信息
     *
     * @param service
     * @return
     */
    Service find(Service service);

    /**
     * 查询服务总数
     *
     * @param map
     * @return
     */
    int findTotalCount(HashMap<String, Object> map);

    /**
     * 查询服务列表
     *
     * @param map
     * @return
     */
    List<Service> findList(HashMap<String, Object> map);

    /**
     * 查询服务列表
     *
     * @param map
     * @return
     */
    List<Service> findServiceList(HashMap<String, Object> map);

    /**
     * 通过id查询版本列表
     *
     * @param id
     * @return
     */
    Service findServiceVersionListById(String id);

    /**
     * 查询服务、服务组、项目信息
     *
     * @param map
     * @return
     */
    Service findServiceJoinGroupJoinProjectByMap(HashMap<String, Object> map);
}
