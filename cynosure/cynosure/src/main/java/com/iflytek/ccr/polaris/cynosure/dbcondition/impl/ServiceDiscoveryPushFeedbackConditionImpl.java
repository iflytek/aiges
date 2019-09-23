package com.iflytek.ccr.polaris.cynosure.dbcondition.impl;

import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IServiceDiscoveryPushFeedbackCondition;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceDiscoveryPushFeedback;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceDiscoveryPushFeedbackDetail;
import com.iflytek.ccr.polaris.cynosure.mapper.ServiceDiscoveryPushFeedbackMapper;
import com.iflytek.ccr.polaris.cynosure.request.servicediscovery.ServiceDiscoveryFeedBackRequestBody;
import com.iflytek.ccr.polaris.cynosure.util.SnowflakeIdWorker;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.Date;
import java.util.HashMap;
import java.util.List;

/**
 * 服务发现推送反馈条件接口实现
 *
 * @author sctang2
 * @create 2017-12-12 17:26
 **/
@Service
public class ServiceDiscoveryPushFeedbackConditionImpl extends BaseService implements IServiceDiscoveryPushFeedbackCondition {
    @Autowired
    private ServiceDiscoveryPushFeedbackMapper serviceDiscoveryPushFeedbackMapper;

    @Override
    public int findTotalCount(HashMap<String, Object> map) {
        return this.serviceDiscoveryPushFeedbackMapper.findTotalCount(map);
    }

    @Override
    public int delete(String pushId) {
        return this.serviceDiscoveryPushFeedbackMapper.deleteByPushId(pushId);
    }

    @Override
    public int deleteByPushIds(List<String> pushIds) {
        return this.serviceDiscoveryPushFeedbackMapper.deleteByPushIds(pushIds);
    }

    @Override
    public List<ServiceDiscoveryPushFeedbackDetail> findList(HashMap<String, Object> map) {
        return this.serviceDiscoveryPushFeedbackMapper.findList(map);
    }

    @Override
    public ServiceDiscoveryPushFeedback add(ServiceDiscoveryFeedBackRequestBody body) {
        String pushId = body.getPushId();
        String consumerVersion = body.getConsumerVersion();
        String providerVersion = body.getProviderVersion();
        String project = body.getProject();
        String group = body.getGroup();
        String consumer = body.getConsumer();
        String provider = body.getProvider();
        String addr = body.getAddr();
        int updateStatus = body.getUpdateStatus();
        int loadStatus = body.getLoadStatus();
        Date updateTime = body.getUpdateTime();
        Date loadTime = body.getLoadTime();
        String type = body.getType();
        String apiVersion = body.getApiVersion();

        //新增
        Date now = new Date();
        ServiceDiscoveryPushFeedback serviceDiscoveryPushFeedback = new ServiceDiscoveryPushFeedback();
        serviceDiscoveryPushFeedback.setId(SnowflakeIdWorker.getId());
        serviceDiscoveryPushFeedback.setPushId(pushId);
        serviceDiscoveryPushFeedback.setProject(project);
        serviceDiscoveryPushFeedback.setServiceGroup(group);
        serviceDiscoveryPushFeedback.setConsumerService(consumer);
        serviceDiscoveryPushFeedback.setConsumerVersion(consumerVersion);
        serviceDiscoveryPushFeedback.setProviderService(provider);
        serviceDiscoveryPushFeedback.setProviderVersion(providerVersion);
        serviceDiscoveryPushFeedback.setAddr(addr);
        serviceDiscoveryPushFeedback.setUpdateStatus((byte) updateStatus);
        serviceDiscoveryPushFeedback.setLoadStatus((byte) loadStatus);
        serviceDiscoveryPushFeedback.setUpdateTime(updateTime);
        serviceDiscoveryPushFeedback.setLoadTime(loadTime);
        serviceDiscoveryPushFeedback.setCreateTime(now);
        serviceDiscoveryPushFeedback.setApiVersion(apiVersion);
        serviceDiscoveryPushFeedback.setType(type);
        this.serviceDiscoveryPushFeedbackMapper.insert(serviceDiscoveryPushFeedback);
        return serviceDiscoveryPushFeedback;
    }
}
