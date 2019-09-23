package com.iflytek.ccr.polaris.cynosure.mapper;

import com.iflytek.ccr.polaris.cynosure.domain.ServiceApiVersion;

import java.util.HashMap;
import java.util.List;

/**
 * 服务版本持久层接口
 *
 * @author sctang2
 * @create 2017-11-17 14:52
 **/
public interface ServiceApiVersionMapper {
    /**
     * 保存版本
     *
     * @param serviceApiVersion
     * @return
     */
    int insert(ServiceApiVersion serviceApiVersion);

    /**
     * 删除 版本
     *
     * @param id
     * @return
     */
    int deleteById(String id);

    /**
     * 更新版本
     *
     * @param serviceApiVersion
     * @return
     */
    int updateById(ServiceApiVersion serviceApiVersion);

    /**
     * 通过id查询版本
     *
     * @param id
     * @return
     */
    ServiceApiVersion findById(String id);

    /**
     * 查询版本
     *
     * @param serviceApiVersion
     * @return
     */
    ServiceApiVersion find(ServiceApiVersion serviceApiVersion);

    /**
     * 查询版本总数
     *
     * @param map
     * @return
     */
    int findTotalCount(HashMap<String, Object> map);

    /**
     * 查询版本列表
     *
     * @param map
     * @return
     */
    List<ServiceApiVersion> findList(HashMap<String, Object> map);

    /**
     * 查询服务版本列表
     *
     * @param map
     * @return
     */
    List<ServiceApiVersion> findServiceApiVersionList(HashMap<String, Object> map);

    /**
     * 通过id查询服务配置列表
     *
     * @param id
     * @return
     */
    ServiceApiVersion findServiceConfigListById(String id);

    /**
     * 通过id查询服务版本、服务、服务组、项目信息
     *
     * @param id
     * @return
     */
    ServiceApiVersion findVersionJoinServiceJoinGroupJoinProjectById(String id);

    /**
     * 通过服务id集合查询所有api版本
     * @param map
     * @return
     */
    List<ServiceApiVersion> findApiVersionList(HashMap<String, Object> map);
}
