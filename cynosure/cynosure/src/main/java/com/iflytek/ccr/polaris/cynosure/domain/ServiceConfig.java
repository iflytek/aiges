package com.iflytek.ccr.polaris.cynosure.domain;

import java.util.Date;

/**
 * 服务配置模型
 *
 * @author sctang2
 * @create 2017-11-21 11:32
 **/
public class ServiceConfig {
    //配置id
    private String id;

    //用户id
    private String userId;

    //版本id
    private String versionId;

    //灰度组id
    private String grayId;

    //配置名称
    private String name;

    //配置路径
    private String path;

    //配置描述
    private String description;

    //创建时间
    private Date createTime;

    //更新时间
    private Date updateTime;

    //配置内容
    private byte[] content;

    //配置内容md5
    private String md5;

    //服务版本
    private ServiceVersion serviceVersion;

    //服务
    private Service service;

    //集群
    private Cluster cluster;

    //项目
    private Project project;

//    //灰度组
//    private String grayGroupId;

//    public String getGrayGroupId() {
//        return grayGroupId;
//    }
//
//    public void setGrayGroupId(String grayGroupId) {
//        this.grayGroupId = grayGroupId;
//    }

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getUserId() {
        return userId;
    }

    public void setUserId(String userId) {
        this.userId = userId;
    }

    public String getVersionId() {
        return versionId;
    }

    public void setVersionId(String versionId) {
        this.versionId = versionId;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getPath() {
        return path;
    }

    public void setPath(String path) {
        this.path = path;
    }

    public String getDescription() {
        return description;
    }

    public void setDescription(String description) {
        this.description = description;
    }

    public Date getCreateTime() {
        return createTime;
    }

    public void setCreateTime(Date createTime) {
        this.createTime = createTime;
    }

    public Date getUpdateTime() {
        return updateTime;
    }

    public void setUpdateTime(Date updateTime) {
        this.updateTime = updateTime;
    }

    public byte[] getContent() {
        return content;
    }

    public void setContent(byte[] content) {
        this.content = content;
    }

    public ServiceVersion getServiceVersion() {
        return serviceVersion;
    }

    public void setServiceVersion(ServiceVersion serviceVersion) {
        this.serviceVersion = serviceVersion;
    }

    public Service getService() {
        return service;
    }

    public void setService(Service service) {
        this.service = service;
    }

    public Cluster getCluster() {
        return cluster;
    }

    public void setCluster(Cluster cluster) {
        this.cluster = cluster;
    }

    public Project getProject() {
        return project;
    }

    public void setProject(Project project) {
        this.project = project;
    }

    public String getMd5() {
        return md5;
    }

    public void setMd5(String md5) {
        this.md5 = md5;
    }

    public String getGrayId() {
        return grayId;
    }

    public void setGrayId(String grayId) {
        this.grayId = grayId;
    }
}
