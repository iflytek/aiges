package com.iflytek.ccr.polaris.cynosure.request.serviceconfig;

import com.iflytek.ccr.polaris.cynosure.request.BaseRequestBody;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

import java.io.Serializable;

/**
 * 查询服务配置列表-请求
 *
 * @author sctang2
 * @create 2017-11-21 14:26
 **/
@ApiModel("查询服务配置列表请求参数")
public class QueryServiceConfigRequestBody extends BaseRequestBody implements Serializable {
	private static final long serialVersionUID = -5333839997546227035L;

	// 项目名称
	@ApiModelProperty("项目名称")
	private String project;

	// 集群名称
	@ApiModelProperty("集群名称")
	private String cluster;

	// 服务名称
	@ApiModelProperty("服务名称")
	private String service;

	// 服务版本名称
	@ApiModelProperty("服务版本名称")
	private String version;

	// 灰度组名称
	@ApiModelProperty("灰度组名称")
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
}
