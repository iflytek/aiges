package com.iflytek.ccr.polaris.companion.common;

/**
 * 反馈对象
 */
public class ServiceValue {

    private String project;
    private String group;
    private String service;
    private String apiVersion;


    public String getApiVersion() {
        return apiVersion;
    }

    public void setApiVersion(String apiVersion) {
        this.apiVersion = apiVersion;
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

    @Override
    public String toString() {
        return "ServiceValue{" +
                "project='" + project + '\'' +
                ", group='" + group + '\'' +
                ", service='" + service + '\'' +
                ", apiVersion='" + apiVersion + '\'' +
                '}';
    }
}

