package com.iflytek.ccr.polaris.cynosure.dbcondition.impl;

import com.alibaba.fastjson.JSON;
import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.customdomain.SearchCondition;
import com.iflytek.ccr.polaris.cynosure.dbcondition.ILastestSearchCondition;
import com.iflytek.ccr.polaris.cynosure.domain.LastestSearch;
import com.iflytek.ccr.polaris.cynosure.mapper.LastestSearchMapper;
import com.iflytek.ccr.polaris.cynosure.util.SnowflakeIdWorker;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.dao.DuplicateKeyException;
import org.springframework.stereotype.Service;

import java.util.Date;

/**
 * 最近搜索条件接口实现
 *
 * @author sctang2
 * @create 2018-01-26 10:21
 **/
@Service
public class LastestSearchConditionImpl extends BaseService implements ILastestSearchCondition {
    private static final EasyLogger logger = EasyLoggerFactory.getInstance(LastestSearchConditionImpl.class);

    @Autowired
    private LastestSearchMapper lastestSearchMapper;

    @Override
    public int insert(String url, String condition) {
        String userId = this.getUserId();

        //新增
        Date now = new Date();
        LastestSearch lastestSearch = new LastestSearch();
        lastestSearch.setId(SnowflakeIdWorker.getId());
        lastestSearch.setPreCondition(condition);
        lastestSearch.setCreateTime(now);
        lastestSearch.setUserId(userId);
        lastestSearch.setUrl(url);
        try {
            return this.lastestSearchMapper.insert(lastestSearch);
        } catch (DuplicateKeyException ex) {
            logger.warn("lastest search duplicate key " + ex.getMessage());
            return 1;
        }
    }

    @Override
    public int update(String url, String condition) {
        String userId = this.getUserId();

        //更新
        Date now = new Date();
        LastestSearch lastestSearch = new LastestSearch();
        lastestSearch.setUserId(userId);
        lastestSearch.setUrl(url);
        lastestSearch.setUpdateTime(now);
        lastestSearch.setPreCondition(condition);
        return this.lastestSearchMapper.update(lastestSearch);
    }

    @Override
    public SearchCondition find(String url) {
        String userId = this.getUserId();

        //更新
        LastestSearch lastestSearchCondition = new LastestSearch();
        lastestSearchCondition.setUserId(userId);
        lastestSearchCondition.setUrl(url);
        LastestSearch lastestSearch = this.lastestSearchMapper.find(lastestSearchCondition);
        if (null == lastestSearch) {
            return null;
        }
        String condition = lastestSearch.getPreCondition();
        return JSON.parseObject(condition, SearchCondition.class);
    }

    @Override
    public int syncSearchCondition(String url, String condition) {
        //查询搜索条件
        int success;
        SearchCondition search = this.find(url);
        if (null == search) {
            success = this.insert(url, condition);
        } else {
            success = this.update(url, condition);
        }
        return success;
    }
}
