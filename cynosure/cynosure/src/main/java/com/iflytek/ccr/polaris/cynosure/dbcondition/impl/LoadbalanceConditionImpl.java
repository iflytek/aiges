package com.iflytek.ccr.polaris.cynosure.dbcondition.impl;

import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.dbcondition.ILoadbalanceCondition;
import com.iflytek.ccr.polaris.cynosure.domain.LoadBalance;
import com.iflytek.ccr.polaris.cynosure.mapper.LoadBalanceMapper;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.HashMap;
import java.util.List;

/**
 * 负载均衡条件接口实现
 *
 * @author sctang2
 * @create 2017-12-11 10:48
 **/
@Service
public class LoadbalanceConditionImpl extends BaseService implements ILoadbalanceCondition {
    @Autowired
    private LoadBalanceMapper loadBalanceMapper;

    @Override
    public int findTotalCount(HashMap<String, Object> map) {
        return this.loadBalanceMapper.findTotalCount(map);
    }

    @Override
    public List<LoadBalance> findList(HashMap<String, Object> map) {
        return this.loadBalanceMapper.findList(map);
    }

    @Override
    public LoadBalance findByEnglishName(String abbr) {
        return this.loadBalanceMapper.findByEnglishName(abbr);
    }
}
