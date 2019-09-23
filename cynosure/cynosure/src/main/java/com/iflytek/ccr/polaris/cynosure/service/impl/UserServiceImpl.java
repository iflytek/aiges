package com.iflytek.ccr.polaris.cynosure.service.impl;

import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IUserCondition;
import com.iflytek.ccr.polaris.cynosure.domain.User;
import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.request.BaseRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.user.AddUserRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.user.EditUserRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.user.LoginRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.user.ModifyPasswordRequestBody;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.user.LoginResponseBody;
import com.iflytek.ccr.polaris.cynosure.service.IUserService;
import com.iflytek.ccr.polaris.cynosure.util.MD5Util;
import com.iflytek.ccr.polaris.cynosure.util.PagingUtil;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.*;

/**
 * 用户业务逻辑接口实现
 *
 * @author sctang2
 * @create 2017-11-09 16:01
 **/
@Service
public class UserServiceImpl extends BaseService implements IUserService {
    @Autowired
    private IUserCondition userConditionImpl;

    @Override
    public Response<LoginResponseBody> login(LoginRequestBody body) {
        String account = body.getAccount();
        String password = body.getPassword();
        //通过用户名查询账号信息
        User user = this.userConditionImpl.findByAccount(account);
        if (null == user) {
            //不存在该用户
            return new Response<>(SystemErrCode.ERRCODE_USER_NOT_EXISTS, SystemErrCode.ERRMSG_USER_NOT_EXISTS);
        }

        //检查用户名密码
        String md5Password = MD5Util.getSaltMD5("cynosure", password);
        if (!md5Password.equals(user.getPassword())) {
            //账号密码错误
            return new Response<>(SystemErrCode.ERRCODE_USER_PASSWORD_INCORRECT, SystemErrCode.ERRMSG_USER_PASSWORD_INCORRECT);
        }

        //生成session
        this.httpServletRequest.getSession().setAttribute("user", user);

        //创建登录结果
        LoginResponseBody result = this.createLoginResult(user);
        return new Response<>(result);
    }

    /**
     * 创建登录结果
     *
     * @param user
     * @return
     */
    private LoginResponseBody createLoginResult(User user) {
        LoginResponseBody result = new LoginResponseBody();
        result.setId(user.getId());
        result.setAccount(user.getAccount());
        result.setCreateTime(user.getCreateTime());
        result.setEmail(user.getEmail());
        result.setPhone(user.getPhone());
        result.setUpdateTime(user.getUpdateTime());
        result.setUserName(user.getUserName());
        result.setRoleType(user.getRoleType());
        return result;
    }

    @Override
    public Response<LoginResponseBody> add(AddUserRequestBody body) {
        String account = body.getAccount();
        //通过用户名查询账号信息
        User user = this.userConditionImpl.findByAccount(account);
        if (null != user) {
            //已存在该用户
            return new Response<>(SystemErrCode.ERRCODE_USER_EXISTS, SystemErrCode.ERRMSG_USER_EXISTS);
        }

        //创建用户
        User newUser = this.userConditionImpl.add(body);

        //创建登录结果
        LoginResponseBody result = this.createLoginResult(newUser);
        return new Response<>(result);
    }

    @Override
    public Response<LoginResponseBody> edit(EditUserRequestBody body) {
        String id = body.getId();
        //通过用户id查询用户信息
        User user = this.userConditionImpl.findById(id);
        if (null == user) {
            //不存在该用户
            return new Response<>(SystemErrCode.ERRCODE_USER_NOT_EXISTS, SystemErrCode.ERRMSG_USER_NOT_EXISTS);
        }

        //修改用户信息
        User newUser = this.userConditionImpl.updateById(id, body);

        //创建登录结果
        String account = user.getAccount();
        String phone = user.getPhone();
        String email = user.getEmail();
        if (StringUtils.isNotBlank(account)) {
            newUser.setAccount(account);
        }
        if (StringUtils.isNotBlank(phone)) {
            newUser.setPhone(account);
        }
        if (StringUtils.isNotBlank(email)) {
            newUser.setEmail(email);
        }
        LoginResponseBody result = this.createLoginResult(newUser);
        return new Response<>(result);
    }

    @Override
    public Response<String> modifyPassword(ModifyPasswordRequestBody body) {
        //比较密码
        String newPassword = body.getNewPassword();
        String confirmPassword = body.getConfirmPassword();
        if (!newPassword.equals(confirmPassword)) {
            return new Response<>(SystemErrCode.ERRCODE_USER_PASSWORD_CONFIRM, SystemErrCode.ERRMSG_USER_PASSWORD_CONFIRM);
        }

        String id = this.getUserId();
        //通过id查询用户信息
        User userResult = this.userConditionImpl.findById(id);
        if (null == userResult) {
            //不存在该用户
            return new Response<>(SystemErrCode.ERRCODE_USER_NOT_EXISTS, SystemErrCode.ERRMSG_USER_NOT_EXISTS);
        }

        //检查密码
        String oldPassword = MD5Util.getSaltMD5("cynosure", body.getOldPassword());
        if (!oldPassword.equals(userResult.getPassword())) {
            //密码错误
            return new Response<>(SystemErrCode.ERRCODE_USER_PASSWORD_INCORRECT, SystemErrCode.ERRMSG_USER_PASSWORD_INCORRECT);
        }

        //更新密码
        this.userConditionImpl.updateById(id, body);

        //清空当前会话
        httpServletRequest.getSession().invalidate();
        return new Response<>(null);
    }

    @Override
    public Response<LoginResponseBody> find(String userId) {
        //通过id查询用户信息
        User user = this.userConditionImpl.findById(userId);
        if (null == user) {
            //不存在该用户
            return new Response<>(SystemErrCode.ERRCODE_USER_NOT_EXISTS, SystemErrCode.ERRMSG_USER_NOT_EXISTS);
        }

        //创建登录结果
        LoginResponseBody result = this.createLoginResult(user);
        return new Response<>(result);
    }

    @Override
    public Response<QueryPagingListResponseBody> findList(BaseRequestBody body) {
        //创建分页查询条件
        HashMap<String, Object> map = PagingUtil.createCondition(body);

        //查询总数
        int totalCount = this.userConditionImpl.findTotalCount(map);

        //创建分页结果
        QueryPagingListResponseBody result = PagingUtil.createResult(body, totalCount);
        if (0 == totalCount) {
            return new Response<>(result);
        }

        //查询列表
        List<LoginResponseBody> list = new ArrayList<>();
        Optional<List<User>> userList = Optional.ofNullable(this.userConditionImpl.findList(map));
        userList.ifPresent(x -> {
            x.forEach(y -> {
                LoginResponseBody userDetail = this.createLoginResult(y);
                list.add(userDetail);
            });
        });
        result.setList(list);
        return new Response<>(result);
    }

    @Override
    public Response<String> logout() {
        this.httpServletRequest.getSession().invalidate();
        return new Response<>(null);
    }
}
