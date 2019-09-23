package com.iflytek.ccr.polaris.cynosure.service.impl;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.companionservice.ConfigCenter;
import com.iflytek.ccr.polaris.cynosure.companionservice.GrayConfigCenter;
import com.iflytek.ccr.polaris.cynosure.companionservice.domain.ServiceProviderConsumerResult;
import com.iflytek.ccr.polaris.cynosure.companionservice.domain.ServiceResult;
import com.iflytek.ccr.polaris.cynosure.customdomain.FileContent;
import com.iflytek.ccr.polaris.cynosure.dbcondition.*;
import com.iflytek.ccr.polaris.cynosure.dbtransactional.GrayGroupTransactional;
import com.iflytek.ccr.polaris.cynosure.domain.*;
import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.request.graygroup.*;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.graygroup.AddGrayGroupAndConfigResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.graygroup.GrayGroupDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.serviceconfig.ServiceConfigDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.servicediscovery.QueryServiceDiscoveryConsumerResponseBody;
import com.iflytek.ccr.polaris.cynosure.service.IGrayService;
import com.iflytek.ccr.polaris.cynosure.service.ILastestSearchService;
import com.iflytek.ccr.polaris.cynosure.util.PagingUtil;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.*;

/**
 * 灰度组业务接口实现
 *
 * @author sctang2
 * @create 2017-11-17 16:03
 **/
@Service
public class GrayServiceImpl extends BaseService implements IGrayService {
    private final EasyLogger logger = EasyLoggerFactory.getInstance(GrayServiceImpl.class);

    @Autowired
    private IProjectCondition projectConditionImpl;

    @Autowired
    private IProjectMemberCondition projectMemberConditionImpl;

    @Autowired
    private IClusterCondition clusterConditionImpl;

    @Autowired
    private IServiceCondition serviceConditionImpl;

    @Autowired
    private IGrayGroupCondition grayGroupConditionImpl;

    @Autowired
    private IServiceVersionCondition serviceVersionConditionImpl;

    @Autowired
    private GrayGroupTransactional grayGroupTransactionalImpl;

    @Autowired
    private IServiceConfigCondition serviceConfigConditionImpl;

    @Autowired
    private ILastestSearchService lastestSearchServiceImpl;

    @Autowired
    private GrayConfigCenter grayConfigCenter;

    @Autowired
    private ConfigCenter configCenter;

    @Autowired
    private IRegionCondition regionConditionImpl;

    @Autowired
    private InstanceManageCondition instanceManageConditionImpl;

    @Override
    public Response<AddGrayGroupAndConfigResponseBody> add(AddGrayGroupRequestBody body) {
        String projectId = null;
        String clusterId = null;
        String serviceId = null;

        //通过项目名称查询项目信息
        String projectName = body.getProject();
        Project project = this.projectConditionImpl.findByName(projectName);
        if (null == project) {
            //不存在该项目
            return new Response<>(SystemErrCode.ERRCODE_PROJECT_NOT_EXISTS, SystemErrCode.ERRMSG_PROJECT_NOT_EXISTS);
        } else {
            projectId = project.getId();
        }

        //根据id和集群名称查询集群信息
        String clusterName = body.getCluster();
        Cluster cluster = null;
        if (StringUtils.isNotBlank(projectId)) {
            cluster = this.clusterConditionImpl.find(projectId, clusterName);
            if (null == cluster) {
                //不存在该集群
                return new Response<>(SystemErrCode.ERRCODE_CLUSTER_NOT_EXISTS, SystemErrCode.ERRMSG_CLUSTER_NOT_EXISTS);
            } else {
                clusterId = cluster.getId();
            }
        }

        //根据服务名和集群id查询服务
        String serviceName = body.getService();
        com.iflytek.ccr.polaris.cynosure.domain.Service service = null;
        if (StringUtils.isNotBlank(clusterId)) {
            service = this.serviceConditionImpl.find(serviceName, clusterId);
            if (null == service) {
                //不存在该服务
                return new Response<>(SystemErrCode.ERRCODE_SERVICE_NOT_EXISTS, SystemErrCode.ERRMSG_SERVICE_NOT_EXISTS);
            } else {
                serviceId = service.getId();
            }
        }

        //根据版本和服务id查询服务版本
        String versionName = body.getVersion();
        ServiceVersion serviceVersion = null;
        if (StringUtils.isNotBlank(serviceId)) {
            serviceVersion = this.serviceVersionConditionImpl.find(versionName, serviceId);
            if (null == serviceVersion) {
                //不存在该版本
                return new Response<>(SystemErrCode.ERRCODE_SERVICE_VERSION_NOT_EXISTS, SystemErrCode.ERRMSG_SERVICE_VERSION_NOT_EXISTS);
            }
        }

        //通过id查询服务
        String versionId = body.getVersionId();
        ServiceVersion version = this.serviceVersionConditionImpl.findById(versionId);
        if (null == version) {
            //不存在该版本
            return new Response<>(SystemErrCode.ERRCODE_SERVICE_VERSION_NOT_EXISTS, SystemErrCode.ERRMSG_SERVICE_VERSION_NOT_EXISTS);
        }

        //根据灰度组名称和版本id查询灰度组
        String grayName = body.getName();
        GrayGroup grayGroup = this.grayGroupConditionImpl.find(grayName, versionId);
        if (null != grayGroup) {
            return new Response<>(SystemErrCode.ERRCODE_GRAY_GROUP_EXISTS, SystemErrCode.ERRMSG_GRAY_GROUP_EXISTS);
        }

        //校验老的版本是下的推送实例是否重复
        String instanceContent = body.getContent();
        if (StringUtils.isNotBlank(instanceContent)) {
            List<String> contentList = Arrays.asList(instanceContent.split(","));
            List<String> versionContent = this.instanceManageConditionImpl.findTotal(body.getVersionId(), null);
            if (null != versionContent && !versionContent.isEmpty()) {
                List<String> oldGroupList = new ArrayList<>();
                for (String content : versionContent) {
                    oldGroupList.addAll(Arrays.asList(StringUtils.split(content, ",")));
                }

                //校验灰度组的推送实例是否有重复内容
                oldGroupList.addAll(contentList);
                HashSet<String> compareSet = new HashSet<>(oldGroupList);
                if (oldGroupList.size() != compareSet.size()) {
                    return new Response<>(SystemErrCode.ERRCODE_GRAY_INSTANCE_ARE_USED, SystemErrCode.ERRMSG_GRAY_INSTANCE_ARE_USED);
                }
            }
        }
        return this.grayGroupTransactionalImpl.addGrayGroup(body);
    }

    @Override
    public Response<AddGrayGroupAndConfigResponseBody> addAndFile(AddGrayGroupRequestBody body, List<FileContent> fileContentList) {
        String projectId = null;
        String clusterId = null;
        String serviceId = null;

        //通过项目名称查询项目信息
        String projectName = body.getProject();
        Project project = this.projectConditionImpl.findByName(projectName);
        if (null == project) {
            //不存在该项目
            return new Response<>(SystemErrCode.ERRCODE_PROJECT_NOT_EXISTS, SystemErrCode.ERRMSG_PROJECT_NOT_EXISTS);
        } else {
            projectId = project.getId();
        }

        //根据id和集群名称查询集群信息
        String clusterName = body.getCluster();
        Cluster cluster = null;
        if (StringUtils.isNotBlank(projectId)) {
            cluster = this.clusterConditionImpl.find(projectId, clusterName);
            if (null == cluster) {
                //不存在该集群
                return new Response<>(SystemErrCode.ERRCODE_CLUSTER_NOT_EXISTS, SystemErrCode.ERRMSG_CLUSTER_NOT_EXISTS);
            } else {
                clusterId = cluster.getId();
            }
        }

        //根据服务名和集群id查询服务
        String serviceName = body.getService();
        com.iflytek.ccr.polaris.cynosure.domain.Service service = null;
        if (StringUtils.isNotBlank(clusterId)) {
            service = this.serviceConditionImpl.find(serviceName, clusterId);
            if (null == service) {
                //不存在该服务
                return new Response<>(SystemErrCode.ERRCODE_SERVICE_NOT_EXISTS, SystemErrCode.ERRMSG_SERVICE_NOT_EXISTS);
            } else {
                serviceId = service.getId();
            }
        }

        //根据版本和服务id查询服务版本
        String versionName = body.getVersion();
        ServiceVersion serviceVersion = null;
        if (StringUtils.isNotBlank(serviceId)) {
            serviceVersion = this.serviceVersionConditionImpl.find(versionName, serviceId);
            if (null == serviceVersion) {
                //不存在该版本
                return new Response<>(SystemErrCode.ERRCODE_SERVICE_VERSION_NOT_EXISTS, SystemErrCode.ERRMSG_SERVICE_VERSION_NOT_EXISTS);
            }
        }

        //通过id查询版本
        String versionId = body.getVersionId();
        ServiceVersion version = this.serviceVersionConditionImpl.findById(versionId);
        if (null == version) {
            //不存在该版本
            return new Response<>(SystemErrCode.ERRCODE_SERVICE_VERSION_NOT_EXISTS, SystemErrCode.ERRMSG_SERVICE_VERSION_NOT_EXISTS);
        }

        //根据灰度组名称和版本id查询灰度组
        String grayName = body.getName();
        GrayGroup grayGroup = this.grayGroupConditionImpl.find(grayName, versionId);
        if (null != grayGroup) {
            return new Response<>(SystemErrCode.ERRCODE_GRAY_GROUP_EXISTS, SystemErrCode.ERRMSG_GRAY_GROUP_EXISTS);
        }

        //校验老的版本是下的推送实例是否重复
        String instanceContent = body.getContent();
        if (StringUtils.isNotBlank(instanceContent)) {
            List<String> contentList = Arrays.asList(instanceContent.split(","));
            List<String> versionContent = this.instanceManageConditionImpl.findTotal(body.getVersionId(), null);
            if (null != versionContent && !versionContent.isEmpty()) {
                List<String> oldGroupList = new ArrayList<>();
                for (String content : versionContent) {
                    oldGroupList.addAll(Arrays.asList(StringUtils.split(content, ",")));
                }
                //校验灰度组的推送实例是否有重复内容
                oldGroupList.addAll(contentList);
                HashSet<String> compareSet = new HashSet<>(oldGroupList);
                if (oldGroupList.size() != compareSet.size()) {
                    return new Response<>(SystemErrCode.ERRCODE_GRAY_INSTANCE_ARE_USED, SystemErrCode.ERRMSG_GRAY_INSTANCE_ARE_USED);
                }
            }
        }
        return this.grayGroupTransactionalImpl.addGrayGroupAndFile(body, fileContentList);
    }

    @Override
    public Response<QueryPagingListResponseBody> findList(QueryGrayGroupListRequestBody body) {

        HashMap<String, Object> map = new HashMap<>();
        String versionId = body.getVersionId();
        map.put("versionId", versionId);
        List<GrayGroup> grayGroupList = this.grayGroupConditionImpl.findList(map);
        QueryPagingListResponseBody result = new QueryPagingListResponseBody();
        result.setList(grayGroupList);
        return new Response<>(result);
    }

    @Override
    public Response<AddGrayGroupAndConfigResponseBody> findById(String id) {//把两部分都组合起来然后拼接
        //根据id查询灰度组
        GrayGroup grayGroup = this.grayGroupConditionImpl.findById(id);
        if (null == grayGroup) {
            //不存在该灰度组
            return new Response<>(SystemErrCode.ERRCODE_GRAY_GROUP_NOT_EXISTS, SystemErrCode.ERRMSG_GRAY_GROUP_NOT_EXISTS);
        }

        //创建灰度组详情结果
        AddGrayGroupAndConfigResponseBody result = new AddGrayGroupAndConfigResponseBody();

        //创建灰度组结果
        GrayGroupDetailResponseBody grayGroupResult = this.createGrayGroupResult(grayGroup);
        result.setGrayGroup(grayGroupResult);

        //创建灰度配置结果
        List<ServiceConfig> serviceConfigList = this.serviceConfigConditionImpl.findListByGrayId(id);
        List<ServiceConfigDetailResponseBody> configs = new ArrayList<>();
        if (null != serviceConfigList && !serviceConfigList.isEmpty()) {
            for (ServiceConfig serviceConfig : serviceConfigList) {
                configs.add(this.createServiceConfigResult(serviceConfig));
            }
            result.setConfigs(configs);
        }
        return new Response<>(result);
    }

    @Override
    public Response<String> delete(DeleteGrayGroupRequestBody body) {
        String projectId = null;
        String clusterId = null;
        String serviceId = null;

        //通过项目名称查询项目信息
        String projectName = body.getProject();
        Project project = this.projectConditionImpl.findByName(projectName);
        if (null == project) {
            //不存在该项目
            return new Response<>(SystemErrCode.ERRCODE_PROJECT_NOT_EXISTS, SystemErrCode.ERRMSG_PROJECT_NOT_EXISTS);
        } else {
            projectId = project.getId();
        }

        //根据id和集群名称查询集群信息
        String clusterName = body.getCluster();
        Cluster cluster = null;
        if (StringUtils.isNotBlank(projectId)) {
            cluster = this.clusterConditionImpl.find(projectId, clusterName);
            if (null == cluster) {
                //不存在该集群
                return new Response<>(SystemErrCode.ERRCODE_CLUSTER_NOT_EXISTS, SystemErrCode.ERRMSG_CLUSTER_NOT_EXISTS);
            } else {
                clusterId = cluster.getId();
            }
        }

        //根据服务名和集群id查询服务
        String serviceName = body.getService();
        com.iflytek.ccr.polaris.cynosure.domain.Service service = null;
        if (StringUtils.isNotBlank(clusterId)) {
            service = this.serviceConditionImpl.find(serviceName, clusterId);
            if (null == service) {
                //不存在该服务
                return new Response<>(SystemErrCode.ERRCODE_SERVICE_NOT_EXISTS, SystemErrCode.ERRMSG_SERVICE_NOT_EXISTS);
            } else {
                serviceId = service.getId();
            }
        }

        //根据版本和服务id查询服务版本
        String versionName = body.getVersion();
        ServiceVersion serviceVersion = null;
        if (StringUtils.isNotBlank(serviceId)) {
            serviceVersion = this.serviceVersionConditionImpl.find(versionName, serviceId);
            if (null == serviceVersion) {
                //不存在该版本
                return new Response<>(SystemErrCode.ERRCODE_SERVICE_VERSION_NOT_EXISTS, SystemErrCode.ERRMSG_SERVICE_VERSION_NOT_EXISTS);
            }
        }


        //通过id查询灰度组
        String id = body.getGrayId();
        GrayGroup grayGroup = this.grayGroupConditionImpl.findById(id);
        if (null == grayGroup) {
            //不存在灰度组
            return new Response<>(SystemErrCode.ERRCODE_GRAY_GROUP_NOT_EXISTS, SystemErrCode.ERRMSG_GRAY_GROUP_NOT_EXISTS);
        }
        List<String> grayIdList = new ArrayList<>();
        grayIdList.add(id);
        List<ServiceConfig> grayConfigs = this.serviceConfigConditionImpl.findListByGrayIds(grayIdList);

        if (null != grayConfigs && !grayConfigs.isEmpty()) {
            //同时删除灰度组和灰度配置文件
            List<String> configIds = new ArrayList<>();
            for (ServiceConfig serviceConfig : grayConfigs) {
                configIds.add(serviceConfig.getId());
            }
            this.grayGroupTransactionalImpl.deleteGroupAndConfig(id, configIds);

            //查询集群列表
            List<Region> regionList = this.regionConditionImpl.findList(null);
            if (null == regionList || regionList.isEmpty()) {
                return new Response<>(null);
            }

            //整体删除灰度组缓存中心
            String path = StringUtils.substringBeforeLast(grayConfigs.get(0).getPath(), "/");
            logger.info(path);
            this.grayConfigCenter.deleteGrayGroup(path, id, regionList);
            return new Response<>(null);
        } else {
            //根据id删除灰度组
            this.grayGroupConditionImpl.deleteById(id);

            //查询集群列表
            List<Region> regionList = this.regionConditionImpl.findList(null);
            if (null == regionList || regionList.isEmpty()) {
                return new Response<>(null);
            }

            //整体删除灰度组缓存中心
            String path = StringUtils.substringBeforeLast(this.configCenter.getConfigPath(projectName, clusterName, serviceName, versionName, null), "/");
            logger.info(path);
            this.grayConfigCenter.deleteGrayGroup(path, id, regionList);

            return new Response<>(null);
        }
    }

    @Override
    public Response<GrayGroupDetailResponseBody> edit(EditGrayGroupRequestBody body) {
        //根据id查询灰度组
        String id = body.getId();
        GrayGroup grayGroup = this.grayGroupConditionImpl.findById(id);
        if (null == grayGroup) {
            return new Response<>(SystemErrCode.ERRCODE_GRAY_GROUP_NOT_EXISTS, SystemErrCode.ERRMSG_GRAY_GROUP_NOT_EXISTS);
        }

        //根据id更新灰度组
        GrayGroup updateGrayGroup = this.grayGroupConditionImpl.updateById(id, body);

        //创建灰度组结果
        updateGrayGroup.setCreateTime(grayGroup.getCreateTime());
        updateGrayGroup.setContent(grayGroup.getContent());
        updateGrayGroup.setName(grayGroup.getName());
        updateGrayGroup.setUserId(grayGroup.getUserId());
        updateGrayGroup.setVersionId(grayGroup.getVersionId());
        GrayGroupDetailResponseBody result = this.createGrayGroupResult(updateGrayGroup);
        return new Response<>(result);
    }

    @Override
    public Response<QueryPagingListResponseBody> consumer(QueryCustomDetailRequestBody body) {
        QueryPagingListResponseBody result;
        String grayId = body.getGrayId();

        //通过区域名称查询区域信息
        String regionName = body.getRegion();
        Region region = this.regionConditionImpl.findByName(regionName);
        if (null == region) {
            return new Response<>(SystemErrCode.ERRCODE_REGION_NOT_EXISTS, SystemErrCode.ERRMSG_REGION_NOT_EXISTS);
        }
        //获取服务消费者path
        String project = body.getProject();
        String cluster = body.getCluster();
        String service = body.getService();
        String version = body.getVersion();
        String path = this.grayConfigCenter.getServiceConsumerPath(project, cluster, service, version, grayId);

        //查询服务消费端
        int startIndex = PagingUtil.getStartIndex(body);
        int endIndex = PagingUtil.getEndIndex(body);
        ServiceResult serviceResult = this.grayConfigCenter.findConsumersByPaging(path, region, false, startIndex, endIndex);
        int totalCount = serviceResult.getTotalCount();
        result = PagingUtil.createResult(body, totalCount);

        //创建消费端
        List<ServiceProviderConsumerResult> serviceProviderConsumerResults = serviceResult.getResults();
        if (null != serviceProviderConsumerResults && !serviceProviderConsumerResults.isEmpty()) {
            List<QueryServiceDiscoveryConsumerResponseBody> serviceDiscoveryConsumerList = this.createConsumer(serviceProviderConsumerResults);
            result.setList(serviceDiscoveryConsumerList);
        }
        return new Response<>(result);
    }

    /**
     * 创建消费端
     *
     * @param cacheCenterServiceProviderConsumerResults
     * @return
     */
    private List<QueryServiceDiscoveryConsumerResponseBody> createConsumer(List<ServiceProviderConsumerResult> cacheCenterServiceProviderConsumerResults) {
        List<QueryServiceDiscoveryConsumerResponseBody> results = new ArrayList<>();
        cacheCenterServiceProviderConsumerResults.forEach(x -> {
            QueryServiceDiscoveryConsumerResponseBody result = new QueryServiceDiscoveryConsumerResponseBody();
            result.setAddr(x.getAddr());
            results.add(result);
        });
        return results;
    }
}
