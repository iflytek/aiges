package com.iflytek.ccr.polaris.cynosure.response.user;

import java.io.Serializable;
import java.util.Date;

import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;

/**
 * 登录-响应
 *
 * @author sctang2
 * @create 2017-11-10 9:27
 **/
@ApiModel("登陆响应")
public class LoginResponseBody implements Serializable {
	private static final long serialVersionUID = 5279948331694041065L;

	// 用户id
	@ApiModelProperty("用户id")
	private String id;

	// 账号
	@ApiModelProperty("账户")
	private String account;

	// 用户名
	@ApiModelProperty("用户名")
	private String userName;

	// 手机号
	@ApiModelProperty("手机号")
	private String phone;

	// 邮箱
	@ApiModelProperty("邮箱")
	private String email;

	// 创建时间
	@ApiModelProperty("创建时间")
	private Date createTime;

	// 更新时间
	@ApiModelProperty("更新时间")
	private Date updateTime;

	// 角色类型
	@ApiModelProperty("角色类型")
	private Byte roleType;

	public String getId() {
		return id;
	}

	public void setId(String id) {
		this.id = id;
	}

	public String getAccount() {
		return account;
	}

	public void setAccount(String account) {
		this.account = account;
	}

	public String getUserName() {
		return userName;
	}

	public void setUserName(String userName) {
		this.userName = userName;
	}

	public String getPhone() {
		return phone;
	}

	public void setPhone(String phone) {
		this.phone = phone;
	}

	public String getEmail() {
		return email;
	}

	public void setEmail(String email) {
		this.email = email;
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

	public Byte getRoleType() {
		return roleType;
	}

	public void setRoleType(Byte roleType) {
		this.roleType = roleType;
	}

	@Override
	public String toString() {
		return "LoginResponseBody{" + "id='" + id + '\'' + ", account='" + account + '\'' + ", userName='" + userName + '\'' + ", phone='" + phone + '\'' + ", email='" + email + '\'' + ", createTime="
				+ createTime + ", updateTime=" + updateTime + ", roleType=" + roleType + '}';
	}
}
