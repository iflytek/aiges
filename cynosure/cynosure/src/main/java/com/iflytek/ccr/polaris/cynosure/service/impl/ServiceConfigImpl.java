package com.iflytek.ccr.polaris.cynosure.service.impl;

import com.alibaba.fastjson.JSONArray;
import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.companionservice.ConfigCenter;
import com.iflytek.ccr.polaris.cynosure.companionservice.GrayConfigCenter;
import com.iflytek.ccr.polaris.cynosure.companionservice.ServiceDiscoveryCenter;
import com.iflytek.ccr.polaris.cynosure.companionservice.domain.PushDetailResult;
import com.iflytek.ccr.polaris.cynosure.companionservice.domain.PushResult;
import com.iflytek.ccr.polaris.cynosure.customdomain.FileContent;
import com.iflytek.ccr.polaris.cynosure.customdomain.SearchCondition;
import com.iflytek.ccr.polaris.cynosure.dbcondition.*;
import com.iflytek.ccr.polaris.cynosure.dbtransactional.QuickStartTransactional;
import com.iflytek.ccr.polaris.cynosure.dbtransactional.ServiceConfigTransactional;
import com.iflytek.ccr.polaris.cynosure.domain.*;
import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.exception.GlobalExceptionUtil;
import com.iflytek.ccr.polaris.cynosure.request.IdRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.IdsRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.graygroup.AddGrayConfigRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.graygroup.ServiceGrayConfigHistoryListRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceconfig.*;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.serviceconfig.*;
import com.iflytek.ccr.polaris.cynosure.service.ILastestSearchService;
import com.iflytek.ccr.polaris.cynosure.service.IServiceConfig;
import com.iflytek.ccr.polaris.cynosure.util.PagingUtil;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.io.*;
import java.util.*;
import java.util.stream.Collectors;

/**
 * 配置服务业务逻辑接口实现
 *
 * @author sctang2
 * @create 2017-11-21 11:46
 **/
@Service
public class ServiceConfigImpl extends BaseService implements IServiceConfig {
    private final EasyLogger logger = EasyLoggerFactory.getInstance(ServiceConfigImpl.class);

    @Autowired
    private ILastestSearchService lastestSearchServiceImpl;

    @Autowired
    private IServiceConfigCondition serviceConfigConditionImpl;

    @Autowired
    private IRegionCondition regionConditionImpl;

    @Autowired
    private ConfigCenter configCenter;

    @Autowired
    private GrayConfigCenter grayConfigCenter;

    @Autowired
    private ServiceDiscoveryCenter serviceDiscoveryCenter;

    @Autowired
    private ServiceConfigTransactional serviceConfigTransactional;

    @Autowired
    private IServiceConfigHistoryCondition serviceConfigHistoryConditionImpl;

    @Autowired
    private IServiceConfigPushHistoryCondition servicePushHistoryConditionImpl;

    @Autowired
    private IServiceConfigPushFeedbackCondition serviceConfigPushFeedbackConditionImpl;

    @Autowired
    private IGrayGroupCondition grayGroupConditionImpl;

    @Autowired
    private QuickStartTransactional quickStartTransactional;

    @Autowired
    private IProjectCondition projectConditionImpl;

    @Autowired
    private IClusterCondition clusterConditionImpl;

    @Autowired
    private IServiceCondition serviceConditionImpl;

    @Autowired
    private IServiceVersionCondition serviceVersionConditionImpl;


    @Override
    public Response<QueryPagingListResponseBody> findLastestList(QueryServiceConfigRequestBody body) {
        QueryPagingListResponseBody result;
        String projectName = body.getProject();
        String clusterName = body.getCluster();
        String serviceName = body.getService();
        String versionName = body.getVersion();
        String grayName = body.getGray();

        if (StringUtils.isBlank(grayName)) {
            //查询最近的搜索
            SearchCondition searchCondition = this.lastestSearchServiceImpl.find(projectName, clusterName, serviceName, versionName);
            projectName = searchCondition.getProject();
            clusterName = searchCondition.getCluster();
            serviceName = searchCondition.getService();
            versionName = searchCondition.getVersion();
            if (StringUtils.isBlank(projectName) || StringUtils.isBlank(clusterName) || StringUtils.isBlank(serviceName) || StringUtils.isBlank(versionName)) {
                result = PagingUtil.createResult(body, 0);
                return new Response<>(result);
            }

            //创建分页查询条件
            HashMap<String, Object> map = PagingUtil.createCondition(body);
            map.put("projectName", projectName);
            map.put("clusterName", clusterName);
            map.put("serviceName", serviceName);
            map.put("versionName", versionName);
            map.put("grayId", "0");
            //查询总数
            int totalCount = this.serviceConfigConditionImpl.findTotalCount(map);
            //创建分页结果
            result = PagingUtil.createResult(body, totalCount);

            //保存最近的搜索
            String condition = this.lastestSearchServiceImpl.saveLastestSearch(projectName, clusterName, serviceName, versionName);
            result.setCondition(condition);
            if (0 == totalCount) {
                return new Response<>(result);
            }

            //查询列表
            List<ServiceConfigDetailResponseBody> list = new ArrayList<>();
            Optional<List<ServiceConfig>> serviceConfigList = Optional.ofNullable(this.serviceConfigConditionImpl.findList(map));
            serviceConfigList.ifPresent(x -> {
                x.forEach(y -> {
                    //创建服务配置结果
                    ServiceConfigDetailResponseBody serviceConfigDetail = this.createServiceConfigResult(y);
                    list.add(serviceConfigDetail);
                });
            });
            result.setList(list);
        } else {

            //查询最近的搜索
            SearchCondition searchCondition = this.lastestSearchServiceImpl.find(projectName, clusterName, serviceName, versionName, grayName);
            projectName = searchCondition.getProject();
            clusterName = searchCondition.getCluster();
            serviceName = searchCondition.getService();
            versionName = searchCondition.getVersion();
            grayName = searchCondition.getGray();
            if (StringUtils.isBlank(projectName) || StringUtils.isBlank(clusterName) || StringUtils.isBlank(serviceName) || StringUtils.isBlank(versionName) || StringUtils.isBlank(grayName)) {
                result = PagingUtil.createResult(body, 0);
                return new Response<>(result);
            }

            //创建分页查询条件
            HashMap<String, Object> map = PagingUtil.createCondition(body);
            map.put("projectName", projectName);
            map.put("clusterName", clusterName);
            map.put("serviceName", serviceName);
            map.put("versionName", versionName);
            map.put("grayName", grayName);
            //查询总数
            int totalCount = this.serviceConfigConditionImpl.findGrayTotalCount(map);
            //创建分页结果
            result = PagingUtil.createResult(body, totalCount);

            //保存最近的搜索
            String condition = this.lastestSearchServiceImpl.saveLastestSearch(projectName, clusterName, serviceName, versionName, grayName);
            result.setCondition(condition);
            if (0 == totalCount) {
                return new Response<>(result);
            }

            //查询列表
            List<ServiceConfigDetailResponseBody> list = new ArrayList<>();
            Optional<List<ServiceConfig>> serviceConfigList = Optional.ofNullable(this.serviceConfigConditionImpl.findGrayList(map));
            serviceConfigList.ifPresent(x -> {
                x.forEach(y -> {
                    //创建服务配置结果
                    ServiceConfigDetailResponseBody serviceConfigDetail = this.createServiceConfigResult(y);
                    list.add(serviceConfigDetail);
                });
            });
            result.setList(list);
        }
        return new Response<>(result);
    }

    @Override
    public Response<ServiceConfigDetailResponseBody> edit(EditServiceConfigRequestBody body) {
        //通过id查询服务配置
        String id = body.getId();
        ServiceConfig serviceConfig = this.serviceConfigConditionImpl.findById(id);
        if (null == serviceConfig) {
            //不存在该配置
            return new Response<>(SystemErrCode.ERRCODE_SERVICE_CONFIG_NOT_EXISTS, SystemErrCode.ERRMSG_SERVICE_CONFIG_NOT_EXISTS);
        }
        //校验文件字节长度
        byte[] contentByte = null;
        try{
            contentByte = body.getContent().getBytes("utf-8");
        }catch (Exception e){
            e.printStackTrace();
        }

        if (contentByte.length > (1024 * 1024)) {
            return new Response<>(SystemErrCode.ERRCODE_SERVICE_CONFIG_CONTENT_BYTE_MAX_LENGTH, SystemErrCode.ERRMSG_SERVICE_CONFIG_CONTENT_BYTE_MAX_LENGTH);
        }
        //修改配置信息
        ServiceConfig updateServiceConfig = this.serviceConfigConditionImpl.updateById(id, body);

        //创建服务配置结果
        updateServiceConfig.setCreateTime(serviceConfig.getCreateTime());
        updateServiceConfig.setName(serviceConfig.getName());
        updateServiceConfig.setPath(serviceConfig.getPath());
        updateServiceConfig.setVersionId(serviceConfig.getVersionId());
        updateServiceConfig.setGrayId(serviceConfig.getGrayId());
        ServiceConfigDetailResponseBody result = this.createServiceConfigResult(updateServiceConfig);
        return new Response<>(result);
    }

    @Override
    public Response<ServiceConfigDetailResponseBody> find(String id) {
        //通过id查询服务配置
        ServiceConfig serviceConfig = this.serviceConfigConditionImpl.findById(id);
        if (null == serviceConfig) {
            //不存在该配置
            return new Response<>(SystemErrCode.ERRCODE_SERVICE_CONFIG_NOT_EXISTS, SystemErrCode.ERRMSG_SERVICE_CONFIG_NOT_EXISTS);
        }

        //创建服务配置结果
        ServiceConfigDetailResponseBody result = this.createServiceConfigResult(serviceConfig);
        return new Response<>(result);
    }

    @Override
    public Response<PushServiceConfigResponseBody> push(PushServiceConfigRequestBody body) {
        //通过id查询配置、版本、服务、集群、项目信息
        ServiceConfig serviceConfig = this.serviceConfigConditionImpl.findConfigJoinVersionJoinServiceJoinClusterJoinProjectById(body.getId());

        //不存在该配置
        if (null == serviceConfig) {
            return new Response<>(SystemErrCode.ERRCODE_SERVICE_CONFIG_NOT_EXISTS, SystemErrCode.ERRMSG_SERVICE_CONFIG_NOT_EXISTS);
        }

        //通过区域ids查询区域
        List<String> ids = body.getRegionIds();
        List<Region> regionList = this.regionConditionImpl.findListByIds(ids);

        //不存在该推送区域
        if (null == regionList || regionList.isEmpty()) {
            return new Response<>(SystemErrCode.ERRCODE_REGION_NOT_EXISTS, SystemErrCode.ERRMSG_REGION_NOT_EXISTS);
        }

        //得到要推送的数据信息
        String path = serviceConfig.getPath();
        byte[] content = serviceConfig.getContent();
        String fileName = serviceConfig.getName();
        String grayId = serviceConfig.getGrayId();

        //判断是否是灰度配置文件，灰度和非灰度的推送路径不同
        if (!grayId.equals("0")) {
            //灰度配置文件推送
            GrayGroup grayGroup = this.grayGroupConditionImpl.findById(grayId);
            if (null == grayGroup) {
                return new Response<>(SystemErrCode.ERRCODE_GRAY_GROUP_NOT_EXISTS, SystemErrCode.ERRMSG_GRAY_GROUP_NOT_EXISTS);
            }
            //获取推送实例的地址
            String grayServers = grayGroup.getContent();

            //截取最后一个/的路径
            String grayPath = StringUtils.substringBeforeLast(path, "/");

            //灰度组推送
            PushResult cacheCenterPushResult = this.grayConfigCenter.pushByOneToMany(grayPath, fileName, grayId, grayServers, content, regionList);

            //新增服务配置和推送历史
            PushServiceConfigResponseBody result = this.serviceConfigTransactional.addServiceConfigAndPushHistory(serviceConfig, cacheCenterPushResult);

            //检查推送失败的区域
            List<PushDetailResult> datas = cacheCenterPushResult.getData();
            List<String> failAreaNames = getFailAreas(datas);
            if (failAreaNames.size()!=0){
                Response response = new Response<>(result);
                response.setCode(-1);
                response.setMessage("共有"+failAreaNames.size()+"个区域推送失败："+JSONArray.toJSONString(failAreaNames));
                return response;
            }
            //返回结果
            return new Response<>(result);
        } else {
            //非灰度文件推送
            PushResult cacheCenterPushResult = this.configCenter.pushByOneToMany(path, content, regionList);

            //新增服务配置和推送历史
            PushServiceConfigResponseBody result = this.serviceConfigTransactional.addServiceConfigAndPushHistory(serviceConfig, cacheCenterPushResult);

            //检查推送失败的区域
            List<PushDetailResult> datas = cacheCenterPushResult.getData();
            List<String> failAreaNames = getFailAreas(datas);
            if (failAreaNames.size()!=0){
                Response response = new Response<>(result);
                response.setCode(-1);
                response.setMessage("共有"+failAreaNames.size()+"个区域推送失败："+JSONArray.toJSONString(failAreaNames));
                return response;
            }
            //返回结果
            return new Response<>(result);
        }
    }

    @Override
    public Response<PushServiceConfigResponseBody> batchPush(BatchPushServiceConfigRequestBody body) {
        //通过ids查询配置、版本、服务、集群、项目信息
        List<String> configIds = body.getIds();
        List<ServiceConfig> serviceConfigList = this.serviceConfigConditionImpl.findConfigJoinVersionJoinServiceJoinClusterJoinProjectByIds(configIds);
        if (null == serviceConfigList || serviceConfigList.isEmpty()) {
            //不存在该配置
            return new Response<>(SystemErrCode.ERRCODE_SERVICE_CONFIG_NOT_EXISTS, SystemErrCode.ERRMSG_SERVICE_CONFIG_NOT_EXISTS);
        }

        //通过区域ids查询区域
        List<String> ids = body.getRegionIds();
        List<Region> regionList = this.regionConditionImpl.findListByIds(ids);
        if (null == regionList || regionList.isEmpty()) {
            //不存在该集群
            return new Response<>(SystemErrCode.ERRCODE_REGION_NOT_EXISTS, SystemErrCode.ERRMSG_REGION_NOT_EXISTS);
        }

        //通过多对多推送到缓存中心
        List<String> pathList = serviceConfigList.stream().map(x -> x.getPath()).collect(Collectors.toList());
        List<byte[]> contentList = serviceConfigList.stream().map(x -> x.getContent()).collect(Collectors.toList());
        List<String> fileNameList = serviceConfigList.stream().map(x -> x.getName()).collect(Collectors.toList());

        //判断是否是灰度配置文件
        List<String> grayIdList = serviceConfigList.stream().map(x -> x.getGrayId()).collect(Collectors.toList());
        if (!"0".equals(grayIdList.get(0))) {

            GrayGroup grayGroup = this.grayGroupConditionImpl.findById(grayIdList.get(0));
            if (null == grayGroup) {
                return new Response<>(SystemErrCode.ERRCODE_GRAY_GROUP_NOT_EXISTS, SystemErrCode.ERRMSG_GRAY_GROUP_NOT_EXISTS);
            }
            String grayServers = grayGroup.getContent();

            List<String> grayPathList = new ArrayList<>();
            for (String path : pathList) {
                String grayPath = StringUtils.substringBeforeLast(path, "/");
                grayPathList.add(grayPath);
            }
            PushResult cacheCenterPushResult = this.grayConfigCenter.pushByManyToMany(grayPathList, fileNameList, grayIdList.get(0), grayServers, contentList, regionList);

            //新增服务配置列表和推送历史
            PushServiceConfigResponseBody result = this.serviceConfigTransactional.addServiceConfigsAndPushHistory(serviceConfigList, cacheCenterPushResult);

            //检验推送是否有失败的区域
            List<PushDetailResult> datas = cacheCenterPushResult.getData();
            List<String> failAreaNames = getFailAreas(datas);
            if (failAreaNames.size()!=0){
                Response response = new Response<>(result);
                response.setCode(-1);
                response.setMessage("共有"+failAreaNames.size()+"个区域推送失败："+JSONArray.toJSONString(failAreaNames));
                return response;
            }
            return new Response<>(result);
        } else {
            PushResult cacheCenterPushResult = this.configCenter.pushByManyToMany(pathList, contentList, regionList);

            //新增服务配置列表和推送历史
            PushServiceConfigResponseBody result = this.serviceConfigTransactional.addServiceConfigsAndPushHistory(serviceConfigList, cacheCenterPushResult);

            //检验推送失败的区域
            List<PushDetailResult> datas = cacheCenterPushResult.getData();
            List<String> failAreaNames = getFailAreas(datas);
            if (failAreaNames.size()!=0){
                Response response = new Response<>(result);
                response.setCode(-1);
                response.setMessage("共有"+failAreaNames.size()+"个区域推送失败："+JSONArray.toJSONString(failAreaNames));
                return response;
            }
            return new Response<>(result);
        }
    }

    //获得推送失败区域的信息
    private List<String> getFailAreas(List<PushDetailResult> datas){
        List<String> failAreaNames = new ArrayList<>();
        for (PushDetailResult data: datas) {
            int successed = data.getSuccessed();
            if (successed==-1){
                failAreaNames.add(data.getName());
            }
        }
        return failAreaNames;
    }

    @Override
    public Response<String> delete(IdRequestBody body) {
        //通过id查询服务配置
        String id = body.getId();
        ServiceConfig serviceConfig = this.serviceConfigConditionImpl.findById(id);
        if (null == serviceConfig) {
            //不存在该配置
            return new Response<>(SystemErrCode.ERRCODE_SERVICE_CONFIG_NOT_EXISTS, SystemErrCode.ERRMSG_SERVICE_CONFIG_NOT_EXISTS);
        }

        //删除配置和历史
        this.serviceConfigTransactional.deleteConfigAndHistory(id);

        //查询集群列表
        List<Region> regionList = this.regionConditionImpl.findList(null);
        if (null == regionList || regionList.isEmpty()) {
            return new Response<>(null);
        }

        //删除缓存中心配置,灰度和非灰度上报的路径不同
        String grayId = serviceConfig.getGrayId();
        if ("0".equals(grayId)) {
            String path = serviceConfig.getPath();
            this.configCenter.deleteConf(path, regionList);
        } else {
            //拼接灰度配置路径
            String path = StringUtils.substringBeforeLast(serviceConfig.getPath(), "/") + "/gray" + "/" + grayId + "/" + serviceConfig.getName();
            logger.info(path);
            this.grayConfigCenter.deleteGrayConfig(path, grayId, regionList);
        }

        return new Response<>(null);
    }

    @Override
    public Response<String> batchDelete(IdsRequestBody body) {
        //通过ids查询服务配置列表
        List<String> ids = body.getIds();
        List<ServiceConfig> serviceConfigList = this.serviceConfigConditionImpl.findByIds(ids);
        if (null == serviceConfigList || serviceConfigList.isEmpty()) {
            //不存在该配置
            return new Response<>(SystemErrCode.ERRCODE_SERVICE_CONFIG_NOT_EXISTS, SystemErrCode.ERRMSG_SERVICE_CONFIG_NOT_EXISTS);
        }

        //批量删除配置和历史
        this.serviceConfigTransactional.batchDeleteConfigAndHistory(ids);

        //查询集群列表
        List<Region> regionList = this.regionConditionImpl.findList(null);
        if (null == regionList || regionList.isEmpty()) {
            return new Response<>(null);
        }

        //删除缓存中心配置,区分灰度和非灰度配置
        ServiceConfig serviceConfigFist = serviceConfigList.get(0);
        String grayId = serviceConfigFist.getGrayId();
        if ("0".equals(grayId)) {
            for (ServiceConfig serviceConfig : serviceConfigList) {
                String path = serviceConfig.getPath();
                this.configCenter.deleteConf(path, regionList);
            }
        } else {
            for (ServiceConfig serviceConfig : serviceConfigList) {
                String path = StringUtils.substringBeforeLast(serviceConfig.getPath(), "/") + "/gray" + "/" + grayId + "/" + serviceConfig.getName();
                this.grayConfigCenter.deleteGrayConfig(path, grayId, regionList);
            }
        }
        return new Response<>(null);
    }

    @Override
    public Response<QueryPagingListResponseBody> findServiceConfigHistoryList(ServiceConfigHistoryListRequestBody body) {

        //查询项目Id
        Project project = this.projectConditionImpl.findByName(body.getProject());
        if (project == null){
            return new Response<>(SystemErrCode.ERRCODE_PROJECT_NOT_EXISTS, SystemErrCode.ERRMSG_PROJECT_NOT_EXISTS);
        }

        //根据项目Id和集群名称查询
        Cluster cluster = this.clusterConditionImpl.find(project.getId(), body.getCluster());
        if (cluster == null){
            return new Response<>(SystemErrCode.ERRCODE_CLUSTER_NOT_EXISTS, SystemErrCode.ERRMSG_CLUSTER_NOT_EXISTS);
        }

        //根据集群Id和服务名字查询服务
        com.iflytek.ccr.polaris.cynosure.domain.Service service = this.serviceConditionImpl.find(body.getService(), cluster.getId());
        if (service == null){
            return new Response<>(SystemErrCode.ERRCODE_SERVICE_NOT_EXISTS, SystemErrCode.ERRMSG_SERVICE_NOT_EXISTS);
        }

        //根据服务Id和版本名称查询版本
        ServiceVersion version = this.serviceVersionConditionImpl.find(body.getVersion(), service.getId());
        if (version == null){
            return new Response<>(SystemErrCode.ERRCODE_SERVICE_VERSION_NOT_EXISTS, SystemErrCode.ERRMSG_SERVICE_VERSION_NOT_EXISTS);
        }

        //根据版本Id查询配置列表
        List<ServiceConfig> configs = this.serviceConfigConditionImpl.findConfigsByVersionId(version.getId(), "0");

        //得到配置Id集合
        List<String> configIds = configs.stream().map(x -> x.getId()).distinct().collect(Collectors.toList());

        //不存在数据
        if (configIds.isEmpty()||configIds==null){
            return new Response<>(null);
        }

        //创建分页查询条件
        HashMap<String, Object> map = PagingUtil.createCondition(body);
        map.put("configIds", configIds);


        //查询总数
        int totalCount = this.serviceConfigHistoryConditionImpl.findTotalCount(map);


        //创建分页结果
        QueryPagingListResponseBody result = PagingUtil.createResult(body, totalCount);
        if (0 == totalCount) {
            return new Response<>(result);
        }
        //查询列表
        List<ServiceConfigHistory> serviceConfigHistoryList = this.serviceConfigHistoryConditionImpl.findList(map);

        //创建服务配置历史列表
        List<QueryServiceConfigHistoryResponseBody> list = this.createServiceConfigHistoryList(serviceConfigHistoryList);
        result.setList(list);
        return new Response<>(result);
    }

    /**
     * 创建服务配置历史列表
     *
     * @param serviceConfigHistoryList
     * @return
     */
    private List<QueryServiceConfigHistoryResponseBody> createServiceConfigHistoryList(List<ServiceConfigHistory> serviceConfigHistoryList){
        List<QueryServiceConfigHistoryResponseBody> list = new ArrayList<>();
        //获取推送版本号列表
        List<String> pushVersionList = serviceConfigHistoryList.stream().map(x -> x.getPushVersion()).distinct().collect(Collectors.toList());
        pushVersionList.forEach(x -> {
            QueryServiceConfigHistoryResponseBody serviceConfigHistory = new QueryServiceConfigHistoryResponseBody();
            serviceConfigHistory.setPushVersion(x);
            List<ServiceConfigHistoryResponseBody> histories = new ArrayList<>();
            serviceConfigHistoryList.forEach(y -> {
                String pushVersion = y.getPushVersion();
                if (x.equals(pushVersion)) {
                    ServiceConfigHistoryResponseBody history = new ServiceConfigHistoryResponseBody();
                    history.setId(y.getId());
                    history.setName(y.getServiceConfig().getName());
                    history.setConfigId(y.getConfigId());
                    try {
                        history.setContent(new String(y.getContent(), "UTF-8"));
                    } catch (UnsupportedEncodingException e) {
                        GlobalExceptionUtil.log(e);
                    }
                    history.setCreateTime(y.getCreateTime());
                    serviceConfigHistory.setCreateTime(y.getCreateTime());
                    history.setDesc(y.getDescription());
                    history.setPushVersion(y.getPushVersion());
                    histories.add(history);
                }
            });
            serviceConfigHistory.setHistories(histories);
            list.add(serviceConfigHistory);
        });
        return list;
    }

    @Override
    public Response<List<ServiceConfigDetailResponseBody>> rollback(IdsRequestBody body) {
        //通过id查询服务配置历史
        List<String> ids = body.getIds();
        List<ServiceConfigHistory> serviceConfigHistoryList = this.serviceConfigHistoryConditionImpl.findByIds(ids);
        if (null == serviceConfigHistoryList || serviceConfigHistoryList.isEmpty()) {
            //不存在该配置历史
            return new Response<>(SystemErrCode.ERRCODE_SERVICE_CONFIG_HISTORY_NOT_EXISTS, SystemErrCode.ERRMSG_SERVICE_CONFIG_HISTORY_NOT_EXISTS);
        }

        //通过id查询服务配置
        List<String> configIds = serviceConfigHistoryList.stream().map(x -> x.getConfigId()).collect(Collectors.toList());
        List<ServiceConfig> serviceConfigList = this.serviceConfigConditionImpl.findByIds(configIds);
        if (null == serviceConfigList || serviceConfigList.isEmpty()) {
            //不存在该配置
            return new Response<>(SystemErrCode.ERRCODE_SERVICE_CONFIG_NOT_EXISTS, SystemErrCode.ERRMSG_SERVICE_CONFIG_NOT_EXISTS);
        }

        //批量更新服务配置
        List<ServiceConfig> newServiceConfigList = this.serviceConfigConditionImpl.batchUpdate(serviceConfigHistoryList, serviceConfigList);

        //创建服务配置结果
        List<ServiceConfigDetailResponseBody> result = new ArrayList<>();
        for (ServiceConfig serviceConfig : newServiceConfigList) {
            ServiceConfigDetailResponseBody item = this.createServiceConfigResult(serviceConfig);
            result.add(item);
        }
        return new Response<>(result);
    }

    @Override
    public Response<String> feedback(ServiceConfigFeedBackRequestBody body) {
        //通过id查询服务 推送历史
        String pushId = body.getPushId();
        ServiceConfigPushHistory servicePushHistory = this.servicePushHistoryConditionImpl.findById(pushId);
        if (null == servicePushHistory) {
            //不存在该轨迹
            return new Response<>(SystemErrCode.ERRCODE_TRACK_NOT_EXISTS, SystemErrCode.ERRMSG_TRACK_NOT_EXISTS);
        }

        //新增服务推送反馈
        this.serviceConfigPushFeedbackConditionImpl.add(body);
        return new Response<>(null);
    }

    @Override
    public List<ServiceConfig> findListByIds(List<String> ids) {
        if (ids == null && ids.isEmpty()) {
            return null;
        }
        return serviceConfigConditionImpl.findListByIds(ids);
    }


    @Override
    public Response<GrayConfigListDetailResponseBody> addGrayConfig(AddGrayConfigRequestBody body, List<FileContent> fileContentList) throws UnsupportedEncodingException {
        return this.quickStartTransactional.addGrayConfig(body, fileContentList);
    }

    @Override
    public Response<QueryPagingListResponseBody> findServiceGrayConfigHistoryList(ServiceGrayConfigHistoryListRequestBody body) {
        //查询项目Id
        Project project = this.projectConditionImpl.findByName(body.getProject());
        if (project == null){
            return new Response<>(SystemErrCode.ERRCODE_PROJECT_NOT_EXISTS, SystemErrCode.ERRMSG_PROJECT_NOT_EXISTS);
        }

        //根据项目Id和集群名称查询
        Cluster cluster = this.clusterConditionImpl.find(project.getId(), body.getCluster());
        if (cluster == null){
            return new Response<>(SystemErrCode.ERRCODE_CLUSTER_NOT_EXISTS, SystemErrCode.ERRMSG_CLUSTER_NOT_EXISTS);
        }

        //根据集群Id和服务名字查询服务
        com.iflytek.ccr.polaris.cynosure.domain.Service service = this.serviceConditionImpl.find(body.getService(), cluster.getId());
        if (service == null){
            return new Response<>(SystemErrCode.ERRCODE_SERVICE_NOT_EXISTS, SystemErrCode.ERRMSG_SERVICE_NOT_EXISTS);
        }

        //根据服务Id和版本名称查询版本
        ServiceVersion version = this.serviceVersionConditionImpl.find(body.getVersion(), service.getId());
        if (version == null){
            return new Response<>(SystemErrCode.ERRCODE_SERVICE_VERSION_NOT_EXISTS, SystemErrCode.ERRMSG_SERVICE_VERSION_NOT_EXISTS);
        }

        //根据灰度组名字查询灰度组
        GrayGroup grayGroup = this.grayGroupConditionImpl.find(body.getGray(), version.getId());
        if (grayGroup == null){
            return new Response<>(SystemErrCode.ERRCODE_GRAY_GROUP_NOT_EXISTS, SystemErrCode.ERRMSG_GRAY_GROUP_NOT_EXISTS);
        }

        //根据版本Id和灰度Id查询配置文件列表
        List<ServiceConfig> configs = this.serviceConfigConditionImpl.findConfigsByVersionId(version.getId(), grayGroup.getId());

        //得到配置Id集合
        List<String> configIds = configs.stream().map(x -> x.getId()).distinct().collect(Collectors.toList());

        //不存在数据
        if (configIds.isEmpty()||configIds==null){
            return new Response<>(null);
        }

        //创建分页查询条件
        HashMap<String, Object> map = PagingUtil.createCondition(body);
        map.put("configIds", configIds);


        //查询总数
        int totalCount = this.serviceConfigHistoryConditionImpl.findTotalCount(map);


        //创建分页结果
        QueryPagingListResponseBody result = PagingUtil.createResult(body, totalCount);
        if (0 == totalCount) {
            return new Response<>(result);
        }
        //查询列表
        List<ServiceConfigHistory> serviceConfigHistoryList = this.serviceConfigHistoryConditionImpl.findList(map);

        //创建服务配置历史列表
        List<QueryServiceConfigHistoryResponseBody> list = this.createServiceConfigHistoryList(serviceConfigHistoryList);
        result.setList(list);
        return new Response<>(result);
    }


    @Override
    public Response<DownloadFile> download(DownloadServiceConfigRequestBody body) {
        String id = body.getId();

        //通过id查询服务配置
        ServiceConfig serviceConfig = this.serviceConfigConditionImpl.findById(id);
        if (null == serviceConfig) {
            return null;
        }
        //下载文件
        String fileName=serviceConfig.getName();
        byte[] content = serviceConfig.getContent();
        String contentStr = null;
        try{
            contentStr = new String(content, "utf-8");
        }catch (Exception e){
            e.printStackTrace();
        }

        DownloadFile downloadFile = new DownloadFile();
        downloadFile.setFileName(fileName);
        downloadFile.setContent(contentStr);

        return new Response<>(downloadFile);
    }



}