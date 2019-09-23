package com.iflytek.ccr.polaris.cynosure.response.quickstart;

import com.iflytek.ccr.polaris.cynosure.response.cluster.ClusterDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.project.ProjectDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.service.ServiceDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.serviceconfig.ServiceConfigDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.serviceversion.ServiceVersionDetailResponseBody;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

import java.io.Serializable;
import java.util.List;

/**
 * 新增服务配置-响应
 *
 * @author sctang2
 * @create 2018-01-30 9:13
 **/
@ApiModel("新增服务配置响应")
public class AddServiceConfigResponseBodyByQuickStart implements Serializable {
	private static final long serialVersionUID = -5769822887741353548L;

	// 项目
	@ApiModelProperty("项目")
	private ProjectDetailResponseBody project;

	// 集群
	@ApiModelProperty("集群")
	private ClusterDetailResponseBody cluster;

	// 服务
	@ApiModelProperty("服务")
	private ServiceDetailResponseBody service;

	// 服务版本
	@ApiModelProperty("服务版本")
	private ServiceVersionDetailResponseBody version;

	// 服务配置
	@ApiModelProperty("服务配置")
	private List<ServiceConfigDetailResponseBody> configs;

	public ProjectDetailResponseBody getProject() {
		return project;
	}

	public void setProject(ProjectDetailResponseBody project) {
		this.project = project;
	}

	public ClusterDetailResponseBody getCluster() {
		return cluster;
	}

	public void setCluster(ClusterDetailResponseBody cluster) {
		this.cluster = cluster;
	}

	public ServiceDetailResponseBody getService() {
		return service;
	}

	public void setService(ServiceDetailResponseBody service) {
		this.service = service;
	}

	public ServiceVersionDetailResponseBody getVersion() {
		return version;
	}

	public void setVersion(ServiceVersionDetailResponseBody version) {
		this.version = version;
	}

	public List<ServiceConfigDetailResponseBody> getConfigs() {
		return configs;
	}

	public void setConfigs(List<ServiceConfigDetailResponseBody> configs) {
		this.configs = configs;
	}

	@Override
	public String toString() {
		return "AddServiceConfigResponseBodyByQuickStart{" + "project=" + project + ", cluster=" + cluster + ", service=" + service + ", version=" + version + ", configs=" + configs + '}';
	}
}
