package com.iflytek.ccr.polaris.cynosure.dbcondition;

import com.iflytek.ccr.polaris.cynosure.domain.User;
import com.iflytek.ccr.polaris.cynosure.request.user.AddUserRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.user.EditUserRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.user.ModifyPasswordRequestBody;

import java.util.HashMap;
import java.util.List;

/**
 * 用户条件接口
 *
 * @author sctang2
 * @create 2017-12-20 9:24
 **/
public interface IUserCondition {
    /**
     * 查询用户总数
     *
     * @param map
     * @return
     */
    int findTotalCount(HashMap<String, Object> map);

    /**
     * 查询用户列表
     *
     * @param map
     * @return
     */
    List<User> findList(HashMap<String, Object> map);

    /**
     * 通过账号查询用户信息
     *
     * @param account
     * @return
     */
    User findByAccount(String account);

    /**
     * 通过id查询用户信息
     *
     * @param id
     * @return
     */
    User findById(String id);

    /**
     * 通过id更新用户
     *
     * @param id
     * @param body
     * @return
     */
    int updateById(String id, ModifyPasswordRequestBody body);

    /**
     * 通过id更新用户
     *
     * @param id
     * @param body
     * @return
     */
    User updateById(String id, EditUserRequestBody body);

    /**
     * 新增用户
     *
     * @param body
     * @return
     */
    User add(AddUserRequestBody body);

    /**
     * 通过id删除用户
     *
     * @param id
     * @return
     */
    int deleteById(String id);
}
