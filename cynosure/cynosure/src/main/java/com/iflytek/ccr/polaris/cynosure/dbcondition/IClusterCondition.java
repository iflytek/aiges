package com.iflytek.ccr.polaris.cynosure.dbcondition;

import com.iflytek.ccr.polaris.cynosure.domain.Cluster;
import com.iflytek.ccr.polaris.cynosure.request.cluster.AddClusterRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.cluster.EditClusterRequestBody;

import java.util.HashMap;
import java.util.List;

/**
 * 服务组条件接口
 *
 * @author sctang2
 * @create 2017-12-10 13:51
 **/
public interface IClusterCondition {
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
     * @param projectIds
     * @return
     */
    List<Cluster> findList(List<String> projectIds);

    /**
     * 根据id查询集群信息
     *
     * @param id
     * @return
     */
    Cluster findById(String id);

    /**
     * 根据id和集群名称查询集群信息
     *
     * @param projectId
     * @param name
     * @return
     */
    Cluster find(String projectId, String name);

    /**
     * 创建集群
     *
     * @param body
     * @return
     */
    Cluster add(AddClusterRequestBody body);

    /**
     * 根据id更新集群
     *
     * @param id
     * @param body
     * @return
     */
    Cluster updateById(String id, EditClusterRequestBody body);

    /**
     * 通过id查询服务列表
     *
     * @param id
     * @return
     */
    Cluster findServiceListById(String id);

    /**
     * 通过id删除集群
     *
     * @param id
     * @return
     */
    int deleteById(String id);
}
