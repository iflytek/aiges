package com.iflytek.ccr.polaris.cynosure.dbcondition.impl;

import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IServiceConfigHistoryCondition;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceConfig;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceConfigHistory;
import com.iflytek.ccr.polaris.cynosure.mapper.ServiceConfigHistoryMapper;
import com.iflytek.ccr.polaris.cynosure.util.DateUtil;
import com.iflytek.ccr.polaris.cynosure.util.SnowflakeIdWorker;
import com.iflytek.ccr.polaris.cynosure.util.TimeStampID;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.Date;
import java.util.HashMap;
import java.util.List;

/**
 * 服务配置历史条件接口实现
 *
 * @author sctang2
 * @create 2017-12-10 22:21
 **/
@Service
public class ServiceConfigHistoryConditionImpl extends BaseService implements IServiceConfigHistoryCondition {
    @Autowired
    private ServiceConfigHistoryMapper serviceConfigHistoryMapper;

    @Override
    public int findTotalCount(HashMap<String, Object> map) {
        return this.serviceConfigHistoryMapper.findTotalCount(map);
    }

    @Override
    public List<ServiceConfigHistory> findList(HashMap<String, Object> map) {
        return this.serviceConfigHistoryMapper.findList(map);
    }

    @Override
    public ServiceConfigHistory add(ServiceConfig serviceConfig) {
        //构造配置历史实体
        ServiceConfigHistory serviceConfigHistory = new ServiceConfigHistory(
                SnowflakeIdWorker.getId(), this.getUserId(), serviceConfig.getId(),
                serviceConfig.getDescription(), TimeStampID.getStampID(), new Date(),
                serviceConfig.getContent(), serviceConfig.getMd5());
        //新增
        this.serviceConfigHistoryMapper.insert(serviceConfigHistory);
        return serviceConfigHistory;
    }

    @Override
    public List<ServiceConfigHistory> batchAdd(List<ServiceConfig> serviceConfigList) {
        Date createTime = new Date();
        String userId = this.getUserId();
        String pushVersion = TimeStampID.getStampID();

        //配置历史列表
        List<ServiceConfigHistory> serviceConfigHistoryList = new ArrayList<>();

        //循环加入到配置历史列表中
        for (ServiceConfig serviceConfig : serviceConfigList) {
            ServiceConfigHistory serviceConfigHistory = new ServiceConfigHistory(
                    SnowflakeIdWorker.getId(), userId, serviceConfig.getId(),
                    serviceConfig.getDescription(), pushVersion, createTime,
                    serviceConfig.getContent(), serviceConfig.getMd5());
            serviceConfigHistoryList.add(serviceConfigHistory);
        }

        //批量插入数据库
        this.serviceConfigHistoryMapper.batchInsert(serviceConfigHistoryList);
        return serviceConfigHistoryList;
    }

    @Override
    public List<ServiceConfigHistory> findByIds(List<String> ids) {
        return this.serviceConfigHistoryMapper.findByIds(ids);
    }

    @Override
    public int deleteByConfigId(String configId) {
        return this.serviceConfigHistoryMapper.deleteByConfigId(configId);
    }

    @Override
    public int deleteByConfigIds(List<String> configIds) {
        return this.serviceConfigHistoryMapper.deleteByConfigIds(configIds);
    }

    @Override
    public int findGrayTotalCount(HashMap<String, Object> map) {
        return this.serviceConfigHistoryMapper.findGrayTotalCount(map);
    }

    @Override
    public List<ServiceConfigHistory> findGrayList(HashMap<String, Object> map) {
        return this.serviceConfigHistoryMapper.findGrayList(map);
    }
}
