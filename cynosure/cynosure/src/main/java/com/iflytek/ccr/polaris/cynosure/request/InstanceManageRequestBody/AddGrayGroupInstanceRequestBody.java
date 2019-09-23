package com.iflytek.ccr.polaris.cynosure.request.InstanceManageRequestBody;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.request.BaseRequestBody;
import org.hibernate.validator.constraints.Length;
import org.hibernate.validator.constraints.NotBlank;

import java.io.Serializable;

/**
 * 查看订阅实例-请求
 *
 * @author sctang2
 * @create 2017-11-21 14:21
 **/
public class AddGrayGroupInstanceRequestBody extends BaseRequestBody implements Serializable {
    private static final long serialVersionUID = 4210119734922797764L;

    //项目名称
    @NotBlank(message = SystemErrCode.ERRMSG_PROJECT_NAME_NOT_NULL)
    private String project;

    //集群名称
    @NotBlank(message = SystemErrCode.ERRMSG_CLUSTER_NAME_NOT_NULL)
    private String cluster;

    //服务名
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_NAME_NOT_NULL)
    private String service;

    //服务版本名称
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_VERSION_NOT_NULL)
    private String version;

    //版本id
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_VERSION_ID_NOT_NULL)
    @Length(max = 50, message = SystemErrCode.ERRMSG_SERVICE_VERSION_ID_MAX_LENGTH)
    private String versionId;

    //灰度组id
    @NotBlank(message = SystemErrCode.ERRMSG_GRAY_GROUP_ID_NOT_NULL)
    @Length(max = 50, message = SystemErrCode.ERRMSG_GRAY_GROUP_ID_MAX_LENGTH)
    private String grayId;

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

    public String getVersionId() {
        return versionId;
    }

    public void setVersionId(String versionId) {
        this.versionId = versionId;
    }

    public String getGrayId() {
        return grayId;
    }

    public void setGrayId(String grayId) {
        this.grayId = grayId;
    }

    @Override
    public String toString() {
        return "AddGrayGroupInstanceRequestBody{" +
                "project='" + project + '\'' +
                ", cluster='" + cluster + '\'' +
                ", service='" + service + '\'' +
                ", version='" + version + '\'' +
                ", versionId='" + versionId + '\'' +
                ", grayId='" + grayId + '\'' +
                '}';
    }
}
