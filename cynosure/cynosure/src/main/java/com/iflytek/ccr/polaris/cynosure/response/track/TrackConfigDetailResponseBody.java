package com.iflytek.ccr.polaris.cynosure.response.track;

import java.io.Serializable;

/**
 * 轨迹配置明细-响应
 *
 * @author sctang2
 * @create 2017-11-25 10:24
 **/
public class TrackConfigDetailResponseBody implements Serializable {
    private static final long serialVersionUID = 7357891616914444028L;

    //反馈id
    private String id;

    //项目名称
    private String project;

    //服务组名称
    private String cluster;

    //服务名称
    private String service;

    //版本号
    private String version;

    //配置名称
    private String config;

    //地址
    private String addr;

    //更新状态
    private int updateStatus;

    //更新时间
    private long updateTime;

    //加载状态
    private int loadStatus;

    //加载时间
    private long loadTime;

    //灰度标识
    private String grayGroupId;

    //所属灰度组名字
    private String grayGroupName;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

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

    public String getConfig() {
        return config;
    }

    public void setConfig(String config) {
        this.config = config;
    }

    public String getAddr() {
        return addr;
    }

    public void setAddr(String addr) {
        this.addr = addr;
    }

    public int getUpdateStatus() {
        return updateStatus;
    }

    public void setUpdateStatus(int updateStatus) {
        this.updateStatus = updateStatus;
    }

    public long getUpdateTime() {
        return updateTime;
    }

    public void setUpdateTime(long updateTime) {
        this.updateTime = updateTime;
    }

    public int getLoadStatus() {
        return loadStatus;
    }

    public void setLoadStatus(int loadStatus) {
        this.loadStatus = loadStatus;
    }

    public long getLoadTime() {
        return loadTime;
    }

    public void setLoadTime(long loadTime) {
        this.loadTime = loadTime;
    }

    public String getGrayGroupId() {
        return grayGroupId;
    }

    public void setGrayGroupId(String grayGroupId) {
        this.grayGroupId = grayGroupId;
    }

    public String getGrayGroupName() {
        return grayGroupName;
    }

    public void setGrayGroupName(String grayGroupName) {
        this.grayGroupName = grayGroupName;
    }

    @Override
    public String toString() {
        return "TrackConfigDetailResponseBody{" +
                "id='" + id + '\'' +
                ", project='" + project + '\'' +
                ", cluster='" + cluster + '\'' +
                ", service='" + service + '\'' +
                ", version='" + version + '\'' +
                ", config='" + config + '\'' +
                ", addr='" + addr + '\'' +
                ", updateStatus=" + updateStatus +
                ", updateTime=" + updateTime +
                ", loadStatus=" + loadStatus +
                ", loadTime=" + loadTime +
                '}';
    }
}
