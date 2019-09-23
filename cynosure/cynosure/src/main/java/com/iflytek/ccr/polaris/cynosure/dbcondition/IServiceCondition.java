package com.iflytek.ccr.polaris.cynosure.dbcondition;

import com.iflytek.ccr.polaris.cynosure.domain.Service;
import com.iflytek.ccr.polaris.cynosure.request.service.AddServiceRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.service.EditServiceRequestBody;

import java.util.HashMap;
import java.util.List;

/**
 * 服务条件接口
 *
 * @author sctang2
 * @create 2017-12-10 15:24
 **/
public interface IServiceCondition {
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
     * @param clusterIds
     * @return
     */
    List<Service> findList(List<String> clusterIds);

    /**
     * 新增服务
     *
     * @param body
     * @return
     */
    Service add(AddServiceRequestBody body);

    /**
     * 根据id删除服务
     *
     * @param id
     * @return
     */
    int deleteById(String id);

    /**
     * 通过id查询服务
     *
     * @param id
     * @return
     */
    Service findById(String id);

    /**
     * 根据id更新服务
     *
     * @param id
     * @param body
     * @return
     */
    Service updateById(String id, EditServiceRequestBody body);

    /**
     * 根据服务名和集群id查询服务
     *
     * @param name
     * @param clusterId
     * @return
     */
    Service find(String name, String clusterId);

    /**
     * 通过id查询版本列表
     *
     * @param id
     * @return
     */
    Service findServiceVersionListById(String id);

    /**
     * 通过服务id查询服务、服务组、项目信息
     *
     * @param serviceId
     * @return
     */
    Service findServiceJoinGroupJoinProjectByServiceId(String serviceId);

    /**
     * 通过名称查询服务、服务组、项目信息
     *
     * @param project
     * @param group
     * @param service
     * @return
     */
    Service findServiceJoinGroupJoinProjectByName(String project, String group, String service);
}
