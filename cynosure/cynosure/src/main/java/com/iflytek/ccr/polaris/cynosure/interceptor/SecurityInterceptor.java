package com.iflytek.ccr.polaris.cynosure.interceptor;

import com.alibaba.fastjson.JSON;
import com.alibaba.fastjson.serializer.SerializerFeature;
import com.iflytek.ccr.polaris.cynosure.annotation.Access;
import com.iflytek.ccr.polaris.cynosure.domain.User;
import com.iflytek.ccr.polaris.cynosure.enums.DBEnumInt;
import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.interceptor.AuthenticationHandler.AuthenticType;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import org.springframework.stereotype.Service;
import org.springframework.web.method.HandlerMethod;
import org.springframework.web.servlet.HandlerInterceptor;
import org.springframework.web.servlet.ModelAndView;
import javax.annotation.Resource;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.PrintWriter;
import java.lang.reflect.Method;
import java.util.HashSet;
import java.util.Set;

/**
 * 拦截器，用于登录和权限验证
 *
 * @author sctang2
 * @create 2017-11-10 14:44
 **/
@Service
public class SecurityInterceptor implements HandlerInterceptor {
	@Resource(name = "tokenAuthenticationHandler")
	private AuthenticationHandler tokeAnuthenticationHandler;
	@Resource(name = "namePwdAuthenticationHandler")
	private AuthenticationHandler namePwdAuthenticationHandler;

	@Override
	public boolean preHandle(HttpServletRequest httpServletRequest, HttpServletResponse httpServletResponse, Object handler) throws Exception {
		final String uri = httpServletRequest.getRequestURI();
		if (uri.contains("swagger-ui.html") || uri.contains("swagger-resources") || uri.contains("api-docs")) {
			return true;
		}
//		// 登录验证
//		User user = (User) httpServletRequest.getSession().getAttribute("user");
//		// 设置检查session的轮询时间
//		httpServletRequest.getSession().setMaxInactiveInterval(propUtil.MAXINTERVAL);
//		if (null == user) {
//			this.output(httpServletResponse, SystemErrCode.ERRCODE_USER_NOT_LOGIN, SystemErrCode.ERRMSG_USER_NOT_LOGIN);
//			return false;
//		}
		User user = null;
		if (AuthenticationHandler.getAuthenticType(httpServletRequest) == AuthenticType.NAME_PWD) {
			user = namePwdAuthenticationHandler.authentication(httpServletRequest, httpServletResponse);
		} else {
			user = tokeAnuthenticationHandler.authentication(httpServletRequest, httpServletResponse);
		}
		if (user == null) {
			return false;
		}
		// 权限验证
		HandlerMethod handlerMethod = (HandlerMethod) handler;
		Method method = handlerMethod.getMethod();
		Access access = method.getAnnotation(Access.class);
		if (access == null) {
			// 如果注解为null, 说明不需要拦截, 直接放过
			return true;
		}
		if (access.authorities().length > 0) {
			// 如果权限配置不为空, 则取出配置值
			String[] authorities = access.authorities();
			Set<String> authSet = new HashSet<String>();
			for (String authority : authorities) {
				// 将权限加入一个set集合中
				authSet.add(authority);
			}
			int roleType = user.getRoleType();
			String auth;
			if (DBEnumInt.ROLE_TYPE_ADMIN.getIndex() == roleType) {
				auth = "admin";
			} else {
				auth = "user";
			}
			if (authSet.contains(auth)) {
				// 校验通过返回true, 否则拦截请求
				return true;
			}
		}
		this.output(httpServletResponse, SystemErrCode.ERRCODE_NOT_AUTH, SystemErrCode.ERRMSG_NOT_AUTH);
		return false;
	}

	/**
	 * 输出
	 *
	 * @param httpServletResponse
	 * @param code
	 * @param msg
	 * @throws Exception
	 */
	private void output(HttpServletResponse httpServletResponse, int code, String msg) throws Exception {
		Response<String> response = new Response<>(code, msg);
		httpServletResponse.setCharacterEncoding("UTF-8");
		httpServletResponse.setContentType("application/json");
		PrintWriter out = httpServletResponse.getWriter();
		out.print(JSON.toJSONString(response, SerializerFeature.WriteMapNullValue));
		out.flush();
	}

	@Override
	public void postHandle(HttpServletRequest httpServletRequest, HttpServletResponse httpServletResponse, Object o, ModelAndView modelAndView) throws Exception {

	}

	@Override
	public void afterCompletion(HttpServletRequest httpServletRequest, HttpServletResponse httpServletResponse, Object o, Exception e) throws Exception {

	}
}
