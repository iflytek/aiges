package com.iflytek.ccr.polaris.cynosure.request.serviceconfig;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

import org.hibernate.validator.constraints.Length;
import org.hibernate.validator.constraints.NotBlank;

import java.io.Serializable;

/**
 * 编辑服务配置-请求
 *
 * @author sctang2
 * @create 2017-11-21 13:59
 **/
@ApiModel("编辑服务配置请求参数")
public class EditServiceConfigRequestBody implements Serializable {
	private static final long serialVersionUID = 7690217038065030241L;

	// 配置id
	@NotBlank(message = SystemErrCode.ERRMSG_SERVICE_CONFIG_ID_NOT_NULL)
	@ApiModelProperty(required = true, notes = "配置id")
	private String id;

	// 配置描述
	@Length(max = 500, message = SystemErrCode.ERRMSG_SERVICE_CONFIG_DESC_MAX_LENGTH)
	@ApiModelProperty(notes = "配置描述")
	private String desc;

	// 配置内容
	@NotBlank(message = SystemErrCode.ERRMSG_SERVICE_CONFIG_CONTENT_NOT_NULL)
//    @Length(max = 20240, message = SystemErrCode.ERRMSG_SERVICE_CONFIG_CONTENT_MAX_LENGTH)
	@ApiModelProperty(required = true, notes = "配置内容")
	private String content;

	public String getId() {
		return id;
	}

	public void setId(String id) {
		this.id = id;
	}

	public String getDesc() {
		return desc;
	}

	public void setDesc(String desc) {
		this.desc = desc;
	}

	public String getContent() {
		return content;
	}

	public void setContent(String content) {
		this.content = content;
	}

	@Override
	public String toString() {
		return "EditServiceConfigRequestBody{" + "id='" + id + '\'' + ", desc='" + desc + '\'' + ", content='" + content + '\'' + '}';
	}
}
