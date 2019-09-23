package com.iflytek.ccr.polaris.cynosure.request.graygroup;

import com.iflytek.ccr.polaris.cynosure.request.BaseRequestBody;

import java.io.Serializable;

/**
 * 查询灰度组列表-请求
 *
 * @author sctang2
 * @create 2017-11-21 10:02
 **/
public class QueryGrayGroupRequestBody extends BaseRequestBody implements Serializable {
    private static final long serialVersionUID = 4223996556306317491L;

    //项目名称
    private String project;

    //集群名称
    private String cluster;

    //服务名称
    private String service;

    //服务版本
    private String version;

    public String getProject() {
        return project;
    }

    public void setProject(String project) {
        this.project = project;
    }

    public String getCluster() {
        return cluster;
    }

    public void setCluster(String cluster) {
        this.cluster = cluster;
    }

    public String getService() {
        return service;
    }

    public void setService(String service) {
        this.service = service;
    }

    public String getVersion() {
        return version;
    }

    public void setVersion(String version) {
        this.version = version;
    }

    @Override
    public String toString() {
        return "QueryGrayGroupRequestBody{" +
                "project='" + project + '\'' +
                ", cluster='" + cluster + '\'' +
                ", service='" + service + '\'' +
                ", version='" + version + '\'' +
                '}';
    }
}
