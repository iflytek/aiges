package com.iflytek.ccr.polaris.cynosure.request.graygroup;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import org.hibernate.validator.constraints.Length;
import org.hibernate.validator.constraints.NotBlank;

import java.io.Serializable;

/**
 * 删除灰度组-请求
 *
 * @author sctang2
 * @create 2017-11-21 14:21
 **/
public class DeleteGrayGroupRequestBody implements Serializable {
    private static final long serialVersionUID = 4229119734922797764L;

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

    public String getGrayId() {
        return grayId;
    }

    public void setGrayId(String grayId) {
        this.grayId = grayId;
    }

    @Override
    public String toString() {
        return "DeleteGrayGroupRequestBody{" +
                "project='" + project + '\'' +
                ", cluster='" + cluster + '\'' +
                ", service='" + service + '\'' +
                ", version='" + version + '\'' +
                ", grayId='" + grayId + '\'' +
                '}';
    }
}
