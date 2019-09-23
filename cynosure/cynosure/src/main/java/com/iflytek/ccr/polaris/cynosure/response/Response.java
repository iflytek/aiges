package com.iflytek.ccr.polaris.cynosure.response;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

/**
 * 响应实体
 *
 * @author sctang2
 * @create 2017-11-09 15:59
 **/
@ApiModel("响应实体")
public class Response<T> {
	static final int SUCCESS = 0;
	static String SUCESSMSG = "成功";
	@ApiModelProperty("响应码,0:成功")
	private Integer code;
	@ApiModelProperty("响应消息")
	private String message;
	@ApiModelProperty("响应数据")
	private T data;

	public Integer getCode() {
		return code;
	}

	public void setCode(Integer code) {
		this.code = code;
	}

	public String getMessage() {
		return message;
	}

	public void setMessage(String message) {
		this.message = message;
	}

	public T getData() {
		return data;
	}

	public void setData(T data) {
		this.data = data;
	}

	@Override
	public String toString() {
		return "Response{" + "code=" + code + ", message='" + message + '\'' + ", data=" + data + '}';
	}

	public Response(Integer code, String message) {
		this.code = code;
		this.message = message;
	}

	public Response(T data) {
		this.code = SUCCESS;
		this.message = SUCESSMSG;
		this.data = data;
	}
}
