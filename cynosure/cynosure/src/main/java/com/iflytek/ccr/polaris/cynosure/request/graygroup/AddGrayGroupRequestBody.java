package com.iflytek.ccr.polaris.cynosure.request.graygroup;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import org.hibernate.validator.constraints.Length;
import org.hibernate.validator.constraints.NotBlank;

import java.io.Serializable;
import java.util.List;

/**
 * 新增灰度组-请求
 *
 * @author sctang2
 * @create 2017-11-21 14:21
 **/
public class AddGrayGroupRequestBody implements Serializable {
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

    //名称
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

    //当前版本正常配置id列表
    private List<String> ids;

    //推送实例列表
    private String content;

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

    public List<String> getIds() {
        return ids;
    }

    public void setIds(List<String> ids) {
        this.ids = ids;
    }

    public String getContent() {
        return content;
    }

    public void setContent(String content) {
        this.content = content;
    }

    @Override
    public String toString() {
        return "AddGrayGroupRequestBody{" +
                "project='" + project + '\'' +
                ", cluster='" + cluster + '\'' +
                ", service='" + service + '\'' +
                ", version='" + version + '\'' +
                ", name='" + name + '\'' +
                ", desc='" + desc + '\'' +
                ", versionId='" + versionId + '\'' +
                ", ids=" + ids +
                ", content='" + content + '\'' +
                '}';
    }
}
