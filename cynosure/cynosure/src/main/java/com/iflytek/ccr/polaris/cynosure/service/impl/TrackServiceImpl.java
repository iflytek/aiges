package com.iflytek.ccr.polaris.cynosure.service.impl;

import com.alibaba.fastjson.JSON;
import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.customdomain.SearchCondition;
import com.iflytek.ccr.polaris.cynosure.dbcondition.*;
import com.iflytek.ccr.polaris.cynosure.dbtransactional.ServiceConfigPushTransactional;
import com.iflytek.ccr.polaris.cynosure.dbtransactional.ServiceDiscoveryPushTransactional;
import com.iflytek.ccr.polaris.cynosure.domain.*;
import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.request.BaseRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.IdRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.IdsRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.track.QueryTrackDetailRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.track.QueryTrackRequestBody;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.quickstart.QueryClusterResponseBodyByQuickStart;
import com.iflytek.ccr.polaris.cynosure.response.quickstart.QueryProjectResponseBodyByQuickStart;
import com.iflytek.ccr.polaris.cynosure.response.quickstart.QueryServiceResponseBodyByQuickStart;
import com.iflytek.ccr.polaris.cynosure.response.quickstart.QueryVersionResponseBodyByQuickStart;
import com.iflytek.ccr.polaris.cynosure.response.track.*;
import com.iflytek.ccr.polaris.cynosure.service.ILastestSearchService;
import com.iflytek.ccr.polaris.cynosure.service.ITrackService;
import com.iflytek.ccr.polaris.cynosure.util.PagingUtil;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Optional;

import static java.util.stream.Collectors.toList;

/**
 * 轨迹业务逻辑接口实现
 *
 * @author sctang2
 * @create 2017-11-24 11:57
 **/
@Service
public class TrackServiceImpl extends BaseService implements ITrackService {

    private final EasyLogger logger = EasyLoggerFactory.getInstance(TrackServiceImpl.class);
    @Autowired
    private ILastestSearchService lastestSearchServiceImpl;

    @Autowired
    private IServiceConfigPushHistoryCondition serviceConfigPushHistoryConditionImpl;

    @Autowired
    private IServiceDiscoveryPushHistoryCondition serviceDiscoveryPushHistoryConditionImpl;

    @Autowired
    private IServiceConfigPushFeedbackCondition serviceConfigPushFeedbackConditionImpl;

    @Autowired
    private IServiceDiscoveryPushFeedbackCondition serviceDiscoveryPushFeedbackConditionImpl;

    @Autowired
    private ServiceConfigPushTransactional serviceConfigPushTransactional;

    @Autowired
    private ServiceDiscoveryPushTransactional serviceDiscoveryPushTransactional;

    @Autowired
    private IProjectCondition projectConditionImpl;

    @Autowired
    private IClusterCondition clusterConditionImpl;

    @Autowired
    private IServiceCondition serviceConditionImpl;

    @Autowired
    private IServiceApiVersionCondition serviceApiVersionCondition;

    @Override
    public Response<QueryPagingListResponseBody> findLastestConfigList(QueryTrackRequestBody body) {
        QueryPagingListResponseBody result;
        String projectName = body.getProject();
        String clusterName = body.getCluster();
        String serviceName = body.getService();
        String versionName = body.getVersion();
        Integer filterGray = body.getFilterGray();

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
        map.put("filterGray", filterGray);

        //查询总数
        int totalCount = this.serviceConfigPushHistoryConditionImpl.findTotalCount(map);

        //创建分页结果
        result = PagingUtil.createResult(body, totalCount);

        //保存最近的搜索
        String condition = this.lastestSearchServiceImpl.saveLastestSearch(projectName, clusterName, serviceName, versionName);
        result.setCondition(condition);
        if (0 == totalCount) {
            return new Response<>(result);
        }

        //查询列表
        List<TrackConfigResponseBody> list = new ArrayList<>();
        Optional<List<ServiceConfigPushHistory>> serviceConfigPushHistoryList = Optional.ofNullable(this.serviceConfigPushHistoryConditionImpl.findList(map));
        serviceConfigPushHistoryList.ifPresent(x -> {
            x.forEach(y -> {
                TrackConfigResponseBody trackConfig = this.createConfigPushResult(y);
                list.add(trackConfig);
            });
        });
        result.setList(list);
        return new Response<>(result);
    }

    @Override
    public Response<QueryPagingListResponseBody> findLastestDiscoveryList(QueryTrackRequestBody body) {
        QueryPagingListResponseBody result;
        String projectName = body.getProject();
        String clusterName = body.getCluster();
        String serviceName = body.getService();
        String versionName = body.getVersion();

        //查询最近的搜索
        SearchCondition searchCondition = this.lastestSearchServiceImpl.find(projectName, clusterName, serviceName, versionName);
        projectName = searchCondition.getProject();
        clusterName = searchCondition.getCluster();
        serviceName = searchCondition.getService();
        versionName = searchCondition.getVersion();

        if (StringUtils.isBlank(projectName) || StringUtils.isBlank(clusterName) || StringUtils.isBlank(serviceName)) {
            result = PagingUtil.createResult(body, 0);
            return new Response<>(result);
        }

        //创建分页查询条件
        HashMap<String, Object> map = PagingUtil.createCondition(body);
        map.put("projectName", projectName);
        map.put("clusterName", clusterName);
        map.put("serviceName", serviceName);
        map.put("versionName", versionName);

        //查询总数
        int totalCount = this.serviceDiscoveryPushHistoryConditionImpl.findTotalCount(map);

        //创建分页结果
        result = PagingUtil.createResult(body, totalCount);

        //保存最近的搜索
        String condition = this.lastestSearchServiceImpl.saveLastestSearch(projectName, clusterName, serviceName, versionName);
        result.setCondition(condition);
        if (0 == totalCount) {
            return new Response<>(result);
        }

        //查询列表
        List<TrackDiscoveryResponseBody> list = new ArrayList<>();
        Optional<List<ServiceDiscoveryPushHistory>> serviceDiscoveryPushHistoryList = Optional.ofNullable(this.serviceDiscoveryPushHistoryConditionImpl.findList(map));
        serviceDiscoveryPushHistoryList.ifPresent(x -> {
            x.forEach(y -> {
                TrackDiscoveryResponseBody trackDiscovery = this.createDiscoveryPushResult(y);
                list.add(trackDiscovery);
            });
        });
        result.setList(list);
        return new Response<>(result);
    }

    @Override
    public Response<QueryPagingListResponseBody> findConfig(QueryTrackDetailRequestBody body) {
        QueryPagingListResponseBody result;
        String pushId = body.getPushId();
        String condition = JSON.toJSONString(body);
        if (StringUtils.isEmpty(pushId)) {
            result = PagingUtil.createResult(body, 0);
            return new Response<>(result);
        }

        //创建分页查询条件
        HashMap<String, Object> map = PagingUtil.createCondition(body);
        map.put("pushId", pushId);

        //查询总数
        int totalCount = this.serviceConfigPushFeedbackConditionImpl.findTotalCount(map);

        //创建分页结果
        result = PagingUtil.createResult(body, totalCount);
        result.setCondition(condition);
        if (0 == totalCount) {
            return new Response<>(result);
        }

        //查询列表
        List<TrackConfigDetailResponseBody> list = new ArrayList<>();
        Optional<List<ServiceConfigPushFeedback>> serviceConfigPushFeedbackList = Optional.ofNullable(this.serviceConfigPushFeedbackConditionImpl.findList(map));
        serviceConfigPushFeedbackList.ifPresent(x -> {
            x.forEach(y -> {
                TrackConfigDetailResponseBody trackConfigDetail = this.createConfigFeedbackResult(y);
                list.add(trackConfigDetail);
            });
        });
        result.setList(list);
        return new Response<>(result);
    }

    @Override
    public Response<QueryPagingListResponseBody> findDiscovery(QueryTrackDetailRequestBody body) {
        QueryPagingListResponseBody result;
        String pushId = body.getPushId();
        if (StringUtils.isEmpty(pushId)) {
            result = PagingUtil.createResult(body, 0);
            return new Response<>(result);
        }

        //创建分页查询条件
        HashMap<String, Object> map = PagingUtil.createCondition(body);
        map.put("pushId", pushId);

        //查询总数
        int totalCount = this.serviceDiscoveryPushFeedbackConditionImpl.findTotalCount(map);

        //创建分页结果
        result = PagingUtil.createResult(body, totalCount);
        if (0 == totalCount) {
            return new Response<>(result);
        }

        //查询列表
        List<TrackDiscoveryDetailResponseBody> list = new ArrayList<>();
        Optional<List<ServiceDiscoveryPushFeedbackDetail>> serviceDiscoveryPushFeedbackList = Optional.ofNullable(this.serviceDiscoveryPushFeedbackConditionImpl.findList(map));
        serviceDiscoveryPushFeedbackList.ifPresent(x -> {
            x.forEach(y -> {
                TrackDiscoveryDetailResponseBody trackConfigDetail = this.createDiscoveryFeedbackResult(y);
                list.add(trackConfigDetail);
            });
        });
        result.setList(list);
        return new Response<>(result);
    }

    @Override
    public Response<String> deleteConfig(IdRequestBody body) {
        //删除推送和反馈历史
        String id = body.getId();
        int success = this.serviceConfigPushTransactional.deletePush(id);
        if (success <= 0) {
            //不存在该轨迹
            return new Response<>(SystemErrCode.ERRCODE_TRACK_NOT_EXISTS, SystemErrCode.ERRMSG_TRACK_NOT_EXISTS);
        }
        return new Response<>(null);
    }

    @Override
    public Response<String> batchDeleteConfig(IdsRequestBody body) {
        //批量删除推送和反馈历史
        List<String> ids = body.getIds();
        int success = this.serviceConfigPushTransactional.batchDeletePush(ids);
        if (success <= 0) {
            //不存在该轨迹
            return new Response<>(SystemErrCode.ERRCODE_TRACK_NOT_EXISTS, SystemErrCode.ERRMSG_TRACK_NOT_EXISTS);
        }
        return new Response<>(null);
    }

    @Override
    public Response<String> deleteDiscovery(IdRequestBody body) {
        //删除推送和反馈历史
        String id = body.getId();
        int success = this.serviceDiscoveryPushTransactional.deletePush(id);
        if (success <= 0) {
            //不存在该轨迹
            return new Response<>(SystemErrCode.ERRCODE_TRACK_NOT_EXISTS, SystemErrCode.ERRMSG_TRACK_NOT_EXISTS);
        }
        return new Response<>(null);
    }

    @Override
    public Response<String> batchDeleteDiscovery(IdsRequestBody body) {
        //批量删除推送和反馈历史
        List<String> ids = body.getIds();
        int success = this.serviceDiscoveryPushTransactional.batchDeletePush(ids);
        if (success <= 0) {
            //不存在该轨迹
            return new Response<>(SystemErrCode.ERRCODE_TRACK_NOT_EXISTS, SystemErrCode.ERRMSG_TRACK_NOT_EXISTS);
        }
        return new Response<>(null);
    }

    @Override
    public Response<QueryPagingListResponseBody> findList(BaseRequestBody body) {
        QueryPagingListResponseBody result;
        //查询项目列表
        HashMap<String, Object> map = new HashMap<>();
        map.put("userId", this.getUserId());
        List<Project> projectList = this.projectConditionImpl.findList(map);
        if (null == projectList || projectList.isEmpty()) {
            result = PagingUtil.createResult(body, 0);
            return new Response<>(result);
        }
        int totalCount = projectList.size();
        result = PagingUtil.createResult(body, totalCount);

        //查询集群列表
        List<String> projectIds = projectList.stream().map(x -> x.getId()).collect(toList());
        List<Cluster> clusterList = this.clusterConditionImpl.findList(projectIds);

        //查询服务列表
        List<com.iflytek.ccr.polaris.cynosure.domain.Service> serviceList = new ArrayList<>();
        if (null != clusterList && !clusterList.isEmpty()) {
            List<String> clusterIds = clusterList.stream().map(x -> x.getId()).collect(toList());
            serviceList = this.serviceConditionImpl.findList(clusterIds);
        }

        //查询版本列表
        List<ServiceApiVersion> serviceApiVersionList = new ArrayList<>();
        if (null != serviceList && !serviceList.isEmpty()) {
            List<String> serviceIds = serviceList.stream().map(x -> x.getId()).collect(toList());
            serviceApiVersionList = this.serviceApiVersionCondition.findList(serviceIds);
        }

        //创建列表
        List<QueryProjectResponseBodyByQuickStart> list = this.createList(projectList, clusterList, serviceList, serviceApiVersionList);
        result.setList(list);
        return new Response<>(result);
    }

    /**
     * 创建配置反馈结果
     *
     * @param serviceConfigPushFeedback
     * @return
     */
    private TrackConfigDetailResponseBody createConfigFeedbackResult(ServiceConfigPushFeedback serviceConfigPushFeedback) {
        TrackConfigDetailResponseBody result = new TrackConfigDetailResponseBody();
        result.setId(serviceConfigPushFeedback.getId());
        result.setProject(serviceConfigPushFeedback.getProject());
        result.setCluster(serviceConfigPushFeedback.getServiceGroup());
        result.setService(serviceConfigPushFeedback.getService());
        result.setVersion(serviceConfigPushFeedback.getVersion());
        result.setConfig(serviceConfigPushFeedback.getConfig());
        result.setAddr(serviceConfigPushFeedback.getAddr());
        result.setUpdateStatus(serviceConfigPushFeedback.getUpdateStatus());
        result.setUpdateTime(serviceConfigPushFeedback.getUpdateTime().getTime());
        result.setLoadStatus(serviceConfigPushFeedback.getLoadStatus());
        result.setLoadTime(serviceConfigPushFeedback.getLoadTime().getTime());
        result.setGrayGroupId(serviceConfigPushFeedback.getGrayGroupId());
        result.setGrayGroupName(serviceConfigPushFeedback.getGrayGroupName());
        return result;
    }

    /**
     * 创建发现反馈结果
     *
     * @param serviceDiscoveryPushFeedbackDetail
     * @return
     */
    private TrackDiscoveryDetailResponseBody createDiscoveryFeedbackResult(ServiceDiscoveryPushFeedbackDetail serviceDiscoveryPushFeedbackDetail) {
        TrackDiscoveryDetailResponseBody result = new TrackDiscoveryDetailResponseBody();
        result.setId(serviceDiscoveryPushFeedbackDetail.getId());
        result.setPushId(serviceDiscoveryPushFeedbackDetail.getPushId());
        result.setProject(serviceDiscoveryPushFeedbackDetail.getProject());
        result.setCluster(serviceDiscoveryPushFeedbackDetail.getServiceGroup());
        result.setConsumerService(serviceDiscoveryPushFeedbackDetail.getConsumerService());
        result.setConsumerVersion(serviceDiscoveryPushFeedbackDetail.getConsumerVersion());
        result.setProviderService(serviceDiscoveryPushFeedbackDetail.getProviderService());
        result.setProviderVersion(serviceDiscoveryPushFeedbackDetail.getProviderVersion());
        result.setAddr(serviceDiscoveryPushFeedbackDetail.getAddr());
        result.setUpdateStatus(serviceDiscoveryPushFeedbackDetail.getUpdateStatus());
        result.setUpdateTime(serviceDiscoveryPushFeedbackDetail.getUpdateTime().getTime());
        result.setLoadStatus(serviceDiscoveryPushFeedbackDetail.getLoadStatus());
        result.setLoadTime(serviceDiscoveryPushFeedbackDetail.getLoadTime().getTime());
        result.setApiVersion(serviceDiscoveryPushFeedbackDetail.getApiVersion());
        result.setType(serviceDiscoveryPushFeedbackDetail.getType());
        result.setTypeName(serviceDiscoveryPushFeedbackDetail.getTypeName());
        return result;
    }

    /**
     * 创建配置推送结果
     *
     * @param serviceConfigPushHistory
     * @return
     */
    private TrackConfigResponseBody createConfigPushResult(ServiceConfigPushHistory serviceConfigPushHistory) {
        TrackConfigResponseBody result = new TrackConfigResponseBody();
        result.setId(serviceConfigPushHistory.getId());
        String configText = serviceConfigPushHistory.getServiceConfigText();
        List<TrackConfig> trackConfigs = JSON.parseArray(configText, TrackConfig.class);
        result.setConfigs(trackConfigs);
        String regionText = serviceConfigPushHistory.getClusterText();
        List<TrackRegion> trackRegions = JSON.parseArray(regionText, TrackRegion.class);
        result.setRegions(trackRegions);
        result.setPushTime(serviceConfigPushHistory.getPushTime());
        result.setGrayGroupId(serviceConfigPushHistory.getGrayId());
        return result;
    }

    /**
     * 创建服务发现推送结果
     *
     * @param serviceDiscoveryPushHistory
     * @return
     */
    private TrackDiscoveryResponseBody createDiscoveryPushResult(ServiceDiscoveryPushHistory serviceDiscoveryPushHistory) {
        TrackDiscoveryResponseBody result = new TrackDiscoveryResponseBody();
        result.setId(serviceDiscoveryPushHistory.getId());
        String regionText = serviceDiscoveryPushHistory.getClusterText();
        List<TrackRegion> trackRegions = JSON.parseArray(regionText, TrackRegion.class);
        result.setRegions(trackRegions);
        result.setPushTime(serviceDiscoveryPushHistory.getPushTime());
        return result;
    }

    /**
     * 创建列表
     *
     * @param projectList
     * @param clusterList
     * @param serviceList
     * @param serviceVersionList
     * @return
     */
    private List<QueryProjectResponseBodyByQuickStart> createList(List<Project> projectList, List<Cluster> clusterList, List<com.iflytek.ccr.polaris.cynosure.domain.Service> serviceList, List<ServiceApiVersion> serviceVersionList) {
        List<QueryProjectResponseBodyByQuickStart> projectResults = new ArrayList<>();
        projectList.forEach(x -> {
            QueryProjectResponseBodyByQuickStart project = new QueryProjectResponseBodyByQuickStart();
            String projectId = x.getId();
            String projectName = x.getName();
            project.setId(projectId);
            project.setName(projectName);
            List<QueryClusterResponseBodyByQuickStart> clusterResults = new ArrayList<>();
            clusterList.forEach(y -> {
                if (projectId.equals(y.getProjectId())) {
                    String clusterId = y.getId();
                    String clusterName = y.getName();
                    QueryClusterResponseBodyByQuickStart cluster = new QueryClusterResponseBodyByQuickStart();
                    cluster.setId(clusterId);
                    cluster.setName(clusterName);
                    List<QueryServiceResponseBodyByQuickStart> serviceResults = new ArrayList<>();
                    serviceList.forEach(z -> {
                        String serviceId = z.getId();
                        String serviceName = z.getName();
                        List<QueryVersionResponseBodyByQuickStart> versionResults = new ArrayList<>();
                        if (z.getGroupId().equals(clusterId)) {
                            QueryServiceResponseBodyByQuickStart service = new QueryServiceResponseBodyByQuickStart();
                            service.setId(serviceId);
                            service.setName(serviceName);
                            serviceVersionList.forEach(n -> {
                                if (n.getServiceId().equals(serviceId)) {
                                    QueryVersionResponseBodyByQuickStart version = new QueryVersionResponseBodyByQuickStart();
                                    version.setId(n.getId());
                                    version.setName(n.getApiVersion());
                                    versionResults.add(version);
                                }
                            });
                            if (null != versionResults && !versionResults.isEmpty()) {
                                service.setChildren(versionResults);
                            }
                            serviceResults.add(service);
                        }
                    });
                    if (null != serviceResults && !serviceResults.isEmpty()) {
                        cluster.setChildren(serviceResults);
                    }
                    clusterResults.add(cluster);
                }
            });
            if (null != clusterResults && !clusterResults.isEmpty()) {
                project.setChildren(clusterResults);
            }
            projectResults.add(project);
        });
        return projectResults;
    }
}
