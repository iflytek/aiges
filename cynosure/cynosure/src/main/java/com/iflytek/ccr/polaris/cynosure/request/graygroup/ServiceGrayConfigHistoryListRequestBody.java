package com.iflytek.ccr.polaris.cynosure.request.graygroup;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.request.BaseRequestBody;
import org.hibernate.validator.constraints.Length;
import org.hibernate.validator.constraints.NotBlank;

import javax.validation.constraints.NotNull;
import java.io.Serializable;

/**
 * 查询服务配置历史列表-请求
 *
 * @author sctang2
 * @create 2017-12-12 15:57
 **/
public class ServiceGrayConfigHistoryListRequestBody extends BaseRequestBody implements Serializable {
    private static final long serialVersionUID = 4884803719607172202L;

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

    //灰度组名称
    @NotBlank(message = SystemErrCode.ERRMSG_GRAY_GROUP_NAME_NOT_NULL)
    @Length(min = 1, max = 100, message = SystemErrCode.ERRMSG_GRAY_GROUP_NAME_MAX_LENGTH)
    @NotNull(message = SystemErrCode.ERRMSG_GRAY_GROUP_NAME_NOT_NULL)
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

    @Override
    public String toString() {
        return "ServiceConfigHistoryListRequestBody{" +
                "project='" + project + '\'' +
                ", cluster='" + cluster + '\'' +
                ", service='" + service + '\'' +
                ", version='" + version + '\'' +
                ", gray='" + gray + '\'' +
                '}';
    }
}
