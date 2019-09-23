package com.iflytek.ccr.polaris.cynosure.dbcondition.impl;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IServiceConfigPushFeedbackCondition;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceConfigPushFeedback;
import com.iflytek.ccr.polaris.cynosure.mapper.ServiceConfigPushFeedbackMapper;
import com.iflytek.ccr.polaris.cynosure.request.serviceconfig.ServiceConfigFeedBackRequestBody;
import com.iflytek.ccr.polaris.cynosure.util.SnowflakeIdWorker;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.Date;
import java.util.HashMap;
import java.util.List;

/**
 * 服务配置推送反馈条件接口实现
 *
 * @author sctang2
 * @create 2017-12-11 9:05
 **/
@Service
public class ServiceConfigPushFeedbackConditionImpl extends BaseService implements IServiceConfigPushFeedbackCondition {
    private final EasyLogger logger = EasyLoggerFactory.getInstance(ServiceConfigPushFeedbackConditionImpl.class);
    @Autowired
    private ServiceConfigPushFeedbackMapper serviceConfigPushFeedbackMapper;

    @Override
    public int findTotalCount(HashMap<String, Object> map) {
        return this.serviceConfigPushFeedbackMapper.findTotalCount(map);
    }

    @Override
    public int delete(String pushId) {
        return this.serviceConfigPushFeedbackMapper.deleteByPushId(pushId);
    }

    @Override
    public int deleteByPushIds(List<String> pushIds) {
        return this.serviceConfigPushFeedbackMapper.deleteByPushIds(pushIds);
    }

    @Override
    public List<ServiceConfigPushFeedback> findList(HashMap<String, Object> map) {
        return this.serviceConfigPushFeedbackMapper.findList(map);
    }

    @Override
    public ServiceConfigPushFeedback add(ServiceConfigFeedBackRequestBody body) {
        String pushId = body.getPushId();
        String version = body.getVersion();
        String project = body.getProject();
        String group = body.getGroup();
        String service = body.getService();
        String config = body.getConfig();
        String addr = body.getAddr();
        int updateStatus = body.getUpdateStatus();
        int loadStatus = body.getLoadStatus();
        Date updateTime = body.getUpdateTime();
        Date loadTime = body.getLoadTime();
        String grayGroupId = body.getGrayGroupId();
        if (grayGroupId == null){
            grayGroupId = "0";
        }

        //新增
        Date now = new Date();
        ServiceConfigPushFeedback serviceConfigPushFeedback = new ServiceConfigPushFeedback();
        serviceConfigPushFeedback.setId(SnowflakeIdWorker.getId());
        serviceConfigPushFeedback.setPushId(pushId);
        serviceConfigPushFeedback.setProject(project);
        serviceConfigPushFeedback.setServiceGroup(group);
        serviceConfigPushFeedback.setService(service);
        serviceConfigPushFeedback.setVersion(version);
        serviceConfigPushFeedback.setConfig(config);
        serviceConfigPushFeedback.setAddr(addr);
        serviceConfigPushFeedback.setUpdateStatus((byte) updateStatus);
        serviceConfigPushFeedback.setLoadStatus((byte) loadStatus);
        serviceConfigPushFeedback.setUpdateTime(updateTime);
        serviceConfigPushFeedback.setLoadTime(loadTime);
        serviceConfigPushFeedback.setCreateTime(now);
        serviceConfigPushFeedback.setGrayGroupId(grayGroupId);

        this.serviceConfigPushFeedbackMapper.insert(serviceConfigPushFeedback);
        return serviceConfigPushFeedback;
    }
}
