package com.iflytek.ccr.polaris.cynosure.request.track;

import com.iflytek.ccr.polaris.cynosure.request.BaseRequestBody;

import java.io.Serializable;

/**
 * 查询轨迹列表-请求
 *
 * @author sctang2
 * @create 2018-01-03 14:06
 **/
public class QueryTrackRequestBody extends BaseRequestBody implements Serializable {
    private static final long serialVersionUID = -6245623164289908346L;

    //项目名称
    private String project;

    //集群名称
    private String cluster;

    //服务名称
    private String service;

    //服务版本名称
    private String version;

    //灰度推送过滤字段，取值为0,1,-1
    private Integer filterGray;

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

    public Integer getFilterGray() {
        return filterGray;
    }

    public void setFilterGray(Integer filterGray) {
        this.filterGray = filterGray;
    }

    @Override
    public String toString() {
        return "QueryTrackRequestBody{" +
                "project='" + project + '\'' +
                ", cluster='" + cluster + '\'' +
                ", service='" + service + '\'' +
                ", version='" + version + '\'' +
                '}';
    }
}
