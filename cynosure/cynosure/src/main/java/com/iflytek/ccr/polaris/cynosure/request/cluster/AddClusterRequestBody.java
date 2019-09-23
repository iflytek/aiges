package com.iflytek.ccr.polaris.cynosure.request.cluster;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import org.hibernate.validator.constraints.Length;
import org.hibernate.validator.constraints.NotBlank;

import java.io.Serializable;

/**
 * 添加集群-请求
 *
 * @author sctang2
 * @create 2017-11-16 9:03
 **/
public class AddClusterRequestBody implements Serializable {
    private static final long serialVersionUID = -5132518697884366268L;

    //集群名称
    @NotBlank(message = SystemErrCode.ERRMSG_CLUSTER_NAME_NOT_NULL)
    @Length(min = 1, max = 100, message = SystemErrCode.ERRMSG_CLUSTER_NAME_MAX_LENGTH)
    private String name;

    //集群描述
    @Length(max = 500, message = SystemErrCode.ERRMSG_CLUSTER_DESC_MAX_LENGTH)
    private String desc;

    //项目id
    @NotBlank(message = SystemErrCode.ERRMSG_PROJECT_ID_NOT_NULL)
    @Length(max = 50, message = SystemErrCode.ERRMSG_PROJECT_ID_MAX_LENGTH)
    private String projectId;

    public AddClusterRequestBody() {
    }

    public AddClusterRequestBody(String name, String desc, String projectId) {
        this.name = name;
        this.desc = desc;
        this.projectId = projectId;
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

    public String getProjectId() {
        return projectId;
    }

    public void setProjectId(String projectId) {
        this.projectId = projectId;
    }

    @Override
    public String toString() {
        return "AddClusterRequestBody{" +
                "name='" + name + '\'' +
                ", desc='" + desc + '\'' +
                ", projectId='" + projectId + '\'' +
                '}';
    }
}
