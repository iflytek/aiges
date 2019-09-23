package com.iflytek.ccr.polaris.cynosure.request.graygroup;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import org.hibernate.validator.constraints.Length;
import org.hibernate.validator.constraints.NotBlank;

import java.io.Serializable;

/**
 * 新增灰度配置请求-明细
 * Created by DELL-5490 on 2018/7/11.
 */
public class AddGrayConfigRequestBody implements Serializable {
    private static final long serialVersionUID = 4219119734922794764L;

    //项目名称
    @NotBlank(message = SystemErrCode.ERRMSG_PROJECT_NAME_NOT_NULL)
    @Length(min = 1, max = 100, message = SystemErrCode.ERRMSG_PROJECT_NAME_MAX_LENGTH)
    private String project;

    //集群名称
    @NotBlank(message = SystemErrCode.ERRMSG_CLUSTER_NAME_NOT_NULL)
    @Length(min = 1, max = 100, message = SystemErrCode.ERRMSG_CLUSTER_NAME_MAX_LENGTH)
    private String cluster;

    //服务名称
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_NAME_NOT_NULL)
    @Length(min = 1, max = 100, message = SystemErrCode.ERRMSG_SERVICE_NAME_MAX_LENGTH)
    private String service;

    //服务版本
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_VERSION_NOT_NULL)
    @Length(min = 1, max = 20, message = SystemErrCode.ERRMSG_SERVICE_VERSION_MAX_LENGTH)
    private String version;

    //灰度组id
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_VERSION_ID_NOT_NULL)
    @Length(max = 50, message = SystemErrCode.ERRMSG_SERVICE_VERSION_ID_MAX_LENGTH)
    private String grayId;

    //灰度组名称
    @NotBlank(message = SystemErrCode.ERRMSG_GRAY_GROUP_NAME_NOT_NULL)
    @Length(max = 100, message = SystemErrCode.ERRMSG_GRAY_GROUP_NAME_MAX_LENGTH)
    private String name;

    //描述
    @Length(max = 500, message = SystemErrCode.ERRMSG_SERVICE_CONFIG_DESC_MAX_LENGTH)
    private String desc;

    //版本id
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_VERSION_ID_NOT_NULL)
    @Length(max = 50, message = SystemErrCode.ERRMSG_SERVICE_VERSION_ID_MAX_LENGTH)
    private String versionId;

    public String getGrayId() {
        return grayId;
    }

    public void setGrayId(String grayId) {
        this.grayId = grayId;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getDesc() {
        return desc;
    }

    public void setDesc(String desc) {
        this.desc = desc;
    }

    public String getVersionId() {
        return versionId;
    }

    public void setVersionId(String versionId) {
        this.versionId = versionId;
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

    @Override
    public String toString() {
        return "AddGrayConfigRequestBody{" +
                "project='" + project + '\'' +
                ", cluster='" + cluster + '\'' +
                ", service='" + service + '\'' +
                ", version='" + version + '\'' +
                ", grayId='" + grayId + '\'' +
                ", name='" + name + '\'' +
                ", desc='" + desc + '\'' +
                ", versionId='" + versionId + '\'' +
                '}';
    }
}
