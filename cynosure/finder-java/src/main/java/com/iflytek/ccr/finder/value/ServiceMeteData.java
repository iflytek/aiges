package com.iflytek.ccr.finder.value;

/**
 * 服务元数据
 */
public class ServiceMeteData {

    /**
     * 项目名称
     */
    private String project;

    /**
     * 集群名称
     */
    private String group;

    /**
     * 服务名称
     */
    private String service;

    /**
     *  组件版本
     */
    private String version;

    /**
     *  组件唯一标识（一般场景推荐ip：port）
     */
    private String address;

    public ServiceMeteData(String project, String group, String service, String version, String address) {
        this.project = project;
        this.group = group;
        this.service = service;
        this.version = version;
        this.address = address;
    }

    public String getAddress() {
        return address;
    }

    public void setAddress(String address) {
        this.address = address;
    }

    public String getProject() {
        return project;
    }

    public void setProject(String project) {
        this.project = project;
    }

    public String getGroup() {
        return group;
    }

    public void setGroup(String group) {
        this.group = group;
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
        return "ServiceMeteData{" +
                "project='" + project + '\'' +
                ", group='" + group + '\'' +
                ", service='" + service + '\'' +
                ", version='" + version + '\'' +
                ", address='" + address + '\'' +
                '}';
    }
}
