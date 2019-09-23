package com.iflytek.ccr.polaris.cynosure.service.impl;

import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.customdomain.FileContent;
import com.iflytek.ccr.polaris.cynosure.dbcondition.*;
import com.iflytek.ccr.polaris.cynosure.dbtransactional.QuickStartTransactional;
import com.iflytek.ccr.polaris.cynosure.domain.Cluster;
import com.iflytek.ccr.polaris.cynosure.domain.Project;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceApiVersion;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceVersion;
import com.iflytek.ccr.polaris.cynosure.request.BaseRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.quickstart.AddServiceConfigRequestBodyByQuickStart;
import com.iflytek.ccr.polaris.cynosure.request.quickstart.AddServiceRequestBodyByQuickStart;
import com.iflytek.ccr.polaris.cynosure.request.quickstart.AddServiceVersionRequestBodyByQuickStart;
import com.iflytek.ccr.polaris.cynosure.request.quickstart.AddVersionRequestBodyByQuickStart;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.quickstart.*;
import com.iflytek.ccr.polaris.cynosure.service.IQuickStartService;
import com.iflytek.ccr.polaris.cynosure.service.IServiceApiVersion;
import com.iflytek.ccr.polaris.cynosure.util.PagingUtil;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;

import static java.util.stream.Collectors.toList;

/**
 * 快速创建业务逻辑接口实现
 *
 * @author sctang2
 * @create 2018-01-29 14:03
 **/
@Service
public class QuickStartServiceImpl extends BaseService implements IQuickStartService {
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

    @Autowired
    private IServiceApiVersion serviceApiVersionImpl;

    @Autowired
    private IServiceConfigCondition iServiceConfigConditionImpl;

    @Override
    public Response<AddServiceResponseBodyByQuickStart> addService(AddServiceRequestBodyByQuickStart body) {
        return this.quickStartTransactional.addService(body, true);
    }
    @Override
    public Response<AddVersionResponseBodyByQuickStart> addVersion(AddVersionRequestBodyByQuickStart body) {
        return this.quickStartTransactional.addVersion(body, true);
    }

    @Override
    public Response<AddServiceConfigResponseBodyByQuickStart> addServiceVersion(AddServiceVersionRequestBodyByQuickStart body) {
        return this.quickStartTransactional.addServiceVersion(body);
    }

    @Override
    public Response<AddServiceConfigResponseBodyByQuickStart> addServiceVersionAndFile(AddServiceVersionRequestBodyByQuickStart body, List<FileContent> fileContentList) {
        return this.quickStartTransactional.addServiceVersionAndFile(body, fileContentList);
    }

    @Override
    public Response<AddServiceConfigResponseBodyByQuickStart> addServiceConfig(AddServiceConfigRequestBodyByQuickStart body, List<FileContent> fileContentList) {
        return this.quickStartTransactional.addServiceConfig(body, fileContentList);
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
        List<ServiceVersion> serviceVersionList = new ArrayList<>();
        if (null != serviceList && !serviceList.isEmpty()) {
            List<String> serviceIds = serviceList.stream().map(x -> x.getId()).collect(toList());
            serviceVersionList = this.serviceVersionConditionImpl.findList(serviceIds);
        }

        //创建列表
        List<QueryProjectResponseBodyByQuickStart> list = this.createList(projectList, clusterList, serviceList, serviceVersionList);
        result.setList(list);
        return new Response<>(result);
    }

    @Override
    public Response<QueryPagingListResponseBody> findList1(BaseRequestBody body) {
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
            serviceApiVersionList = this.serviceApiVersionImpl.findList1(serviceIds);
        }

        //创建列表
        List<QueryProjectResponseBodyByQuickStart> list = this.createList1(projectList, clusterList, serviceList, serviceApiVersionList);
        result.setList(list);
        return new Response<>(result);
    }

    /**
     * 创建列表
     *
     * @param projectList
     * @param clusterList
     * @param serviceList
     * @param serviceApiVersionList
     * @return
     */
    private List<QueryProjectResponseBodyByQuickStart> createList1(List<Project> projectList, List<Cluster> clusterList, List<com.iflytek.ccr.polaris.cynosure.domain.Service> serviceList, List<ServiceApiVersion> serviceApiVersionList) {
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
                            serviceApiVersionList.forEach(n -> {
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


    /**
     * 创建列表
     *
     * @param projectList
     * @param clusterList
     * @param serviceList
     * @param serviceVersionList
     * @return
     */
    private List<QueryProjectResponseBodyByQuickStart> createList(List<Project> projectList, List<Cluster> clusterList, List<com.iflytek.ccr.polaris.cynosure.domain.Service> serviceList, List<ServiceVersion> serviceVersionList) {
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
                                    version.setName(n.getVersion());
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
