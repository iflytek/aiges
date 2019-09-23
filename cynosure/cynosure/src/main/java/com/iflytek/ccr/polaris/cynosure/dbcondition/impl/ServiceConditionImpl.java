package com.iflytek.ccr.polaris.cynosure.dbcondition.impl;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IServiceCondition;
import com.iflytek.ccr.polaris.cynosure.mapper.ServiceMapper;
import com.iflytek.ccr.polaris.cynosure.request.service.AddServiceRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.service.EditServiceRequestBody;
import com.iflytek.ccr.polaris.cynosure.util.SnowflakeIdWorker;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.dao.DuplicateKeyException;
import org.springframework.stereotype.Service;

import java.util.Date;
import java.util.HashMap;
import java.util.List;

/**
 * 服务条件接口实现
 *
 * @author sctang2
 * @create 2017-12-10 15:24
 **/
@Service
public class ServiceConditionImpl extends BaseService implements IServiceCondition {
    private static final EasyLogger logger = EasyLoggerFactory.getInstance(ServiceConditionImpl.class);

    @Autowired
    private ServiceMapper serviceMapper;

    @Override
    public int findTotalCount(HashMap<String, Object> map) {
        return this.serviceMapper.findTotalCount(map);
    }

    @Override
    public List<com.iflytek.ccr.polaris.cynosure.domain.Service> findList(HashMap<String, Object> map) {
        return this.serviceMapper.findList(map);
    }

    @Override
    public List<com.iflytek.ccr.polaris.cynosure.domain.Service> findList(List<String> clusterIds) {
        HashMap<String, Object> map = new HashMap<>();
        map.put("clusterIds", clusterIds);
        return this.serviceMapper.findServiceList(map);
    }

    @Override
    public com.iflytek.ccr.polaris.cynosure.domain.Service add(AddServiceRequestBody body) {
        String name = body.getName();
        String desc = body.getDesc();
        String userId = this.getUserId();
        String clusterId = body.getClusterId();

        //新增
        Date now = new Date();
        com.iflytek.ccr.polaris.cynosure.domain.Service service = new com.iflytek.ccr.polaris.cynosure.domain.Service();
        service.setId(SnowflakeIdWorker.getId());
        service.setName(name);
        service.setDescription(desc);
        service.setUserId(userId);
        service.setGroupId(clusterId);
        service.setCreateTime(now);
        try {
            this.serviceMapper.insert(service);
            return service;
        } catch (DuplicateKeyException ex) {
            logger.warn("service duplicate key " + ex.getMessage());
            return this.find(name, clusterId);
        }
    }

    @Override
    public int deleteById(String id) {
        return this.serviceMapper.deleteById(id);
    }

    @Override
    public com.iflytek.ccr.polaris.cynosure.domain.Service findById(String id) {
        return this.serviceMapper.findById(id);
    }

    @Override
    public com.iflytek.ccr.polaris.cynosure.domain.Service updateById(String id, EditServiceRequestBody body) {
        String desc = body.getDesc();

        //更新
        Date now = new Date();
        com.iflytek.ccr.polaris.cynosure.domain.Service service = new com.iflytek.ccr.polaris.cynosure.domain.Service();
        service.setId(id);
        service.setDescription(desc);
        service.setUpdateTime(now);
        this.serviceMapper.updateById(service);
        return service;
    }

    @Override
    public com.iflytek.ccr.polaris.cynosure.domain.Service find(String name, String clusterId) {
        com.iflytek.ccr.polaris.cynosure.domain.Service service = new com.iflytek.ccr.polaris.cynosure.domain.Service();
        service.setGroupId(clusterId);
        service.setName(name);
        return this.serviceMapper.find(service);
    }

    @Override
    public com.iflytek.ccr.polaris.cynosure.domain.Service findServiceVersionListById(String id) {
        return this.serviceMapper.findServiceVersionListById(id);
    }

    @Override
    public com.iflytek.ccr.polaris.cynosure.domain.Service findServiceJoinGroupJoinProjectByServiceId(String serviceId) {
        HashMap<String, Object> map = new HashMap<>();
        map.put("serviceId", serviceId);
        return this.serviceMapper.findServiceJoinGroupJoinProjectByMap(map);
    }

    @Override
    public com.iflytek.ccr.polaris.cynosure.domain.Service findServiceJoinGroupJoinProjectByName(String project, String group, String service) {
        HashMap<String, Object> map = new HashMap<>();
        map.put("project", project);
        map.put("group", group);
        map.put("service", service);
        return this.serviceMapper.findServiceJoinGroupJoinProjectByMap(map);
    }
}
