package com.iflytek.ccr.polaris.cynosure.mapper;

import com.iflytek.ccr.polaris.cynosure.domain.User;

import java.util.HashMap;
import java.util.List;

/**
 * 用户持久化接口
 *
 * @author sctang2
 * @create 2017-11-10 9:18
 **/
public interface UserMapper {
    /**
     * 新增用户
     *
     * @param user
     * @return
     */
    int insert(User user);

    /**
     * 通过id删除用户
     *
     * @param id
     * @return
     */
    int deleteById(String id);

    /**
     * 通过id更新用户
     *
     * @param user
     * @return
     */
    int updateById(User user);

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
}
