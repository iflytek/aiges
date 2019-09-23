package com.iflytek.ccr.polaris.cynosure.dbtransactional;

import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.customdomain.FileContent;
import com.iflytek.ccr.polaris.cynosure.dbcondition.*;
import com.iflytek.ccr.polaris.cynosure.domain.Cluster;
import com.iflytek.ccr.polaris.cynosure.domain.Project;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceConfig;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceVersion;
import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.request.cluster.AddClusterRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.graygroup.AddGrayConfigRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.project.AddProjectRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.quickstart.AddServiceConfigRequestBodyByQuickStart;
import com.iflytek.ccr.polaris.cynosure.request.quickstart.AddServiceRequestBodyByQuickStart;
import com.iflytek.ccr.polaris.cynosure.request.quickstart.AddServiceVersionRequestBodyByQuickStart;
import com.iflytek.ccr.polaris.cynosure.request.quickstart.AddVersionRequestBodyByQuickStart;
import com.iflytek.ccr.polaris.cynosure.request.service.AddServiceRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceconfig.AddServiceConfigRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceversion.AddServiceVersionRequestBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.cluster.ClusterDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.project.ProjectDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.quickstart.AddServiceConfigResponseBodyByQuickStart;
import com.iflytek.ccr.polaris.cynosure.response.quickstart.AddServiceResponseBodyByQuickStart;
import com.iflytek.ccr.polaris.cynosure.response.quickstart.AddVersionResponseBodyByQuickStart;
import com.iflytek.ccr.polaris.cynosure.response.service.ServiceDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.serviceconfig.GrayConfigListDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.serviceconfig.ServiceConfigDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.serviceversion.ServiceVersionDetailResponseBody;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.io.UnsupportedEncodingException;
import java.util.ArrayList;
import java.util.List;

/**
 * 快速开始事务
 *
 * @author sctang2
 * @create 2018-01-29 14:27
 **/
@Service
public class QuickStartTransactional extends BaseService {
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
    private IServiceConfigCondition serviceConfigConditionImpl;

    /**
     * 新增服务
     *
     * @param body
     * @param isLogin
     * @return
     */
    @Transactional
    public Response<AddServiceResponseBodyByQuickStart> addService(AddServiceRequestBodyByQuickStart body, boolean isLogin) {
        String projectId = null;
        String clusterId = null;
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
        String clusterName = body.getCluster();
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

        //根据服务名和集群id查询服务
        String serviceName = body.getService();
        com.iflytek.ccr.polaris.cynosure.domain.Service service = null;
        if (StringUtils.isNotBlank(clusterId)) {
            service = this.serviceConditionImpl.find(serviceName, clusterId);
            if (null != service && isLogin) {
                //已存在该服务
                return new Response<>(SystemErrCode.ERRCODE_SERVICE_EXISTS, SystemErrCode.ERRMSG_SERVICE_EXISTS);
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
        if (null == service) {
            service = this.addService(clusterId, serviceName);
        }

        //创建快速结果
        AddServiceResponseBodyByQuickStart result = new AddServiceResponseBodyByQuickStart();

        //创建项目结果
        ProjectDetailResponseBody projectResult = this.createProjectResult(project);
        result.setProject(projectResult);

        //创建集群结果
        ClusterDetailResponseBody clusterResult = this.createClusterResult(cluster);
        result.setCluster(clusterResult);

        //创建服务结果
        ServiceDetailResponseBody serviceResult = this.createServiceResult(service);
        result.setService(serviceResult);
        return new Response<>(result);
    }
    /**
     * 新增版本
     *
     * @param body
     * @param isLogin
     * @return
     */
    @Transactional
    public Response<AddVersionResponseBodyByQuickStart> addVersion(AddVersionRequestBodyByQuickStart body, boolean isLogin) {
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
        String clusterName = body.getCluster();
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
            service = this.serviceConditionImpl.find(serviceName, clusterId);
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
        String versionName = body.getVersion();
        ServiceVersion version = null;
        if (StringUtils.isNotBlank(serviceId)) {
            version = this.serviceVersionConditionImpl.find(versionName, serviceId);
            if (null != version && isLogin) {
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
        if (null == version) {
            version = this.addServiceVersion(serviceId, versionName, null);
        }

        //创建快速结果
        AddVersionResponseBodyByQuickStart result = new AddVersionResponseBodyByQuickStart();

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
        ServiceVersionDetailResponseBody versionResult = this.createServiceVersionResult(version);
        result.setVersion(versionResult);

        return new Response<>(result);
    }

    /**
     * 新增服务版本(不含配置文件拖拽上传)
     *
     * @param body
     * @return
     */
    @Transactional
    public Response<AddServiceConfigResponseBodyByQuickStart> addServiceVersion(AddServiceVersionRequestBodyByQuickStart body) {
        String projectId = null;
        String clusterId = null;
        String serviceId = null;
        String versionId = null;
        //判断是否为管理员
        boolean isAdmin = this.isAdmin();

        //通过项目名称查询项目信息
        String projectName = body.getProject();
        Project project = this.projectConditionImpl.findByName(projectName);
        if (null == project) {
            if (!isAdmin) {
                //没有权限执行此操作
                return new Response<>(SystemErrCode.ERRCODE_NOT_AUTH, SystemErrCode.ERRMSG_NOT_AUTH);
            }
        } else {
            projectId = project.getId();
        }

        //根据id和集群名称查询集群信息
        String clusterName = body.getCluster();
        Cluster cluster = null;
        if (StringUtils.isNotBlank(projectId)) {
            cluster = this.clusterConditionImpl.find(projectId, clusterName);
            if (null != cluster) {
                clusterId = cluster.getId();
            }
        }

        //根据服务名和集群id查询服务
        String serviceName = body.getService();
        com.iflytek.ccr.polaris.cynosure.domain.Service service = null;
        if (StringUtils.isNotBlank(clusterId)) {
            service = this.serviceConditionImpl.find(serviceName, clusterId);
            if (null != service) {
                serviceId = service.getId();
            }
        }

        //根据版本和服务id查询服务版本
        ServiceVersion serviceVersion = null;
        String version = body.getVersion();
        String desc = body.getDesc();
        if (StringUtils.isNotBlank(serviceId)) {
            serviceVersion = this.serviceVersionConditionImpl.find(version, serviceId);
            if (null != serviceVersion) {
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

        //是否同步配置（直接查询然后放进去，现在改为从body中接收configId）
        //这一部分还得看前端发给我的是不是完整的数组，我是不是要从数组中再查询完整的配置文件信息
        List<ServiceConfig> serviceConfigList = new ArrayList<>();
        if (null != body.getIds() && !body.getIds().isEmpty()) {
            serviceConfigList = this.serviceConfigConditionImpl.findByIds(body.getIds());
        } else {
            serviceConfigList = null;
        }
        //        List<ServiceConfig> serviceConfigList = body.getServiceConfigs();

//        boolean sync = body.isSync();
//        if (sync) {
//            serviceConfigList = this.serviceConfigConditionImpl.findNewList(serviceId);
//        }

        //新增版本
        if (StringUtils.isBlank(versionId)) {
            serviceVersion = this.addServiceVersion(serviceId, version, desc);
            versionId = serviceVersion.getId();
        }

        //批量新增配置（新增服务版本）
        List<ServiceConfig> newServiceConfigList = this.batchVersionAddServiceConfig(versionId, body, serviceConfigList);

        //创建快速结果
        AddServiceConfigResponseBodyByQuickStart result = new AddServiceConfigResponseBodyByQuickStart();

        //创建项目结果
        ProjectDetailResponseBody projectResult = this.createProjectResult(project);
        result.setProject(projectResult);

        //创建集群结果
        ClusterDetailResponseBody clusterResult = this.createClusterResult(cluster);
        result.setCluster(clusterResult);

        //创建服务结果
        ServiceDetailResponseBody serviceResult = this.createServiceResult(service);
        result.setService(serviceResult);

        //创建服务版本结果
        ServiceVersionDetailResponseBody serviceVersionResult = this.createServiceVersionResult(serviceVersion);
        result.setVersion(serviceVersionResult);

        //创建服务版本配置结果
        List<ServiceConfigDetailResponseBody> configs = new ArrayList<>();
        if (null != newServiceConfigList && !newServiceConfigList.isEmpty()) {
            for (ServiceConfig serviceConfig : newServiceConfigList) {
                configs.add(this.createServiceConfigResult(serviceConfig));
            }
        }
        result.setConfigs(configs);
        return new Response<>(result);
    }

    /**
     * 新增服务版本(配置文件拖拽上传)
     *
     * @param body
     * @return
     */
    @Transactional
    public Response<AddServiceConfigResponseBodyByQuickStart> addServiceVersionAndFile(AddServiceVersionRequestBodyByQuickStart body, List<FileContent> fileContentList) {
        String projectId = null;
        String clusterId = null;
        String serviceId = null;
        String versionId = null;
        //判断是否为管理员
        boolean isAdmin = this.isAdmin();

        //通过项目名称查询项目信息
        String projectName = body.getProject();
        Project project = this.projectConditionImpl.findByName(projectName);
        if (null == project) {
            if (!isAdmin) {
                //没有权限执行此操作
                return new Response<>(SystemErrCode.ERRCODE_NOT_AUTH, SystemErrCode.ERRMSG_NOT_AUTH);
            }
        } else {
            projectId = project.getId();
        }

        //根据id和集群名称查询集群信息
        String clusterName = body.getCluster();
        Cluster cluster = null;
        if (StringUtils.isNotBlank(projectId)) {
            cluster = this.clusterConditionImpl.find(projectId, clusterName);
            if (null != cluster) {
                clusterId = cluster.getId();
            }
        }

        //根据服务名和集群id查询服务
        String serviceName = body.getService();
        com.iflytek.ccr.polaris.cynosure.domain.Service service = null;
        if (StringUtils.isNotBlank(clusterId)) {
            service = this.serviceConditionImpl.find(serviceName, clusterId);
            if (null != service) {
                serviceId = service.getId();
            }
        }

        //根据版本和服务id查询服务版本
        ServiceVersion serviceVersion = null;
        String version = body.getVersion();
        String desc = body.getDesc();
        if (StringUtils.isNotBlank(serviceId)) {
            serviceVersion = this.serviceVersionConditionImpl.find(version, serviceId);
            if (null != serviceVersion) {
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

        //是否同步配置（直接查询然后放进去，现在改为从body中接收configId）
        List<ServiceConfig> serviceConfigList;
        if (null != body.getIds() && !body.getIds().isEmpty()) {
            serviceConfigList = this.serviceConfigConditionImpl.findByIds(body.getIds());
        } else {
            serviceConfigList = null;
        }

        //新增版本
        if (StringUtils.isBlank(versionId)) {
            serviceVersion = this.addServiceVersion(serviceId, version, desc);
            versionId = serviceVersion.getId();
        }

        //批量新增配置（新增服务版本）
        List<ServiceConfig> newServiceConfigList = this.batchVersionAddServiceConfigAndFile(versionId, body, fileContentList, serviceConfigList);

        //创建快速结果
        AddServiceConfigResponseBodyByQuickStart result = new AddServiceConfigResponseBodyByQuickStart();

        //创建项目结果
        ProjectDetailResponseBody projectResult = this.createProjectResult(project);
        result.setProject(projectResult);

        //创建集群结果
        ClusterDetailResponseBody clusterResult = this.createClusterResult(cluster);
        result.setCluster(clusterResult);

        //创建服务结果
        ServiceDetailResponseBody serviceResult = this.createServiceResult(service);
        result.setService(serviceResult);

        //创建服务版本结果
        ServiceVersionDetailResponseBody serviceVersionResult = this.createServiceVersionResult(serviceVersion);
        result.setVersion(serviceVersionResult);

        //创建服务版本配置结果
        List<ServiceConfigDetailResponseBody> configs = new ArrayList<>();
        if (null != newServiceConfigList && !newServiceConfigList.isEmpty()) {
            for (ServiceConfig serviceConfig : newServiceConfigList) {
                configs.add(this.createServiceConfigResult(serviceConfig));
            }
        }
        result.setConfigs(configs);
        return new Response<>(result);
    }

    /**
     * 新增服务配置
     *
     * @param body
     * @param fileContentList
     * @return
     */
    @Transactional
    public synchronized Response<AddServiceConfigResponseBodyByQuickStart> addServiceConfig(AddServiceConfigRequestBodyByQuickStart body, List<FileContent> fileContentList) {
        String projectId = null;
        String clusterId = null;
        String serviceId = null;
        String versionId = null;
        String grayId = "0";
        //判断是否为管理员
        boolean isAdmin = this.isAdmin();

        //通过项目名称查询项目信息
        String projectName = body.getProject();
        Project project = this.projectConditionImpl.findByName(projectName);
        if (null == project) {
            if (!isAdmin) {
                //没有权限执行此操作
                return new Response<>(SystemErrCode.ERRCODE_NOT_AUTH, SystemErrCode.ERRMSG_NOT_AUTH);
            }
        } else {
            projectId = project.getId();
        }

        //根据id和集群名称查询集群信息
        String clusterName = body.getCluster();
        Cluster cluster = null;
        if (StringUtils.isNotBlank(projectId)) {
            cluster = this.clusterConditionImpl.find(projectId, clusterName);
            if (null != cluster) {
                clusterId = cluster.getId();
            }
        }

        //根据服务名和集群id查询服务
        String serviceName = body.getService();
        com.iflytek.ccr.polaris.cynosure.domain.Service service = null;
        if (StringUtils.isNotBlank(clusterId)) {
            service = this.serviceConditionImpl.find(serviceName, clusterId);
            if (null != service) {
                serviceId = service.getId();
            }
        }

        //根据版本和服务id查询服务版本
        String versionName = body.getVersion();
        ServiceVersion serviceVersion = null;
        if (StringUtils.isNotBlank(serviceId)) {
            serviceVersion = this.serviceVersionConditionImpl.find(versionName, serviceId);
            if (null != serviceVersion) {
                versionId = serviceVersion.getId();
            }
        }

        //通过版本id，配置名称列表查询服务配置
        List<ServiceConfig> serviceConfigList = null;
        if (StringUtils.isNotBlank(versionId)) {
            List<String> names = new ArrayList<>();
            for (FileContent fileContent : fileContentList) {
                names.add(fileContent.getFileName());
            }
            serviceConfigList = this.serviceConfigConditionImpl.find(versionId, names, grayId);
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
        if (StringUtils.isBlank(versionId)) {
            serviceVersion = this.addServiceVersion(serviceId, versionName, null);
            versionId = serviceVersion.getId();
        }

        //批量新增配置
        List<ServiceConfig> addServiceConfigList = this.batchAddServiceConfig(versionId, body, fileContentList, serviceConfigList);

        //批量更新配置
        List<ServiceConfig> updateServiceConfigList = this.serviceConfigConditionImpl.batchUpdate(body, fileContentList, serviceConfigList);

        //创建快速结果
        AddServiceConfigResponseBodyByQuickStart result = new AddServiceConfigResponseBodyByQuickStart();

        //创建项目结果
        ProjectDetailResponseBody projectResult = this.createProjectResult(project);
        result.setProject(projectResult);

        //创建集群结果
        ClusterDetailResponseBody clusterResult = this.createClusterResult(cluster);
        result.setCluster(clusterResult);

        //创建服务结果
        ServiceDetailResponseBody serviceResult = this.createServiceResult(service);
        result.setService(serviceResult);

        //创建服务版本结果
        ServiceVersionDetailResponseBody serviceVersionResult = this.createServiceVersionResult(serviceVersion);
        result.setVersion(serviceVersionResult);

        //创建服务版本配置结果
        List<ServiceConfigDetailResponseBody> configs = new ArrayList<>();
        if (null != addServiceConfigList && !addServiceConfigList.isEmpty()) {
            for (ServiceConfig serviceConfig : addServiceConfigList) {
                configs.add(this.createServiceConfigResult(serviceConfig));
            }
        }
        if (null != updateServiceConfigList && !updateServiceConfigList.isEmpty()) {
            for (ServiceConfig serviceConfig : updateServiceConfigList) {
                configs.add(this.createServiceConfigResult(serviceConfig));
            }
        }
        result.setConfigs(configs);
        return new Response<>(result);
    }

    /**
     * 新增灰度配置
     *
     * @param body
     * @param fileContentList
     * @return
     */
    @Transactional
    public synchronized Response<GrayConfigListDetailResponseBody> addGrayConfig(AddGrayConfigRequestBody body, List<FileContent> fileContentList) throws UnsupportedEncodingException {
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

        String versionId = body.getVersionId();
        String grayId = body.getGrayId();
        //通过版本id，配置名称列表和灰度组id查询服务配置
        List<ServiceConfig> serviceConfigList = null;
        if (StringUtils.isNotBlank(versionId)) {
            List<String> names = new ArrayList<>();
            for (FileContent fileContent : fileContentList) {
                names.add(fileContent.getFileName());
            }
            serviceConfigList = this.serviceConfigConditionImpl.find(versionId, names, grayId);
        }

        List<FileContent> newFileContents = new ArrayList<>();
        if (null == serviceConfigList || serviceConfigList.isEmpty()) {
            newFileContents.addAll(fileContentList);
        } else {
            for (FileContent fileContent : fileContentList) {
                String name = fileContent.getFileName();
                boolean isExsit = false;
                for (ServiceConfig serviceConfig : serviceConfigList) {
                    String configName = serviceConfig.getName();
                    if (name.equals(configName)) {
                        isExsit = true;
                        break;
                    }
                }
                if (!isExsit) {
                    newFileContents.add(fileContent);
                }
            }
        }
        //批量新增配置
        List<ServiceConfig> addServiceConfigList = this.serviceConfigConditionImpl.batchAddGrayFile(versionId, body, newFileContents);

        //批量更新配置
        List<ServiceConfig> updateServiceConfigList = this.serviceConfigConditionImpl.batchUpdateGrayConfig(body, fileContentList, serviceConfigList);

        //创建灰度配置结果
        List<ServiceConfig> newServiceConfigList = new ArrayList<>();

        if (null != addServiceConfigList && !addServiceConfigList.isEmpty()) {
            for (ServiceConfig serviceConfig : addServiceConfigList) {
                newServiceConfigList.add(serviceConfig);
            }
        }
        if (null != updateServiceConfigList && !updateServiceConfigList.isEmpty()) {
            for (ServiceConfig serviceConfig : updateServiceConfigList) {
                newServiceConfigList.add(serviceConfig);
            }
        }

        GrayConfigListDetailResponseBody result = this.createGrayConfigResult(newServiceConfigList);
        return new Response<>(result);
    }

    /**
     * 批量新增配置（在新增服务版本时，无文件上传）
     *
     * @param versionId
     * @param addServiceVersion
     * @param serviceConfigList
     * @return
     */
    private List<ServiceConfig> batchVersionAddServiceConfig(String versionId, AddServiceVersionRequestBodyByQuickStart addServiceVersion, List<ServiceConfig> serviceConfigList) {
        AddServiceConfigRequestBody body = new AddServiceConfigRequestBody();
        body.setVersionId(versionId);
        AddServiceConfigRequestBodyByQuickStart addServiceConfig = new AddServiceConfigRequestBodyByQuickStart();
        addServiceConfig.setProject(addServiceVersion.getProject());
        addServiceConfig.setCluster(addServiceVersion.getCluster());
        addServiceConfig.setService(addServiceVersion.getService());
        addServiceConfig.setVersion(addServiceVersion.getVersion());
        return this.serviceConfigConditionImpl.batchVersionAdd(body, addServiceConfig, serviceConfigList);
    }

    /**
     * 批量新增配置(含拖拽文件上传)
     *
     * @param versionId
     * @param addServiceVersion
     * @param serviceConfigList
     * @return
     */
    private List<ServiceConfig> batchVersionAddServiceConfigAndFile(String versionId, AddServiceVersionRequestBodyByQuickStart addServiceVersion, List<FileContent> fileContentList, List<ServiceConfig> serviceConfigList) {
        List<FileContent> newFileContents = new ArrayList<>();
        if (null == serviceConfigList || serviceConfigList.isEmpty()) {
            newFileContents.addAll(fileContentList);
//            serviceConfigList.forEach(x -> {
//                FileContent fileContent = new FileContent();
//                fileContent.setFileName(x.getName());
//                fileContent.setContent(x.getContent());
//                fileContentList.add(fileContent);
//            });
        }

        AddServiceConfigRequestBody body = new AddServiceConfigRequestBody();
        body.setVersionId(versionId);
        AddServiceConfigRequestBodyByQuickStart addServiceConfig = new AddServiceConfigRequestBodyByQuickStart();
        addServiceConfig.setProject(addServiceVersion.getProject());
        addServiceConfig.setCluster(addServiceVersion.getCluster());
        addServiceConfig.setService(addServiceVersion.getService());
        addServiceConfig.setVersion(addServiceVersion.getVersion());
        return this.serviceConfigConditionImpl.batchVersionAddAndFile(body, addServiceConfig, serviceConfigList, fileContentList);
    }

    /**
     * 批量新增配置
     *
     * @param versionId
     * @param addServiceConfig
     * @param fileContentList
     * @param serviceConfigList
     * @return
     */
    private List<ServiceConfig> batchAddServiceConfig(String versionId, AddServiceConfigRequestBodyByQuickStart addServiceConfig, List<FileContent> fileContentList, List<ServiceConfig> serviceConfigList) {
        List<FileContent> newFileContents = new ArrayList<>();
        if (null == serviceConfigList || serviceConfigList.isEmpty()) {
            newFileContents.addAll(fileContentList);
        } else {
            for (FileContent fileContent : fileContentList) {
                String name = fileContent.getFileName();
                boolean isExsit = false;
                for (ServiceConfig serviceConfig : serviceConfigList) {
                    String configName = serviceConfig.getName();
                    if (name.equals(configName)) {
                        isExsit = true;
                        break;
                    }
                }
                if (!isExsit) {
                    newFileContents.add(fileContent);
                }
            }
        }
        AddServiceConfigRequestBody body = new AddServiceConfigRequestBody();
        body.setVersionId(versionId);
        return this.serviceConfigConditionImpl.batchAdd(body, addServiceConfig, newFileContents);
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
     * @param version
     * @param desc
     * @return
     */
    private ServiceVersion addServiceVersion(String serviceId, String version, String desc) {
        AddServiceVersionRequestBody body = new AddServiceVersionRequestBody();
        body.setServiceId(serviceId);
        body.setVersion(version);
        body.setDesc(desc);
        return this.serviceVersionConditionImpl.add(body);
    }

    /**
     * 创建灰度配置结果
     *
     * @param serviceConfigs
     * @return
     */
    protected GrayConfigListDetailResponseBody createGrayConfigResult(List<ServiceConfig> serviceConfigs) {
        GrayConfigListDetailResponseBody result = new GrayConfigListDetailResponseBody();
        result.setServiceConfigList(serviceConfigs);
        return result;
    }
}
