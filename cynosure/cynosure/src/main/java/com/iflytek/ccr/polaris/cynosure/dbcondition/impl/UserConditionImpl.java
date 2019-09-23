package com.iflytek.ccr.polaris.cynosure.dbcondition.impl;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IUserCondition;
import com.iflytek.ccr.polaris.cynosure.domain.User;
import com.iflytek.ccr.polaris.cynosure.enums.DBEnumInt;
import com.iflytek.ccr.polaris.cynosure.mapper.UserMapper;
import com.iflytek.ccr.polaris.cynosure.request.user.AddUserRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.user.EditUserRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.user.ModifyPasswordRequestBody;
import com.iflytek.ccr.polaris.cynosure.util.MD5Util;
import com.iflytek.ccr.polaris.cynosure.util.SnowflakeIdWorker;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.dao.DuplicateKeyException;
import org.springframework.stereotype.Service;

import java.util.Date;
import java.util.HashMap;
import java.util.List;

/**
 * 用户条件接口实现
 *
 * @author sctang2
 * @create 2017-12-20 9:25
 **/
@Service
public class UserConditionImpl extends BaseService implements IUserCondition {
    private static final EasyLogger logger = EasyLoggerFactory.getInstance(UserConditionImpl.class);

    @Autowired
    private UserMapper userMapper;

    @Override
    public int findTotalCount(HashMap<String, Object> map) {
        return this.userMapper.findTotalCount(map);
    }

    @Override
    public List<User> findList(HashMap<String, Object> map) {
        return this.userMapper.findList(map);
    }

    @Override
    public User findByAccount(String account) {
        return this.userMapper.findByAccount(account);
    }

    @Override
    public User findById(String id) {
        return this.userMapper.findById(id);
    }

    @Override
    public int updateById(String id, ModifyPasswordRequestBody body) {
        String password = body.getNewPassword();

        //更新密码
        String newPassword = MD5Util.getSaltMD5("cynosure", password);
        User updateUser = new User();
        updateUser.setId(id);
        updateUser.setPassword(newPassword);
        return this.userMapper.updateById(updateUser);
    }

    @Override
    public User updateById(String id, EditUserRequestBody body) {
        String userName = body.getUserName();
        String phone = body.getPhone();
        String email = body.getEmail();

        //更新
        Date now = new Date();
        User user = new User();
        user.setId(id);
        user.setUserName(userName);
        user.setPhone(phone);
        user.setEmail(email);
        user.setUpdateTime(now);
        this.userMapper.updateById(user);
        return user;
    }

    @Override
    public User add(AddUserRequestBody body) {
        String account = body.getAccount();
        String email = body.getEmail();
        String userName = body.getUserName();
        String phone = body.getPhone();

        //新增
        Date now = new Date();
        User user = new User();
        user.setAccount(account);
        user.setCreateTime(now);
        user.setUpdateTime(now);
        user.setEmail(email);
        user.setId(SnowflakeIdWorker.getId());
        String password = MD5Util.getSaltMD5("cynosure", "123456");
        user.setPassword(password);
        user.setPhone(phone);
        user.setUserName(userName);
        user.setRoleType((byte) DBEnumInt.ROLE_TYPE_USER.getIndex());
        try {
            this.userMapper.insert(user);
            return user;
        } catch (DuplicateKeyException ex) {
            logger.warn("user duplicate key " + ex.getMessage());
            return this.findByAccount(account);
        }
    }

    @Override
    public int deleteById(String id) {
        return this.userMapper.deleteById(id);
    }
}
