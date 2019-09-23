package com.iflytek.ccr.polaris.cynosure.service;

import com.iflytek.ccr.polaris.cynosure.request.BaseRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.user.AddUserRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.user.EditUserRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.user.LoginRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.user.ModifyPasswordRequestBody;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.user.LoginResponseBody;

/**
 * 用户业务逻辑接口
 *
 * @author sctang2
 * @create 2017-11-09 16:01
 **/
public interface IUserService {
    /**
     * 登录
     *
     * @param body
     * @return
     */
    Response<LoginResponseBody> login(LoginRequestBody body);

    /**
     * 新增用户
     *
     * @param body
     * @return
     */
    Response<LoginResponseBody> add(AddUserRequestBody body);

    /**
     * 编辑用户
     *
     * @param body
     * @return
     */
    Response<LoginResponseBody> edit(EditUserRequestBody body);

    /**
     * 修改密码
     *
     * @param body
     * @return
     */
    Response<String> modifyPassword(ModifyPasswordRequestBody body);

    /**
     * 查询用户详情
     *
     * @param id
     * @return
     */
    Response<LoginResponseBody> find(String id);

    /**
     * 查询用户列表
     *
     * @param body
     * @return
     */
    Response<QueryPagingListResponseBody> findList(BaseRequestBody body);

    /**
     * 退出
     *
     * @return
     */
    Response<String> logout();
}
