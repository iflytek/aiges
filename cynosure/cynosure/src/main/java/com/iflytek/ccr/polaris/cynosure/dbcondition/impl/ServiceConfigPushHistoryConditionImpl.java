package com.iflytek.ccr.polaris.cynosure.dbcondition.impl;

import com.alibaba.fastjson.JSONArray;
import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.companionservice.domain.PushResult;
import com.iflytek.ccr.polaris.cynosure.customdomain.ServiceConfigText;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IServiceConfigPushHistoryCondition;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceConfig;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceConfigPushHistory;
import com.iflytek.ccr.polaris.cynosure.mapper.ServiceConfigPushHistoryMapper;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.Date;
import java.util.HashMap;
import java.util.List;

/**
 * 服务推送历史条件接口实现
 *
 * @author sctang2
 * @create 2017-12-10 22:24
 **/
@Service
public class ServiceConfigPushHistoryConditionImpl extends BaseService implements IServiceConfigPushHistoryCondition {
    @Autowired
    private ServiceConfigPushHistoryMapper servicePushHistoryMapper;

    @Override
    public int findTotalCount(HashMap<String, Object> map) {
        return this.servicePushHistoryMapper.findTotalCount(map);
    }

    @Override
    public List<ServiceConfigPushHistory> findList(HashMap<String, Object> map) {
        return this.servicePushHistoryMapper.findList(map);
    }

    @Override
    public ServiceConfigPushHistory add(ServiceConfig serviceConfig, PushResult cacheCenterPushResult) {

        List<ServiceConfigText> serviceConfigTexts = new ArrayList<>();
        ServiceConfigText serviceConfigText = new ServiceConfigText();
        serviceConfigText.setName(serviceConfig.getName());
        serviceConfigTexts.add(serviceConfigText);

        //构造推送历史记录实体
        ServiceConfigPushHistory servicePushHistory = new ServiceConfigPushHistory(
                cacheCenterPushResult.getPushId(), this.getUserId(), serviceConfig.getGrayId(),
                serviceConfig.getProject().getName(), serviceConfig.getCluster().getName(), serviceConfig.getService().getName(),
                serviceConfig.getServiceVersion().getVersion(), cacheCenterPushResult.getResult(), JSONArray.toJSONString(serviceConfigTexts),
                new Date());

        //新增
        this.servicePushHistoryMapper.insert(servicePushHistory);
        return servicePushHistory;
    }

    @Override
    public ServiceConfigPushHistory add(List<ServiceConfig> serviceConfigList, PushResult cacheCenterPushResult) {
        ServiceConfig serviceConfig = serviceConfigList.get(0);
        String pushId = cacheCenterPushResult.getPushId();
        String pushResult = cacheCenterPushResult.getResult();
        String userId = this.getUserId();
        String project = serviceConfig.getProject().getName();
        String group = serviceConfig.getCluster().getName();
        String service = serviceConfig.getService().getName();
        String version = serviceConfig.getServiceVersion().getVersion();
        String grayId = serviceConfig.getGrayId();


        //新增
        Date now = new Date();
        ServiceConfigPushHistory servicePushHistory = new ServiceConfigPushHistory();
        servicePushHistory.setId(pushId);
        servicePushHistory.setUserId(userId);
        servicePushHistory.setProject(project);
        servicePushHistory.setServiceGroup(group);
        servicePushHistory.setService(service);
        servicePushHistory.setVersion(version);
        servicePushHistory.setGrayId(grayId);
        servicePushHistory.setPushTime(now);
        List<ServiceConfigText> serviceConfigTexts = new ArrayList<>();
        serviceConfigList.forEach(x -> {
            ServiceConfigText serviceConfigText = new ServiceConfigText();
            serviceConfigText.setName(x.getName());
            serviceConfigTexts.add(serviceConfigText);
        });
        servicePushHistory.setServiceConfigText(JSONArray.toJSONString(serviceConfigTexts));
        servicePushHistory.setClusterText(pushResult);
        this.servicePushHistoryMapper.insert(servicePushHistory);
        return servicePushHistory;
    }

    @Override
    public ServiceConfigPushHistory findById(String id) {
        return this.servicePushHistoryMapper.findById(id);
    }

    @Override
    public int deleteById(String id) {
        return this.servicePushHistoryMapper.deleteById(id);
    }

    @Override
    public int deleteByIds(List<String> ids) {
        return this.servicePushHistoryMapper.deleteByIds(ids);
    }
}
