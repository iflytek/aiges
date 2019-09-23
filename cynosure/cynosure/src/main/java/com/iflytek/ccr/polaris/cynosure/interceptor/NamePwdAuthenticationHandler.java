package com.iflytek.ccr.polaris.cynosure.interceptor;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import com.iflytek.ccr.polaris.cynosure.domain.User;
import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.util.PropUtil;

/**
 * 用户名和密码方式认证
 * 
 * @author jianchen15
 *
 */
@Component("namePwdAuthenticationHandler")
public class NamePwdAuthenticationHandler extends AuthenticationHandler {
	@Autowired
	private PropUtil propUtil;

	@Override
	public User authentication(HttpServletRequest request, HttpServletResponse response) {
		User user = (User) request.getSession().getAttribute("user");
		// 设置检查session的轮询时间
		request.getSession().setMaxInactiveInterval(propUtil.MAXINTERVAL);
		if (user == null) {
			this.output(response, SystemErrCode.ERRCODE_USER_NOT_LOGIN, SystemErrCode.ERRMSG_USER_NOT_LOGIN);
		}
		return user;
	}

}
