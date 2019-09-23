package com.iflytek.ccr.polaris.cynosure.response.serviceconfig;

import java.io.Serializable;
import java.util.Date;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

/**
 * 服务配置明细-响应
 *
 * @author sctang2
 * @create 2017-11-21 11:54
 **/
@ApiModel("服务配置明细")
public class ServiceConfigDetailResponseBody implements Serializable {
	private static final long serialVersionUID = 3465537468615599850L;

	// 配置id
	@ApiModelProperty(notes = "配置id")
	private String id;

	// 版本id
	@ApiModelProperty(notes = "版本id")
	private String versionId;

	// 灰度组配置id
	@ApiModelProperty(notes = "灰度组配置id")
	private String grayId;

	// 配置名称
	@ApiModelProperty(notes = "配置名称")
	private String name;

	// 配置路径
	@ApiModelProperty(notes = "配置路径")
	private String path;

	// 配置内容
	@ApiModelProperty(notes = "配置内容")
	private String content;

	// 配置描述
	@ApiModelProperty(notes = "配置描述")
	private String desc;

	// 创建时间
	@ApiModelProperty(notes = "创建时间")
	private Date createTime;

	// 更新时间
	@ApiModelProperty(notes = "更新时间")
	private Date updateTime;

	public String getId() {
		return id;
	}

	public void setId(String id) {
		this.id = id;
	}

	public String getVersionId() {
		return versionId;
	}

	public void setVersionId(String versionId) {
		this.versionId = versionId;
	}

	public String getName() {
		return name;
	}

	public void setName(String name) {
		this.name = name;
	}

	public String getPath() {
		return path;
	}

	public void setPath(String path) {
		this.path = path;
	}

	public String getContent() {
		return content;
	}

	public void setContent(String content) {
		this.content = content;
	}

	public String getDesc() {
		return desc;
	}

	public void setDesc(String desc) {
		this.desc = desc;
	}

	public Date getCreateTime() {
		return createTime;
	}

	public void setCreateTime(Date createTime) {
		this.createTime = createTime;
	}

	public Date getUpdateTime() {
		return updateTime;
	}

	public void setUpdateTime(Date updateTime) {
		this.updateTime = updateTime;
	}

	public String getGrayId() {
		return grayId;
	}

	public void setGrayId(String grayId) {
		this.grayId = grayId;
	}

	@Override
	public String toString() {
		return "ServiceConfigDetailResponseBody{" + "id='" + id + '\'' + ", versionId='" + versionId + '\'' + ", grayId='" + grayId + '\'' + ", name='" + name + '\'' + ", path='" + path + '\''
				+ ", content='" + content + '\'' + ", desc='" + desc + '\'' + ", createTime=" + createTime + ", updateTime=" + updateTime + '}';
	}
}
