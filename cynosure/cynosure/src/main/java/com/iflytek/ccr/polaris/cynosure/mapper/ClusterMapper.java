package com.iflytek.ccr.polaris.cynosure.mapper;

import com.iflytek.ccr.polaris.cynosure.domain.Cluster;

import java.util.HashMap;
import java.util.List;

/**
 * 集群持久化接口
 *
 * @author sctang2
 * @create 2017-11-15 17:15
 **/
public interface ClusterMapper {
    /**
     * 新增集群
     *
     * @param cluster
     * @return
     */
    int insert(Cluster cluster);

    /**
     * 删除集群
     *
     * @param id
     * @return
     */
    int deleteById(String id);

    /**
     * 更新集群
     *
     * @param cluster
     * @return
     */
    int updateById(Cluster cluster);

    /**
     * 根据id查询集群信息
     *
     * @param id
     * @return
     */
    Cluster findById(String id);

    /**
     * 查询集群信息
     *
     * @param cluster
     * @return
     */
    Cluster find(Cluster cluster);

    /**
     * 查询集群总数
     *
     * @param map
     * @return
     */
    int findTotalCount(HashMap<String, Object> map);

    /**
     * 查询集群列表
     *
     * @param map
     * @return
     */
    List<Cluster> findList(HashMap<String, Object> map);

    /**
     * 查询集群列表
     *
     * @param map
     * @return
     */
    List<Cluster> findClusterList(HashMap<String, Object> map);

    /**
     * 通过id查询集群列表
     *
     * @param id
     * @return
     */
    Cluster findServiceListById(String id);
}
