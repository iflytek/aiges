package com.iflytek.ccr.polaris.cynosure.controller.v1;

import com.iflytek.ccr.polaris.cynosure.annotation.Access;
import com.iflytek.ccr.polaris.cynosure.consts.Constant;
import com.iflytek.ccr.polaris.cynosure.domain.User;
import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.request.BaseRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.user.*;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.user.LoginResponseBody;
import com.iflytek.ccr.polaris.cynosure.service.IUserService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.util.StringUtils;
import org.springframework.validation.BindingResult;
import org.springframework.validation.annotation.Validated;
import org.springframework.web.bind.annotation.*;

import javax.servlet.http.HttpServletRequest;

/**
 * 用户控制器
 *
 * @author sctang2
 * @create 2017-11-09 15:51
 **/
@RestController
@RequestMapping(Constant.API + "/{version}/user")
public class UserController {
	@Autowired
	private IUserService userServiceImpl;

	@Autowired
	private HttpServletRequest httpServletRequest;

	/**
	 * 登录
	 *
	 * @param body
	 * @param result
	 * @return
	 */
	@RequestMapping(value = "/login", method = RequestMethod.POST)
	public Response<LoginResponseBody> login(@Validated @RequestBody LoginRequestBody body, BindingResult result) {
		if (result.hasErrors()) {
			return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
		}
		return this.userServiceImpl.login(body);
	}

	/**
	 * 退出
	 *
	 * @return
	 */
	@RequestMapping(value = "/logout", method = RequestMethod.POST)
	public Response<String> logout() {
		return this.userServiceImpl.logout();
	}

	/**
	 * 修改密码
	 *
	 * @param body
	 * @param result
	 * @return
	 */
	@RequestMapping(value = "/modifyPassword", method = RequestMethod.POST)
	public Response<String> modifyPassword(@Validated @RequestBody ModifyPasswordRequestBody body, BindingResult result) {
		if (result.hasErrors()) {
			return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
		}
		return this.userServiceImpl.modifyPassword(body);
	}

	/**
	 * 新增用户
	 *
	 * @param body
	 * @param result
	 * @return
	 */
	@RequestMapping(value = "/add", method = RequestMethod.POST)
	@Access(authorities = "admin")
	public Response<LoginResponseBody> add(@Validated @RequestBody AddUserRequestBody body, BindingResult result) {
		if (result.hasErrors()) {
			return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
		}
		return this.userServiceImpl.add(body);
	}

	/**
	 * 编辑用户
	 *
	 * @param body
	 * @param result
	 * @return
	 */
	@RequestMapping(value = "/edit", method = RequestMethod.POST)
	@Access(authorities = "admin")
	public Response<LoginResponseBody> edit(@Validated @RequestBody EditUserRequestBody body, BindingResult result) {
		if (result.hasErrors()) {
			return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
		}
		return this.userServiceImpl.edit(body);
	}

	/**
	 * 查询用户列表
	 *
	 * @param body
	 * @return
	 */
	@RequestMapping(value = "/list", method = RequestMethod.GET)
	@Access(authorities = "admin")
	public Response<QueryPagingListResponseBody> list(BaseRequestBody body) {
		return this.userServiceImpl.findList(body);
	}

	/**
	 * 查询用户详情
	 *
	 * @param id
	 * @return
	 */
	@RequestMapping(value = "/detail", method = RequestMethod.GET)
	@Access(authorities = "admin")
	public Response<LoginResponseBody> find(@RequestParam("id") String id) {
		if (StringUtils.isEmpty(id)) {
			return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, SystemErrCode.ERRMSG_ID_NOT_NULL);
		}
		return this.userServiceImpl.find(id);
	}

	/**
	 * 查询个人信息
	 *
	 * @return
	 */
	@RequestMapping(value = "/info", method = RequestMethod.GET)
	public Response<LoginResponseBody> personInfo() {
		User user = (User) httpServletRequest.getSession().getAttribute("user");
		String id = user.getId();
		return this.userServiceImpl.find(id);
	}

	/**
	 * 编辑个人信息
	 *
	 * @param body
	 * @param result
	 * @return
	 */
	@RequestMapping(value = "/editInfo", method = RequestMethod.POST)
	public Response<LoginResponseBody> edit(@Validated @RequestBody EditUserInfoRequestBody body, BindingResult result) {
		if (result.hasErrors()) {
			return new Response<>(SystemErrCode.ERRCODE_INVALID_PARAMETER, result.getFieldError().getDefaultMessage());
		}
		User user = (User) httpServletRequest.getSession().getAttribute("user");
		String id = user.getId();
		EditUserRequestBody newBody = new EditUserRequestBody();
		newBody.setId(id);
		newBody.setEmail(body.getEmail());
		newBody.setPhone(body.getPhone());
		newBody.setUserName(body.getUserName());
		return this.userServiceImpl.edit(newBody);
	}
}
