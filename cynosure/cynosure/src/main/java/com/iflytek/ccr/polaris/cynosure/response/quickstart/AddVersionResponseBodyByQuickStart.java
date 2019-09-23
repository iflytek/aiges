package com.iflytek.ccr.polaris.cynosure.response.quickstart;

import com.iflytek.ccr.polaris.cynosure.response.cluster.ClusterDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.project.ProjectDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.service.ServiceDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.serviceversion.ServiceVersionDetailResponseBody;

import java.io.Serializable;

/**
 * 快速开始创建服务-响应
 *
 * @author sctang2
 * @create 2018-01-29 12:04
 **/
public class AddVersionResponseBodyByQuickStart implements Serializable {
    private static final long serialVersionUID = 7025929839915712530L;

    //项目
    private ProjectDetailResponseBody project;

    //集群
    private ClusterDetailResponseBody cluster;

    //服务
    private ServiceDetailResponseBody service;

    //版本
    private ServiceVersionDetailResponseBody version;

    public ProjectDetailResponseBody getProject() {
        return project;
    }

    public void setProject(ProjectDetailResponseBody project) {
        this.project = project;
    }

    public ClusterDetailResponseBody getCluster() {
        return cluster;
    }

    public void setCluster(ClusterDetailResponseBody cluster) {
        this.cluster = cluster;
    }

    public ServiceDetailResponseBody getService() {
        return service;
    }

    public void setService(ServiceDetailResponseBody service) {
        this.service = service;
    }

    public ServiceVersionDetailResponseBody getVersion() {
        return version;
    }

    public void setVersion(ServiceVersionDetailResponseBody version) {
        this.version = version;
    }

    @Override
    public String toString() {
        return "AddVersionResponseBodyByQuickStart{" +
                "project=" + project +
                ", cluster=" + cluster +
                ", service=" + service +
                ", version=" + version +
                '}';
    }
}
