package com.iflytek.ccr.polaris.cynosure.request.serviceconfig;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

import org.hibernate.validator.constraints.NotBlank;

import javax.validation.constraints.NotNull;
import javax.validation.constraints.Size;
import java.io.Serializable;
import java.util.List;

/**
 * 服务配置推送-请求
 *
 * @author sctang2
 * @create 2017-11-27 17:19
 **/
@ApiModel("服务配置推送请求参数")
public class PushServiceConfigRequestBody implements Serializable {
	private static final long serialVersionUID = 6049857945750309290L;

	// 配置id
	@NotBlank(message = SystemErrCode.ERRMSG_SERVICE_CONFIG_ID_NOT_NULL)
	@ApiModelProperty(required = true, notes = "配置id")
	private String id;

	// 区域id列表,至少应包含一个区域
	@NotNull(message = SystemErrCode.ERRMSG_REGION_IDS_NOT_NULL)
	@Size(min = 1, message = SystemErrCode.ERRMSG_REGION_IDS_IS_NOT_EMPTY)
	@ApiModelProperty(required = true, notes = "推送区域")
	private List<String> regionIds;

	public String getId() {
		return id;
	}

	public void setId(String id) {
		this.id = id;
	}

	public List<String> getRegionIds() {
		return regionIds;
	}

	public void setRegionIds(List<String> regionIds) {
		this.regionIds = regionIds;
	}

	@Override
	public String toString() {
		return "PushServiceConfigRequestBody{" + "id='" + id + '\'' + ", regionIds=" + regionIds + '}';
	}
}
