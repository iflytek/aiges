package com.iflytek.ccr.polaris.cynosure.dbcondition;

import com.iflytek.ccr.polaris.cynosure.domain.ServiceVersion;
import com.iflytek.ccr.polaris.cynosure.request.serviceversion.AddServiceVersionRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceversion.EditServiceVersionRequestBody;

import java.util.HashMap;
import java.util.List;

/**
 * 服务版本条件接口
 *
 * @author sctang2
 * @create 2017-12-10 16:25
 **/
public interface IServiceVersionCondition {
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
     * 查询版本列表
     *
     * @param serviceIds
     * @return
     */
    List<ServiceVersion> findList(List<String> serviceIds);

    /**
     * 新增版本
     *
     * @param body
     * @return
     */
    ServiceVersion add(AddServiceVersionRequestBody body);

    /**
     * 通过id查询服务配置列表
     *
     * @param id
     * @return
     */
    ServiceVersion findServiceConfigListById(String id);

    /**
     * 根据id删除版本
     *
     * @param id
     * @return
     */
    int deleteById(String id);

    /**
     * 根据id查询版本
     *
     * @param id
     * @return
     */
    ServiceVersion findById(String id);

    /**
     * 根据id更新版本
     *
     * @param id
     * @param body
     * @return
     */
    ServiceVersion updateById(String id, EditServiceVersionRequestBody body);

    /**
     * 根据版本和服务id查询服务版本
     *
     * @param version
     * @param serviceId
     * @return
     */
    ServiceVersion find(String version, String serviceId);

    /**
     * 通过id查询服务版本、服务、服务组、项目信息
     *
     * @param id
     * @return
     */
    ServiceVersion findVersionJoinServiceJoinGroupJoinProjectById(String id);
}
