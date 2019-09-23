package com.iflytek.ccr.polaris.cynosure.request.user;

import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

import org.hibernate.validator.constraints.NotBlank;

import java.io.Serializable;

/**
 * 登录-请求
 *
 * @author sctang2
 * @create 2017-11-10 9:36
 **/
@ApiModel("登陆请求参数")
public class LoginRequestBody implements Serializable {
	private static final long serialVersionUID = 636644145376989802L;
	// 账号
	@NotBlank(message = SystemErrCode.ERRMSG_USER_ACCOUNT_NOT_NULL)
	@ApiModelProperty(required = true, notes = "账号")
	private String account;

	// 密码
	@NotBlank(message = SystemErrCode.ERRMSG_USER_PASSWORD_NOT_NULL)
	@ApiModelProperty(required = true, notes = "密码")
	private String password;

	public String getAccount() {
		return account;
	}

	public void setAccount(String account) {
		this.account = account;
	}

	public String getPassword() {
		return password;
	}

	public void setPassword(String password) {
		this.password = password;
	}

	@Override
	public String toString() {
		return "LoginRequestBody{" + "account='" + account + '\'' + ", password='" + password + '\'' + '}';
	}
}