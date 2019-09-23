package com.iflytek.ccr.polaris.cynosure.dbcondition.impl;

import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IServiceVersionCondition;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceVersion;
import com.iflytek.ccr.polaris.cynosure.mapper.ServiceVersionMapper;
import com.iflytek.ccr.polaris.cynosure.request.serviceversion.AddServiceVersionRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceversion.CopyAddServiceVersionRequestBody;
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
public class ServiceVersionConditionImpl extends BaseService implements IServiceVersionCondition {
	private static final EasyLogger logger = EasyLoggerFactory.getInstance(ServiceVersionConditionImpl.class);

	@Autowired
	private ServiceVersionMapper serviceVersionMapper;

	@Override
	public int findTotalCount(HashMap<String, Object> map) {
		return this.serviceVersionMapper.findTotalCount(map);
	}

	@Override
	public List<ServiceVersion> findList(HashMap<String, Object> map) {
		return this.serviceVersionMapper.findList(map);
	}

	@Override
	public List<ServiceVersion> findList(List<String> serviceIds) {
		HashMap<String, Object> map = new HashMap<>();
		map.put("serviceIds", serviceIds);
		return this.serviceVersionMapper.findServiceVersionList(map);
	}

	@Override
	public ServiceVersion add(AddServiceVersionRequestBody body) {
		String version = body.getVersion();
		String serviceId = body.getServiceId();
		String desc = body.getDesc();
		String userId = this.getUserId();

		// 新增
		Date now = new Date();
		ServiceVersion serviceVersion = new ServiceVersion();
		serviceVersion.setId(SnowflakeIdWorker.getId());
		serviceVersion.setVersion(version);
		serviceVersion.setServiceId(serviceId);
		serviceVersion.setUserId(userId);
		serviceVersion.setDescription(desc);
		serviceVersion.setCreateTime(now);
		serviceVersion.setUpdateTime(now);
		if (body instanceof CopyAddServiceVersionRequestBody) {
			Date updateTime = ((CopyAddServiceVersionRequestBody) body).getUpdateTime();
			serviceVersion.setUpdateTime(updateTime == null ? now : updateTime);
		}
		try {
			this.serviceVersionMapper.insert(serviceVersion);
			return serviceVersion;
		} catch (DuplicateKeyException ex) {
			logger.warn("service version duplicate key " + ex.getMessage());
			return this.find(version, serviceId);
		}
	}

	@Override
	public ServiceVersion findServiceConfigListById(String id) {
		return this.serviceVersionMapper.findServiceConfigListById(id);
	}

	@Override
	public int deleteById(String id) {
		return this.serviceVersionMapper.deleteById(id);
	}

	@Override
	public ServiceVersion findById(String id) {
		return this.serviceVersionMapper.findById(id);
	}

	@Override
	public ServiceVersion updateById(String id, EditServiceVersionRequestBody body) {
		String desc = body.getDesc();

		// 更新
		Date now = new Date();
		ServiceVersion serviceVersion = new ServiceVersion();
		serviceVersion.setId(id);
		serviceVersion.setDescription(desc);
		serviceVersion.setUpdateTime(now);
		this.serviceVersionMapper.updateById(serviceVersion);
		return serviceVersion;
	}

	@Override
	public ServiceVersion find(String version, String serviceId) {
		ServiceVersion serviceVersion = new ServiceVersion();
		serviceVersion.setVersion(version);
		serviceVersion.setServiceId(serviceId);
		return this.serviceVersionMapper.find(serviceVersion);
	}

	@Override
	public ServiceVersion findVersionJoinServiceJoinGroupJoinProjectById(String id) {
		return this.serviceVersionMapper.findVersionJoinServiceJoinGroupJoinProjectById(id);
	}
}
