package com.iflytek.ccr.polaris.cynosure.dbcondition.impl;

import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.companionservice.domain.PushResult;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IServiceDiscoveryPushHistoryCondition;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceDiscoveryPushHistory;
import com.iflytek.ccr.polaris.cynosure.mapper.ServiceDiscoveryPushHistoryMapper;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.Date;
import java.util.HashMap;
import java.util.List;

/**
 * 服务发现推送历史条件接口实现
 *
 * @author sctang2
 * @create 2017-12-12 17:27
 **/
@Service
public class ServiceDiscoveryPushHistoryConditionImpl extends BaseService implements IServiceDiscoveryPushHistoryCondition {
    @Autowired
    private ServiceDiscoveryPushHistoryMapper serviceDiscoveryPushHistoryMapper;

    @Override
    public int findTotalCount(HashMap<String, Object> map) {
        return this.serviceDiscoveryPushHistoryMapper.findTotalCount(map);
    }

    @Override
    public List<ServiceDiscoveryPushHistory> findList(HashMap<String, Object> map) {
        return this.serviceDiscoveryPushHistoryMapper.findList(map);
    }

    @Override
    public ServiceDiscoveryPushHistory add(String project, String cluster, String service, String apiVersion, PushResult cacheCenterPushResult) {
        String pushId = cacheCenterPushResult.getPushId();
        String pushResult = cacheCenterPushResult.getResult();
        String userId = this.getUserId();

        //新增
        Date now = new Date();
        ServiceDiscoveryPushHistory serviceDiscoveryPushHistory = new ServiceDiscoveryPushHistory();
        serviceDiscoveryPushHistory.setId(pushId);
        serviceDiscoveryPushHistory.setUserId(userId);
        serviceDiscoveryPushHistory.setClusterText(pushResult);
        serviceDiscoveryPushHistory.setProject(project);
        serviceDiscoveryPushHistory.setServiceGroup(cluster);
        serviceDiscoveryPushHistory.setService(service);
        serviceDiscoveryPushHistory.setVersion(apiVersion);
        serviceDiscoveryPushHistory.setPushTime(now);
        this.serviceDiscoveryPushHistoryMapper.insert(serviceDiscoveryPushHistory);
        return serviceDiscoveryPushHistory;
    }

    @Override
    public ServiceDiscoveryPushHistory findById(String id) {
        return this.serviceDiscoveryPushHistoryMapper.findById(id);
    }

    @Override
    public int deleteById(String id) {
        return this.serviceDiscoveryPushHistoryMapper.deleteById(id);
    }

    @Override
    public int deleteByIds(List<String> ids) {
        return this.serviceDiscoveryPushHistoryMapper.deleteByIds(ids);
    }
}
