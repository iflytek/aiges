package com.iflytek.ccr.polaris.cynosure.request.cluster;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import org.hibernate.validator.constraints.Length;
import org.hibernate.validator.constraints.NotBlank;

import java.io.Serializable;

/**
 * 复制集群-请求
 *
 * @author sctang2
 * @create 2017-11-17 16:04
 **/
public class CopyClusterRequestBody implements Serializable {
    private static final long serialVersionUID = 2388285789083489433L;

    //复制出的集群所归属的项目的id
    @NotBlank(message = SystemErrCode.ERRMSG_PROJECT_ID_NOT_NULL)
    @Length(max = 50, message = SystemErrCode.ERRMSG_PROJECT_ID_MAX_LENGTH)
    private String projectId;

    //被复制的集群的id
    @NotBlank(message = SystemErrCode.ERRMSG_CLUSTER_ID_NOT_NULL)
    @Length(max = 50, message = SystemErrCode.ERRMSG_CLUSTER_ID_MAX_LENGTH)
    private String oldClusterId;

    //新集群的名称
    @NotBlank(message = SystemErrCode.ERRMSG_CLUSTER_NAME_NOT_NULL)
    @Length(min = 1, max = 20, message = SystemErrCode.ERRMSG_SERVICE_VERSION_MAX_LENGTH)
    private String clusterName;

    //集群描述
    @Length(max = 500, message = SystemErrCode.ERRMSG_CLUSTER_DESC_MAX_LENGTH)
    private String desc;


    public String getClusterName() {
        return clusterName;
    }

    public void setClusterName(String clusterName) {
        this.clusterName = clusterName;
    }

    public String getDesc() {
        return desc;
    }

    public void setDesc(String desc) {
        this.desc = desc;
    }

    public String getOldClusterId() {
        return oldClusterId;
    }

    public void setOldClusterId(String oldClusterId) {
        this.oldClusterId = oldClusterId;
    }

    public String getProjectId() {
        return projectId;
    }

    public void setProjectId(String projectId) {
        this.projectId = projectId;
    }

    @Override
    public String toString() {
        return "CopyClusterRequestBody{" +
                "clusterName='" + clusterName + '\'' +
                ", desc='" + desc + '\'' +
                ", oldClusterId='" + oldClusterId + '\'' +
                ", projectId='" + projectId + '\'' +
                '}';
    }
}
