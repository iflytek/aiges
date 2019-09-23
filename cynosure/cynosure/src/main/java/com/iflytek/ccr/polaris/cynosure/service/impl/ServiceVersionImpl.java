package com.iflytek.ccr.polaris.cynosure.service.impl;

import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.customdomain.SearchCondition;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IServiceCondition;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IServiceVersionCondition;
import com.iflytek.ccr.polaris.cynosure.dbtransactional.CopyAndAddTransactional;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceConfig;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceVersion;
import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.request.IdRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceversion.AddServiceVersionRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceversion.CopyServiceVersionRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceversion.EditServiceVersionRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceversion.QueryServiceVersionRequestBody;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.serviceversion.ServiceVersionDetailResponseBody;
import com.iflytek.ccr.polaris.cynosure.service.ILastestSearchService;
import com.iflytek.ccr.polaris.cynosure.service.IServiceVersion;
import com.iflytek.ccr.polaris.cynosure.util.PagingUtil;
import org.apache.commons.lang3.StringUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Optional;

/**
 * 服务历史业务接口实现
 *
 * @author sctang2
 * @create 2017-11-17 16:03
 **/
@Service
public class ServiceVersionImpl extends BaseService implements IServiceVersion {
	@Autowired
	private IServiceVersionCondition serviceVersionConditionImpl;

	@Autowired
	private IServiceCondition serviceConditionimpl;

	@Autowired
	private ILastestSearchService lastestSearchServiceImpl;

	@Autowired
	private CopyAndAddTransactional copyAndAddTransactional;

	@Override
	public Response<QueryPagingListResponseBody> findLastestList(QueryServiceVersionRequestBody body) {
		QueryPagingListResponseBody result;
		String projectName = body.getProject();
		String clusterName = body.getCluster();
		String serviceName = body.getService();
		// 查询最近的搜索
		SearchCondition searchCondition = this.lastestSearchServiceImpl.find(projectName, clusterName, serviceName);
		projectName = searchCondition.getProject();
		clusterName = searchCondition.getCluster();
		serviceName = searchCondition.getService();
		if (StringUtils.isBlank(projectName) || StringUtils.isBlank(clusterName) || StringUtils.isBlank(serviceName)) {
			result = PagingUtil.createResult(body, 0);
			return new Response<>(result);
		}

		// 创建分页查询条件
		HashMap<String, Object> map = PagingUtil.createCondition(body);
		map.put("projectName", projectName);
		map.put("clusterName", clusterName);
		map.put("serviceName", serviceName);

		// 查询总数
		int totalCount = this.serviceVersionConditionImpl.findTotalCount(map);

		// 创建分页结果
		result = PagingUtil.createResult(body, totalCount);

		// 保存最近的搜索
		String condition = this.lastestSearchServiceImpl.saveLastestSearch(projectName, clusterName, serviceName);
		result.setCondition(condition);
		if (0 == totalCount) {
			return new Response<>(result);
		}

		// 查询列表
		List<ServiceVersionDetailResponseBody> list = new ArrayList<>();
		Optional<List<ServiceVersion>> serviceVersionList = Optional.ofNullable(this.serviceVersionConditionImpl.findList(map));
		serviceVersionList.ifPresent(x -> {
			x.forEach(y -> {
				// 创建版本结果
				ServiceVersionDetailResponseBody serviceVersionDetail = this.createServiceVersionResult(y);
				list.add(serviceVersionDetail);
			});
		});
		result.setList(list);
		return new Response<>(result);
	}

	@Override
	public Response<QueryPagingListResponseBody> findList(QueryServiceVersionRequestBody body) {
		String projectName = body.getProject();
		String clusterName = body.getCluster();
		String serviceName = body.getService();
		if (StringUtils.isBlank(projectName) || StringUtils.isBlank(clusterName) || StringUtils.isBlank(serviceName)) {
			return new Response<>(PagingUtil.createResult(body, 0));
		}

		// 创建分页查询条件
		HashMap<String, Object> map = PagingUtil.createCondition(body);
		map.put("projectName", projectName);
		map.put("clusterName", clusterName);
		map.put("serviceName", serviceName);

		// 查询总数
		int totalCount = this.serviceVersionConditionImpl.findTotalCount(map);

		// 创建分页结果
		QueryPagingListResponseBody result = PagingUtil.createResult(body, totalCount);
		if (0 == totalCount) {
			return new Response<>(result);
		}

		// 查询列表
		List<ServiceVersionDetailResponseBody> list = new ArrayList<>();
		Optional<List<ServiceVersion>> serviceVersionList = Optional.ofNullable(this.serviceVersionConditionImpl.findList(map));
		serviceVersionList.ifPresent(x -> {
			x.forEach(y -> {
				// 创建版本结果
				ServiceVersionDetailResponseBody serviceVersionDetail = this.createServiceVersionResult(y);
				list.add(serviceVersionDetail);
			});
		});
		result.setList(list);
		return new Response<>(result);
	}

	@Override
	public Response<ServiceVersionDetailResponseBody> find(String id) {
		// 根据id查询版本
		ServiceVersion serviceVersion = this.serviceVersionConditionImpl.findById(id);
		if (null == serviceVersion) {
			// 不存在该版本
			return new Response<>(SystemErrCode.ERRCODE_SERVICE_VERSION_NOT_EXISTS, SystemErrCode.ERRMSG_SERVICE_VERSION_NOT_EXISTS);
		}

		// 创建版本结果
		ServiceVersionDetailResponseBody result = this.createServiceVersionResult(serviceVersion);
		return new Response<>(result);
	}

	@Override
	public Response<ServiceVersionDetailResponseBody> add(AddServiceVersionRequestBody body) {
		// 通过id查询服务
		String serviceId = body.getServiceId();
		com.iflytek.ccr.polaris.cynosure.domain.Service service = this.serviceConditionimpl.findById(serviceId);
		if (null == service) {
			// 不存在该服务
			return new Response<>(SystemErrCode.ERRCODE_SERVICE_NOT_EXISTS, SystemErrCode.ERRMSG_SERVICE_NOT_EXISTS);
		}

		// 根据版本和服务id查询服务版本
		String version = body.getVersion();
		ServiceVersion serviceVersion = this.serviceVersionConditionImpl.find(version, serviceId);
		if (null != serviceVersion) {
			return new Response<>(SystemErrCode.ERRCODE_SERVICE_VERSION_EXISTS, SystemErrCode.ERRMSG_SERVICE_VERSION_EXISTS);
		}

		// 创建版本
		ServiceVersion newServiceVersion = this.serviceVersionConditionImpl.add(body);

		// 创建版本结果
		ServiceVersionDetailResponseBody result = this.createServiceVersionResult(newServiceVersion);
		return new Response<>(result);
	}

	@Override
	public Response<ServiceVersionDetailResponseBody> edit(EditServiceVersionRequestBody body) {
		// 根据id查询版本
		String id = body.getId();
		ServiceVersion serviceVersion = this.serviceVersionConditionImpl.findById(id);
		if (null == serviceVersion) {
			return new Response<>(SystemErrCode.ERRCODE_SERVICE_VERSION_NOT_EXISTS, SystemErrCode.ERRMSG_SERVICE_VERSION_NOT_EXISTS);
		}

		// 根据id更新版本
		ServiceVersion updateServiceVersion = this.serviceVersionConditionImpl.updateById(id, body);

		// 创建服务版本结果
		updateServiceVersion.setServiceId(id);
		updateServiceVersion.setCreateTime(serviceVersion.getCreateTime());
		updateServiceVersion.setVersion(serviceVersion.getVersion());
		ServiceVersionDetailResponseBody result = this.createServiceVersionResult(updateServiceVersion);
		return new Response<>(result);
	}

	@Override
	public Response<String> delete(IdRequestBody body) {
		// 通过id查询服务配置列表
		String id = body.getId();
		ServiceVersion serviceVersion = this.serviceVersionConditionImpl.findServiceConfigListById(id);
		if (null == serviceVersion) {
			// 不存在该版本
			return new Response<>(SystemErrCode.ERRCODE_SERVICE_VERSION_NOT_EXISTS, SystemErrCode.ERRMSG_SERVICE_VERSION_NOT_EXISTS);
		}

		List<ServiceConfig> serviceConfigList = serviceVersion.getServiceConfigList();
		if (null != serviceConfigList && !serviceConfigList.isEmpty()) {
			// 已创建配置
			return new Response<>(SystemErrCode.ERRCODE_SERVICE_CONFIG_CREATE, SystemErrCode.ERRMSG_SERVICE_CONFIG_CREATE);
		}

		// 根据id删除版本
		this.serviceVersionConditionImpl.deleteById(id);
		return new Response<>(null);
	}

	@Override
	public Response<ServiceVersionDetailResponseBody> copy(CopyServiceVersionRequestBody body) {
		// 根据服务id查询复制后的版本所属的服务是否存在，若不存在，直接返回
		String serviceId = body.getServiceId();
		com.iflytek.ccr.polaris.cynosure.domain.Service service = this.serviceConditionimpl.findById(serviceId);
		if (null == service) {
			return new Response<>(SystemErrCode.ERRCODE_SERVICE_NOT_EXISTS, SystemErrCode.ERRMSG_SERVICE_NOT_EXISTS);
		}

		// 根据新的版本version和所属服务的id查询该版本是否已经存在，若存在，直接返回
		String version = body.getVersion();
		ServiceVersion serviceVersion = this.serviceVersionConditionImpl.find(version, serviceId);
		if (null != serviceVersion) {
			return new Response<>(SystemErrCode.ERRCODE_SERVICE_VERSION_EXISTS, SystemErrCode.ERRMSG_SERVICE_VERSION_EXISTS);
		}

		// 根据版本id查询被复制的版本是否存在，若不存在，直接返回
		String oldVersionId = body.getOldVersionId();
		ServiceVersion versionCopy = this.serviceVersionConditionImpl.findById(oldVersionId);
		if (null == versionCopy) {
			return new Response<>(SystemErrCode.ERRCODE_SERVICE_VERSION_NOT_EXISTS, SystemErrCode.ERRMSG_SERVICE_VERSION_NOT_EXISTS);
		}
		
		// 复制版本
		ServiceVersion newServiceVersion = this.copyAndAddTransactional.copyAndAddVersion(body,versionCopy.getUpdateTime());

		// 版本结果
		ServiceVersionDetailResponseBody result = this.createServiceVersionResult(newServiceVersion);
		return new Response<>(result);
	}
}