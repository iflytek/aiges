package com.iflytek.ccr.polaris.cynosure.interceptor;

import java.io.IOException;
import java.io.PrintWriter;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import com.alibaba.fastjson.JSON;
import com.alibaba.fastjson.serializer.SerializerFeature;
import com.iflytek.ccr.polaris.cynosure.domain.User;
import com.iflytek.ccr.polaris.cynosure.response.Response;

/**
 * 认证处理器
 * 
 * @author jianchen15
 *
 */
public abstract class AuthenticationHandler {
	/**
	 * 认证
	 * 
	 * @param request
	 * @param response
	 * @return
	 */
	public abstract User authentication(HttpServletRequest request, HttpServletResponse response);

	/**
	 * 输出
	 *
	 * @param httpServletResponse
	 * @param code
	 * @param msg
	 * @throws Exception
	 */
	protected void output(HttpServletResponse httpServletResponse, int code, String msg) {
		Response<String> response = new Response<>(code, msg);
		httpServletResponse.setCharacterEncoding("UTF-8");
		httpServletResponse.setContentType("application/json");
		try (PrintWriter out = httpServletResponse.getWriter()) {
			out.print(JSON.toJSONString(response, SerializerFeature.WriteMapNullValue));
		} catch (IOException e) {
			throw new RuntimeException(e.getMessage());
		}

	}

	/**
	 * 获取认证方式
	 * 
	 * @param request
	 * @return
	 */
	public static AuthenticType getAuthenticType(HttpServletRequest request) {
		String appId = request.getHeader("app_id");
		if (appId == null || appId.isEmpty()) {
			return AuthenticType.NAME_PWD;
		}
		return AuthenticType.TOKEN;
	}

	public enum AuthenticType {
		NAME_PWD, TOKEN
	}

}
