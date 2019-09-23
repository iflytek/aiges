package com.iflytek.ccr.polaris.cynosure.mapper;

import com.iflytek.ccr.polaris.cynosure.domain.ServiceVersion;

import java.util.HashMap;
import java.util.List;

/**
 * 服务版本持久层接口
 *
 * @author sctang2
 * @create 2017-11-17 14:52
 **/
public interface ServiceVersionMapper {
    /**
     * 保存版本
     *
     * @param serviceVersion
     * @return
     */
    int insert(ServiceVersion serviceVersion);

    /**
     * 删除版本
     *
     * @param id
     * @return
     */
    int deleteById(String id);

    /**
     * 更新版本
     *
     * @param serviceVersion
     * @return
     */
    int updateById(ServiceVersion serviceVersion);

    /**
     * 通过id查询版本
     *
     * @param id
     * @return
     */
    ServiceVersion findById(String id);

    /**
     * 查询版本
     *
     * @param serviceVersion
     * @return
     */
    ServiceVersion find(ServiceVersion serviceVersion);

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
    List<ServiceVersion> findList(HashMap<String, Object> map);

    /**
     * 查询服务版本列表
     *
     * @param map
     * @return
     */
    List<ServiceVersion> findServiceVersionList(HashMap<String, Object> map);

    /**
     * 通过id查询服务配置列表
     *
     * @param id
     * @return
     */
    ServiceVersion findServiceConfigListById(String id);

    /**
     * 通过id查询服务版本、服务、服务组、项目信息
     *
     * @param id
     * @return
     */
    ServiceVersion findVersionJoinServiceJoinGroupJoinProjectById(String id);
}
