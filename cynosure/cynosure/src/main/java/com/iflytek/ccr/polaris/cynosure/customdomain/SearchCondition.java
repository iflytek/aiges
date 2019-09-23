package com.iflytek.ccr.polaris.cynosure.customdomain;

/**
 * 搜索条件
 *
 * @author sctang2
 * @create 2018-01-26 11:08
 **/
public class SearchCondition {
    //项目名称
    private String project;

    //集群名称
    private String cluster;

    //服务名称
    private String service;

    //版本名称
    private String version;

    //灰度组名称
    private String gray;

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

    public String getGray() {
        return gray;
    }

    public void setGray(String gray) {
        this.gray = gray;
    }
}
