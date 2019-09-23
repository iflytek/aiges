package com.iflytek.ccr.polaris.cynosure.response.servicediscovery;

import com.iflytek.ccr.polaris.cynosure.response.cluster.ClusterDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.project.ProjectDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.service.ServiceDetailResponseBody;

import java.io.Serializable;

/**
 * 快速开始创建服务-响应
 *
 * @author sctang2
 * @create 2018-01-29 12:04
 **/
public class AddApiVersionResponseBody implements Serializable {
    private static final long serialVersionUID = 7025929839915712530L;

    //项目
    private ProjectDetailResponseBody project;

    //集群
    private ClusterDetailResponseBody cluster;

    //服务
    private ServiceDetailResponseBody service;

    //版本
    private ServiceApiVersionDetailResponseBody apiVersion;

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

    public ServiceApiVersionDetailResponseBody getApiVersion() {
        return apiVersion;
    }

    public void setApiVersion(ServiceApiVersionDetailResponseBody apiVersion) {
        this.apiVersion = apiVersion;
    }

    @Override
    public String toString() {
        return "AddApiVersionResponseBody{" +
                "project=" + project +
                ", cluster=" + cluster +
                ", service=" + service +
                ", apiVersion=" + apiVersion +
                '}';
    }
}
