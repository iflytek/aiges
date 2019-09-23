package com.iflytek.ccr.polaris.cynosure.request.servicediscovery;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import org.hibernate.validator.constraints.NotBlank;

import javax.validation.constraints.Max;
import javax.validation.constraints.Min;
import java.io.Serializable;
import java.util.ArrayList;
import java.util.List;

/**
 * 编辑服务发现提供者-请求
 *
 * @author sctang2
 * @create 2017-12-08 9:08
 **/
public class EditServiceDiscoveryProviderRequestBody implements Serializable {
    private static final long serialVersionUID = 8405554706717834060L;

    //项目名称
    @NotBlank(message = SystemErrCode.ERRMSG_PROJECT_NAME_NOT_NULL)
    private String project;

    //集群名称
    @NotBlank(message = SystemErrCode.ERRMSG_CLUSTER_NAME_NOT_NULL)
    private String cluster;

    //服务名称
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_NAME_NOT_NULL)
    private String service;

    //服务版本
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_VERSION_NOT_NULL)
    private String apiVersion;

    //区域名称
    @NotBlank(message = SystemErrCode.ERRMSG_REGION_NAME_NOT_NULL)
    private String region;

    //地址
    @NotBlank(message = SystemErrCode.ERRMSG_SERVICE_CONFIG_ADDR_NOT_NULL)
    private String addr;

    //状态
    private boolean valid;

    //权重
    @Min(value = 0, message = SystemErrCode.ERRMSG_SERVICE_DISCOVERY_WEIGHT_INVALID)
    @Max(value = 100, message = SystemErrCode.ERRMSG_SERVICE_DISCOVERY_WEIGHT_INVALID)
    private int weight;

    //自定义key，value参数
    private List<ServiceParam> params = new ArrayList<>();

    public List<ServiceParam> getParams() {
        return params;
    }

    public void setParams(List<ServiceParam> params) {
        this.params = params;
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

    public String getRegion() {
        return region;
    }

    public void setRegion(String region) {
        this.region = region;
    }

    public String getAddr() {
        return addr;
    }

    public void setAddr(String addr) {
        this.addr = addr;
    }

    public boolean isValid() {
        return valid;
    }

    public void setValid(boolean valid) {
        this.valid = valid;
    }

    public int getWeight() {
        return weight;
    }

    public void setWeight(int weight) {
        this.weight = weight;
    }

    public String getApiVersion() {
        return apiVersion;
    }

    public void setApiVersion(String apiVersion) {
        this.apiVersion = apiVersion;
    }

    @Override
    public String toString() {
        return "EditServiceDiscoveryProviderRequestBody{" +
                "project='" + project + '\'' +
                ", cluster='" + cluster + '\'' +
                ", service='" + service + '\'' +
                ", apiVersion='" + apiVersion + '\'' +
                ", region='" + region + '\'' +
                ", addr='" + addr + '\'' +
                ", valid=" + valid +
                ", weight=" + weight +
                '}';
    }
}
