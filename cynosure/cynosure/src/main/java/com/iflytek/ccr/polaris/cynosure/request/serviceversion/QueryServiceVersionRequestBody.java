package com.iflytek.ccr.polaris.cynosure.request.serviceversion;

import com.iflytek.ccr.polaris.cynosure.request.BaseRequestBody;

import java.io.Serializable;

/**
 * 查询服务版本列表-请求
 *
 * @author sctang2
 * @create 2017-11-21 10:02
 **/
public class QueryServiceVersionRequestBody extends BaseRequestBody implements Serializable {
    private static final long serialVersionUID = 4223996556386317491L;

    //项目名称
    private String project;

    //集群名称
    private String cluster;

    //服务名称
    private String service;

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

    @Override
    public String toString() {
        return "QueryServiceVersionRequestBody{" +
                "project='" + project + '\'' +
                ", cluster='" + cluster + '\'' +
                ", service='" + service + '\'' +
                '}';
    }
}
