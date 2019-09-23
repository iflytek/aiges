package com.iflytek.ccr.polaris.cynosure.request.quickstart;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import org.hibernate.validator.constraints.Length;
import org.hibernate.validator.constraints.NotBlank;
import java.io.Serializable;

/**
 * 新增服务配置-请求
 *
 * @author sctang2
 * @create 2018-01-30 9:12
 **/
public class AddServiceConfigRequestBodyByQuickStart implements Serializable {
	private static final long serialVersionUID = 2050698860258042066L;

	// 项目名称
	@NotBlank(message = SystemErrCode.ERRMSG_PROJECT_NAME_NOT_NULL)
	@Length(min = 1, max = 100, message = SystemErrCode.ERRMSG_PROJECT_NAME_MAX_LENGTH)
	private String project;

	// 集群名称
	@NotBlank(message = SystemErrCode.ERRMSG_CLUSTER_NAME_NOT_NULL)
	@Length(min = 1, max = 100, message = SystemErrCode.ERRMSG_CLUSTER_NAME_MAX_LENGTH)
	private String cluster;

	// 服务名称
	@NotBlank(message = SystemErrCode.ERRMSG_SERVICE_NAME_NOT_NULL)
	@Length(min = 1, max = 100, message = SystemErrCode.ERRMSG_SERVICE_NAME_MAX_LENGTH)
	private String service;

	// 服务版本
	@NotBlank(message = SystemErrCode.ERRMSG_SERVICE_VERSION_NOT_NULL)
	@Length(min = 1, max = 20, message = SystemErrCode.ERRMSG_SERVICE_VERSION_MAX_LENGTH)
	private String version;

	// 配置描述
	@Length(max = 500, message = SystemErrCode.ERRMSG_SERVICE_CONFIG_DESC_MAX_LENGTH)
	private String desc;

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

	public String getDesc() {
		return desc;
	}

	public void setDesc(String desc) {
		this.desc = desc;
	}

	@Override
	public String toString() {
		return "AddServiceConfigRequestBodyByQuickStart{" + "project='" + project + '\'' + ", cluster='" + cluster + '\'' + ", service='" + service + '\'' + ", version='" + version + '\'' + ", desc='"
				+ desc + '\'' + '}';
	}
}
