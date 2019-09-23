package com.iflytek.ccr.polaris.cynosure.dbcondition.impl;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IServiceApiVersionCondition;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceApiVersion;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceVersion;
import com.iflytek.ccr.polaris.cynosure.mapper.ServiceApiVersionMapper;
import com.iflytek.ccr.polaris.cynosure.mapper.ServiceVersionMapper;
import com.iflytek.ccr.polaris.cynosure.request.servicediscovery.AddServiceApiVersionRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceversion.EditServiceVersionRequestBody;
import com.iflytek.ccr.polaris.cynosure.util.SnowflakeIdWorker;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.dao.DuplicateKeyException;
import org.springframework.stereotype.Service;

import java.util.Date;
import java.util.HashMap;
import java.util.List;

/**
 * 服务版本条件接口实现
 *
 * @author sctang2
 * @create 2017-12-10 16:25
 **/
@Service
public class ServiceApiVersionConditionImpl extends BaseService implements IServiceApiVersionCondition {
    private static final EasyLogger logger = EasyLoggerFactory.getInstance(ServiceApiVersionConditionImpl.class);

    @Autowired
    private ServiceVersionMapper    serviceVersionMapper;
    @Autowired
    private ServiceApiVersionMapper serviceApiVersionMapper;

    @Override
    public int findTotalCount(HashMap<String, Object> map) {
        return this.serviceApiVersionMapper.findTotalCount(map);
    }

    @Override
    public List<ServiceApiVersion> findList(HashMap<String, Object> map) {
        return this.serviceApiVersionMapper.findList(map);
    }

    @Override
    public List<ServiceApiVersion> findList(List<String> serviceIds) {
        HashMap<String, Object> map = new HashMap<>();
        map.put("serviceIds", serviceIds);
        return this.serviceApiVersionMapper.findServiceApiVersionList(map);
    }

    @Override
    public ServiceApiVersion add(AddServiceApiVersionRequestBody body) {
        String apiVersion = body.getApiVersion();
        String serviceId = body.getServiceId();
        String desc = body.getDesc();
        String userId = this.getUserId();

        //新增
        Date now = new Date();
        ServiceApiVersion serviceApiVersion = new ServiceApiVersion();
        serviceApiVersion.setId(SnowflakeIdWorker.getId());
        serviceApiVersion.setApiVersion(apiVersion);
        serviceApiVersion.setServiceId(serviceId);
        serviceApiVersion.setUserId(userId);
        serviceApiVersion.setDescription(desc);
        serviceApiVersion.setCreateTime(now);
        serviceApiVersion.setUpdateTime(now);
        try {
            this.serviceApiVersionMapper.insert(serviceApiVersion);
            return serviceApiVersion;
        } catch (DuplicateKeyException ex) {
            logger.warn("service version duplicate key " + ex.getMessage());
            return this.find(apiVersion, serviceId);
        }
    }

    @Override
    public ServiceApiVersion findServiceConfigListById(String id) {
        return this.serviceApiVersionMapper.findServiceConfigListById(id);
    }

    @Override
    public int deleteById(String id) {
        return this.serviceApiVersionMapper.deleteById(id);
    }

    @Override
    public ServiceApiVersion findById(String id) {
        return this.serviceApiVersionMapper.findById(id);
    }

    @Override
    public ServiceApiVersion updateById(String id, EditServiceVersionRequestBody body) {
        String desc = body.getDesc();

        //更新
        Date now = new Date();
        ServiceApiVersion serviceApiVersion = new ServiceApiVersion();
        serviceApiVersion.setId(id);
        serviceApiVersion.setDescription(desc);
        serviceApiVersion.setUpdateTime(now);
        this.serviceApiVersionMapper.updateById(serviceApiVersion);
        return serviceApiVersion;
    }

    @Override
    public ServiceApiVersion find(String version, String serviceId) {
        ServiceApiVersion serviceVersion = new ServiceApiVersion();
        serviceVersion.setApiVersion(version);
        serviceVersion.setServiceId(serviceId);
        return this.serviceApiVersionMapper.find(serviceVersion);
    }

    @Override
    public ServiceVersion findVersionJoinServiceJoinGroupJoinProjectById(String id) {
        return this.serviceVersionMapper.findVersionJoinServiceJoinGroupJoinProjectById(id);
    }
}
