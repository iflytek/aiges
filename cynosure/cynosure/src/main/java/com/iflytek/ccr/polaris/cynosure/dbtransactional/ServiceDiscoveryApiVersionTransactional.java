package com.iflytek.ccr.polaris.cynosure.dbtransactional;

import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.dbcondition.*;
import com.iflytek.ccr.polaris.cynosure.domain.Cluster;
import com.iflytek.ccr.polaris.cynosure.domain.Project;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceApiVersion;
import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.request.cluster.AddClusterRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.project.AddProjectRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.service.AddServiceRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.servicediscovery.AddServiceApiVersionRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.servicediscovery.AddServiceDiscoveryRequestBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.cluster.ClusterDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.project.ProjectDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.service.ServiceDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.servicediscovery.AddApiVersionResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.servicediscovery.ServiceApiVersionDetailResponseBody;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

/**
 * 服务注册事务
 *
 * @author sctang2
 * @create 2017-12-20 9:09
 **/
@Service
public class ServiceDiscoveryApiVersionTransactional extends BaseService {
    @Autowired
    private IProjectCondition projectConditionImpl;

    @Autowired
    private IProjectMemberCondition projectMemberConditionImpl;

    @Autowired
    private IClusterCondition clusterConditionImpl;

    @Autowired
    private IServiceCondition serviceConditionImpl;

    @Autowired
    private IServiceVersionCondition serviceVersionConditionImpl;

    @Autowired
    private IServiceApiVersionCondition serviceApiVersionConditionImpl;

    @Autowired
    private IServiceConfigCondition serviceConfigConditionImpl;

    /**
     * 新增api版本
     *
     * @param body
     * @param isLogin
     * @return
     */
    @Transactional
    public Response<AddApiVersionResponseBody> addApiVersion(AddServiceDiscoveryRequestBody body, boolean isLogin) {
        String projectId = null;
        String clusterId = null;
        String serviceId = null;

        //判断是否为管理员
        boolean isAdmin = true;
        if (isLogin) {
            isAdmin = this.isAdmin();
        }

        //通过项目名称查询项目信息
        String projectName = body.getProject();
        Project project = this.projectConditionImpl.findByName(projectName);
        if (null == project) {
            if (!isLogin) {
                //不存在项目
                return new Response<>(SystemErrCode.ERRCODE_PROJECT_NOT_EXISTS, SystemErrCode.ERRMSG_PROJECT_NOT_EXISTS);
            }
            if (!isAdmin) {
                //没有权限执行此操作
                return new Response<>(SystemErrCode.ERRCODE_NOT_AUTH, SystemErrCode.ERRMSG_NOT_AUTH);
            }
        } else {
            projectId = project.getId();
        }

        //根据id和集群名称查询集群信息
        String clusterName = body.getGroup();
        Cluster cluster = null;
        if (StringUtils.isNotBlank(projectId)) {
            cluster = this.clusterConditionImpl.find(projectId, clusterName);
            if (null != cluster) {
                clusterId = cluster.getId();
            } else {
                if (!isLogin) {
                    //不存在该集群
                    return new Response<>(SystemErrCode.ERRCODE_CLUSTER_NOT_EXISTS, SystemErrCode.ERRMSG_CLUSTER_NOT_EXISTS);
                }
            }
        }

        //根据id和服务名称查询服务信息
        String serviceName = body.getService();
        com.iflytek.ccr.polaris.cynosure.domain.Service service = null;
        if (StringUtils.isNotBlank(clusterId)) {
            service = this.serviceConditionImpl.find(serviceName,clusterId);
            if (null != service) {
                serviceId = service.getId();
            } else {
                if (!isLogin) {
                    //不存在该服务
                    return new Response<>(SystemErrCode.ERRCODE_SERVICE_NOT_EXISTS, SystemErrCode.ERRMSG_SERVICE_NOT_EXISTS);
                }
            }
        }

        //根据版本名和服务id查询版本
        String apiVersionName = body.getApiVersion();
        ServiceApiVersion apiVersion = null;
        if (StringUtils.isNotBlank(serviceId)) {
            apiVersion = this.serviceApiVersionConditionImpl.find(apiVersionName, serviceId);
            if (null != apiVersion && isLogin) {
                //已存在该版本
                return new Response<>(SystemErrCode.ERRCODE_SERVICE_VERSION_EXISTS, SystemErrCode.ERRMSG_SERVICE_VERSION_EXISTS);
            }
        }

        //新增项目
        if (StringUtils.isBlank(projectId)) {
            project = this.addProject(projectName);
            projectId = project.getId();
        }

        //新增集群
        if (StringUtils.isBlank(clusterId)) {
            cluster = this.addCluster(projectId, clusterName);
            clusterId = cluster.getId();
        }

        //新增服务
        if (StringUtils.isBlank(serviceId)) {
            service = this.addService(clusterId, serviceName);
            serviceId = service.getId();
        }

        //新增版本
        if (null == apiVersion) {
            apiVersion = this.addServiceApiVersion(serviceId, apiVersionName, null);
        }

        //创建快速结果
        AddApiVersionResponseBody result = new AddApiVersionResponseBody();

        //创建项目结果
        ProjectDetailResponseBody projectResult = this.createProjectResult(project);
        result.setProject(projectResult);

        //创建集群结果
        ClusterDetailResponseBody clusterResult = this.createClusterResult(cluster);
        result.setCluster(clusterResult);

        //创建服务结果
        ServiceDetailResponseBody serviceResult = this.createServiceResult(service);
        result.setService(serviceResult);

        //创建服务结果
        ServiceApiVersionDetailResponseBody apiVersionResult = this.createServiceApiVersionResult(apiVersion);
        result.setApiVersion(apiVersionResult);

        return new Response<>(null);
    }

    /**
     * 新增项目
     *
     * @param projectName
     * @return
     */
    private Project addProject(String projectName) {
        //新增项目
        AddProjectRequestBody body = new AddProjectRequestBody();
        body.setName(projectName);
        Project project = this.projectConditionImpl.add(body);

        //新增项目创建者
        String projectId = project.getId();
        this.projectMemberConditionImpl.addCreator(projectId);
        return project;
    }

    /**
     * 新增集群
     *
     * @param projectId
     * @param clusterName
     * @return
     */
    private Cluster addCluster(String projectId, String clusterName) {
        AddClusterRequestBody body = new AddClusterRequestBody();
        body.setProjectId(projectId);
        body.setName(clusterName);
        return this.clusterConditionImpl.add(body);
    }

    /**
     * 新增服务
     *
     * @param clusterId
     * @param serviceName
     * @return
     */
    private com.iflytek.ccr.polaris.cynosure.domain.Service addService(String clusterId, String serviceName) {
        AddServiceRequestBody body = new AddServiceRequestBody();
        body.setClusterId(clusterId);
        body.setName(serviceName);
        return this.serviceConditionImpl.add(body);
    }

    /**
     * 新增版本
     *
     * @param serviceId
     * @param apiVersionName
     * @param desc
     * @return
     */
    private ServiceApiVersion addServiceApiVersion(String serviceId, String apiVersionName, String desc) {
        AddServiceApiVersionRequestBody body = new AddServiceApiVersionRequestBody();
        body.setServiceId(serviceId);
        body.setApiVersion(apiVersionName);
        body.setDesc(desc);
        return this.serviceApiVersionConditionImpl.add(body);
    }
}
