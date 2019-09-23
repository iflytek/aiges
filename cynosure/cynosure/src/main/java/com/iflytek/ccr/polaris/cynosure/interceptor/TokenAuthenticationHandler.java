package com.iflytek.ccr.polaris.cynosure.interceptor;

import java.security.NoSuchAlgorithmException;
import java.security.spec.InvalidKeySpecException;
import java.util.Date;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;

import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import com.alibaba.fastjson.JSON;
import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IUserCondition;
import com.iflytek.ccr.polaris.cynosure.domain.User;
import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.util.MD5Util;
import com.iflytek.ccr.polaris.cynosure.util.RSAUtils;

/**
 * 令牌方式认证
 * 
 * @author jianchen15
 *
 */
@Component("tokenAuthenticationHandler")
public class TokenAuthenticationHandler extends AuthenticationHandler {
	private final static String HEAD_TOKEN_NAME = "token";
	private final static String HEAD_APP_ID_NAME = "app_id";
	private final EasyLogger logger = EasyLoggerFactory.getInstance(TokenAuthenticationHandler.class);
	@Value("${private_key}")
	private String privateKey;// 私钥,当前在配置中配置,后期会从数据库中读取
	@Value("${valid_interval}")
	private long   validInterval;// token有效时间间隔(单位ms)
	@Autowired
	private IUserCondition userConditionImpl;

	@Override
	public User authentication(HttpServletRequest request, HttpServletResponse response) {
		String app_id = request.getHeader(HEAD_APP_ID_NAME);
		if (app_id == null || app_id.isEmpty()) {
			logger.error("app_id不正确");
			this.output(response, SystemErrCode.ERRCODE_USER_NOT_LOGIN, SystemErrCode.ERRMSG_APP_ID_ERROR);
			return null;
		}
		try {
			// 后期根据app_id获取私钥,当前私钥在配置中配置
			String tokenStr = RSAUtils.privateDecrypt(request.getHeader(HEAD_TOKEN_NAME), RSAUtils.getPrivateKey(privateKey));
			Token token = JSON.parseObject(tokenStr, Token.class);
			// 验证token是否在有效时间间隔内
			if ((new Date().getTime() - token.getTimestamp()) > validInterval) {
				logger.error("token:{}已失效,超过消息时间间隔:{},token生成日期为:{}", request.getHeader(HEAD_TOKEN_NAME), validInterval, token.getTimestamp());
				this.output(response, SystemErrCode.ERRCODE_USER_NOT_LOGIN, SystemErrCode.ERRMSG_TOKEN_ERROR);
				return null;
			}
			// 验证用户名和密码
			User user = this.userConditionImpl.findByAccount(token.getName());
			if (user == null) {
				logger.error("token:{}验证用户名:{}不正确", token.getName());
				this.output(response, SystemErrCode.ERRCODE_USER_NOT_EXISTS, SystemErrCode.ERRMSG_USER_NOT_EXISTS);
				return null;
			}
			if (StringUtils.isEmpty(token.getPwd()) || //
					!(MD5Util.getSaltMD5("cynosure", token.getPwd()).equals(user.getPassword()))) {
				logger.error("token:{}验证密码:{}不正确", token.getPwd());
				this.output(response, SystemErrCode.ERRCODE_USER_PASSWORD_INCORRECT, SystemErrCode.ERRMSG_USER_PASSWORD_INCORRECT);
				return null;
			}
			return user;
		} catch (NoSuchAlgorithmException | InvalidKeySpecException e) {
			logger.error("{}:解签失败,失败的原因是:{}", request.getHeader(HEAD_TOKEN_NAME), e.getMessage());
			this.output(response, SystemErrCode.ERRCODE_USER_NOT_LOGIN, SystemErrCode.ERRMSG_TOKEN_DECRPT_ERROR);
			return null;
		}
	}

	static class Token {
		private String name;
		private String pwd;
		private long timestamp;

		public String getName() {
			return name;
		}

		public void setName(String name) {
			this.name = name;
		}

		public String getPwd() {
			return pwd;
		}

		public void setPwd(String pwd) {
			this.pwd = pwd;
		}

		public long getTimestamp() {
			return timestamp;
		}

		public void setTimestamp(long timestamp) {
			this.timestamp = timestamp;
		}
	}
}
