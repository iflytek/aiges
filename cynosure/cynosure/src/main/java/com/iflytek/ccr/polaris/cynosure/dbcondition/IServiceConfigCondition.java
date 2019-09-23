package com.iflytek.ccr.polaris.cynosure.dbcondition;

import com.iflytek.ccr.polaris.cynosure.customdomain.FileContent;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceConfig;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceConfigHistory;
import com.iflytek.ccr.polaris.cynosure.request.graygroup.AddGrayConfigRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.graygroup.AddGrayGroupRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.quickstart.AddServiceConfigRequestBodyByQuickStart;
import com.iflytek.ccr.polaris.cynosure.request.serviceconfig.AddServiceConfigRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceconfig.EditServiceConfigRequestBody;

import java.io.UnsupportedEncodingException;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * 服务配置条件接口
 *
 * @author sctang2
 * @create 2017-12-10 17:32
 **/
public interface IServiceConfigCondition {
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
     * 批量新增服务配置（新增服务版本,无配置文件上传）
     *
     * @param body
     * @param addServiceConfig
     * @return
     */
    List<ServiceConfig> batchVersionAdd(AddServiceConfigRequestBody body, AddServiceConfigRequestBodyByQuickStart addServiceConfig, List<ServiceConfig> serviceConfigList);

    /**
     * 批量新增服务配置（新增服务版本含配置文件上传）
     *
     * @param body
     * @param addServiceConfig
     * @return
     */
    List<ServiceConfig> batchVersionAddAndFile(AddServiceConfigRequestBody body, AddServiceConfigRequestBodyByQuickStart addServiceConfig, List<ServiceConfig> serviceConfigList, List<FileContent> fileContentList);

    /**
     * 批量新增服务配置
     *
     * @param body
     * @param addServiceConfig
     * @param fileContentList
     * @return
     */
    List<ServiceConfig> batchAdd(AddServiceConfigRequestBody body, AddServiceConfigRequestBodyByQuickStart addServiceConfig, List<FileContent> fileContentList);

    /**
     * 批量更新服务配置
     *
     * @param body
     * @param fileContentList
     * @param serviceConfigList
     * @return
     */
    List<ServiceConfig> batchUpdate(AddServiceConfigRequestBodyByQuickStart body, List<FileContent> fileContentList, List<ServiceConfig> serviceConfigList);

    /**
     * 批量更新服务配置
     *
     * @param serviceConfigHistoryList
     * @param serviceConfigList
     * @return
     */
    List<ServiceConfig> batchUpdate(List<ServiceConfigHistory> serviceConfigHistoryList, List<ServiceConfig> serviceConfigList);

    /**
     * 新增服务配置
     *
     * @param body
     * @param fileContent
     * @return
     */
    ServiceConfig add(AddServiceConfigRequestBody body, FileContent fileContent);

    /**
     * 通过id查询服务配置
     *
     * @param id
     * @return
     */
    ServiceConfig findById(String id);

    /**
     * 通过ids查询服务配置列表
     *
     * @param ids
     * @return
     */
    List<ServiceConfig> findByIds(List<String> ids);

    /**
     * 通过id删除服务配置
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
     * 根据id更新服务配置
     *
     * @param id
     * @param body
     * @return
     */
    ServiceConfig updateById(String id, EditServiceConfigRequestBody body);

    /**
     * 通过版本id，配置名称，灰度组id，列表查询服务配置
     *
     * @param versionId
     * @param names
     * @return
     */
    List<ServiceConfig> find(String versionId, List<String> names, String grayId);

    /**
     * 通过版本Id查询配置文件集合
     * @param versionId
     * @return
     */
    List<ServiceConfig> findConfigsByVersionId(String versionId, String grayId);

    /**
     * 通过配置名称，版本id查询服务配置
     *
     * @param name
     * @param versionId
     * @return
     */
    ServiceConfig find(String name, String versionId);

    /**
     * 查询最新版本的配置信息
     *
     * @param serviceId
     * @return
     */
    List<ServiceConfig> findNewList(String serviceId);

    /**
     * 通过id查询配置、版本、服务、集群、项目信息
     *
     * @param id
     * @return
     */
    ServiceConfig findConfigJoinVersionJoinServiceJoinClusterJoinProjectById(String id);

    /**
     * 通过ids查询配置、版本、服务、集群、项目信息
     *
     * @param ids
     * @return
     */
    List<ServiceConfig> findConfigJoinVersionJoinServiceJoinClusterJoinProjectByIds(List<String> ids);

    /**
     * 通过ids查询服务配置列表信息
     *
     * @param ids
     * @return
     */
    List<ServiceConfig> findListByIds(List<String> ids);

    /**
     * 批量新增服务配置（新增灰度组,无配置文件上传）
     *
     * @param body
     * @param grayId
     * @param serviceConfigList
     * @return
     */
    List<ServiceConfig> batchGrayGroupAdd(AddGrayGroupRequestBody body, String grayId, List<ServiceConfig> serviceConfigList);

    /**
     * 批量新增服务配置（新增灰度组,含配置文件上传）
     *
     * @param body
     * @param grayId
     * @param serviceConfigList
     * @param fileContentList
     * @return
     */
    List<ServiceConfig> batchGrayGroupAndFileAdd(AddGrayGroupRequestBody body, String grayId, List<ServiceConfig> serviceConfigList, List<FileContent> fileContentList);

    /**
     * 通过ids查询服务配置列表信息
     *
     * @param ids
     * @return
     */
    List<ServiceConfig> findListByGrayIds(List<String> ids);

    /**
     * 通过id查询服务配置列表信息
     *
     * @param id
     * @return
     */
    List<ServiceConfig> findListByGrayId(String id);

    /**
     * 通过灰度组id删除灰度服务配置
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
     * 通过配置名称，版本id，灰度组id校验唯一性
     *
     * @param name
     * @param versionId
     * @return
     */
    ServiceConfig findOnlyConfig(String name, String versionId, String grayId);

    /**
     * 批量新增灰度配置（拖拽配置文件上传）
     *
     * @param versionId
     * @param body
     * @param fileContentList
     * @return
     */
    List<ServiceConfig> batchAddGrayFile(String versionId, AddGrayConfigRequestBody body, List<FileContent> fileContentList) throws UnsupportedEncodingException;


    /**
     * 批量更新灰度配置
     *
     * @param body
     * @param fileContentList
     * @param serviceConfigList
     * @return
     */
    List<ServiceConfig> batchUpdateGrayConfig(AddGrayConfigRequestBody body, List<FileContent> fileContentList, List<ServiceConfig> serviceConfigList);

    /**
     * 复制配置文件
     * @param serviceConfigs
     * @param versionId
     * @param oldGrayId2NewGrayId
     * @return
     */
    int copyConfigs1(List<ServiceConfig> serviceConfigs, String versionId, Map<String, String> oldGrayId2NewGrayId, String path);
}
