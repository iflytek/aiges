package com.iflytek.ccr.polaris.cynosure.mapper;

import com.iflytek.ccr.polaris.cynosure.domain.ServiceConfig;
import org.apache.ibatis.annotations.Param;

import java.util.HashMap;
import java.util.List;

/**
 * 服务配置持久层接口
 *
 * @author sctang2
 * @create 2017-11-21 11:33
 **/
public interface ServiceConfigMapper {
    /**
     * 新增服务配置
     *
     * @param serviceConfig
     * @return
     */
    int insert(ServiceConfig serviceConfig);

    /**
     * 批量新增服务配置
     *
     * @param serviceConfigList
     * @return
     */
    int batchInsert(List<ServiceConfig> serviceConfigList);

    /**
     * 批量更新服务配置
     *
     * @param serviceConfigList
     * @return
     */
    int batchUpdate(List<ServiceConfig> serviceConfigList);

    /**
     * 删除服务配置
     *
     * @param id
     * @return
     */
    int deleteById(String id);

    /**
     * 通过ids删除服务配置
     *
     * @param ids
     * @return
     */
    int deleteByIds(List<String> ids);

    /**
     * 更新服务配置
     *
     * @param serviceConfig
     * @return
     */
    int updateById(ServiceConfig serviceConfig);

    /**
     * 通过id查询服务配置
     *
     * @param id
     * @return
     */
    ServiceConfig findById(String id);

    /**
     * 查询服务配置信息
     *
     * @param serviceConfig
     * @return
     */
    ServiceConfig find(ServiceConfig serviceConfig);

    /**
     * 查询最新版本的配置信息
     *
     * @param map
     * @return
     */
    List<ServiceConfig> findNewList(HashMap<String, Object> map);

    /**
     * 查询服务配置列表
     *
     * @param map
     * @return
     */
    List<ServiceConfig> findListByMap(HashMap<String, Object> map);

    /**
     * 查询服务配置总数
     *
     * @param map
     * @return
     */
    int findTotalCount(HashMap<String, Object> map);

    /**
     * 查询服务配置列表
     *
     * @param map
     * @return
     */
    List<ServiceConfig> findList(HashMap<String, Object> map);

    /**
     * 通过版本Id查询该版本下配置的集合
     * @param map
     * @return
     */
    List<ServiceConfig> findConfigsByVersionId(HashMap<String, Object> map);

    /**
     * 通过map查询配置、版本、服务、集群、项目信息
     *
     * @param map
     * @return
     */
    List<ServiceConfig> findConfigJoinVersionJoinServiceJoinClusterJoinProjectByMap(HashMap<String, Object> map);

    /**
     * 通过ids查询服务配置列表
     *
     * @param ids
     * @return
     */
    List<ServiceConfig> findListByIds(@Param("ids") List<String> ids);

    /**
     * 通过ids查询灰度配置列表
     *
     * @param ids
     * @return
     */
    List<ServiceConfig> findListByGrayIds(@Param("ids") List<String> ids);

    /**
     * 通过id查询灰度配置列表
     *
     * @param grayId
     * @return
     */
    List<ServiceConfig> findListByGrayId(String grayId);

    /**
     * 删除灰度配置
     *
     * @param grayId
     * @return
     */
    int deleteByGrayId(String grayId);

    /**
     * 查询服务灰度配置总数
     *
     * @param map
     * @return
     */
    int findGrayTotalCount(HashMap<String, Object> map);

    /**
     * 查询服务灰度配置列表
     *
     * @param map
     * @return
     */
    List<ServiceConfig> findGrayList(HashMap<String, Object> map);

    /**
     * 查询服务配置信息(用于校验配置文件唯一性)
     *
     * @param serviceConfig
     * @return
     */
    ServiceConfig findOnlyConfig(ServiceConfig serviceConfig);
}
