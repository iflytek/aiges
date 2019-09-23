package com.iflytek.ccr.polaris.cynosure.dbtransactional;

import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.dbcondition.*;
import com.iflytek.ccr.polaris.cynosure.domain.*;
import com.iflytek.ccr.polaris.cynosure.request.cluster.AddClusterRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.cluster.CopyClusterRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.service.AddServiceRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.service.CopyServiceRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceversion.CopyAddServiceVersionRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.serviceversion.CopyServiceVersionRequestBody;
import com.iflytek.ccr.polaris.cynosure.util.MD5Util;
import com.iflytek.ccr.polaris.cynosure.util.PropUtil;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.ArrayList;
import java.util.Date;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * 拷贝和新增事务操作 Created by DELL-5490 on 2018/7/21.
 */
@Service
public class CopyAndAddTransactional extends BaseService {
	@Autowired
	private IServiceVersionCondition serviceVersionConditionImpl;

	@Autowired
	private IServiceConfigCondition serviceConfigConditionImpl;

	@Autowired
	private IGrayGroupCondition grayGroupConditionImpl;

	@Autowired
	private IServiceCondition serviceConditionImpl;

	@Autowired
	private IClusterCondition clusterConditionImpl;

	@Autowired
	private IClusterCondition clusterConditionimpl;

	@Autowired
	private IProjectCondition projectConditionimpl;

	@Autowired
	private PropUtil propUtil;

	/**
	 * 复制版本
	 *
	 * @param body
	 * @return
	 */
	@Transactional
	public ServiceVersion copyAndAddVersion(CopyServiceVersionRequestBody body, Date updateTime) {
		CopyAddServiceVersionRequestBody addServiceVersionRequestBody = new CopyAddServiceVersionRequestBody(body.getVersion(), body.getDesc(), body.getServiceId());
		addServiceVersionRequestBody.setUpdateTime(updateTime);
		String oldVersionId = body.getOldVersionId();
		// (一)往tb_service_version表中新增版本
		ServiceVersion serviceVersion = this.serviceVersionConditionImpl.add(addServiceVersionRequestBody);
		String newVersionId = serviceVersion.getId();

		// (二)根据oldVersionId查询tb_gray_group表，如果被复制的版本属于灰度版本，则需要对新的版本添加灰度组
		HashMap<String, Object> map = new HashMap<>();
		map.put("versionId", oldVersionId);
		List<GrayGroup> grayGroupList = this.grayGroupConditionImpl.findList(map);
		Map<String, String> oldGrayId2NewGrayId = new HashMap<>();

		if (null != grayGroupList && !grayGroupList.isEmpty()) {
			oldGrayId2NewGrayId = this.grayGroupConditionImpl.copy1(grayGroupList, newVersionId);
		}

		// (三)同步版本下的配置文件
		com.iflytek.ccr.polaris.cynosure.domain.Service service = this.serviceConditionImpl.findById(body.getServiceId());
		Cluster cluster = this.clusterConditionimpl.findById(service.getGroupId());
		Project project = this.projectConditionimpl.findById(cluster.getProjectId());
		String project_Group = project.getName() + cluster.getName();
		String project_GroupMD5 = MD5Util.getMD5(project_Group.getBytes());
		String service_Version = service.getName() + body.getVersion();
		String service_VersionMD5 = MD5Util.getMD5(service_Version.getBytes());
		String path = propUtil.CONFIG_PATH + project_GroupMD5 + "/" + service_VersionMD5 + "/";

		List<ServiceConfig> serviceConfigList = this.serviceConfigConditionImpl.find(oldVersionId, null, null);
		if (null != serviceConfigList && !serviceConfigList.isEmpty()) {
			this.serviceConfigConditionImpl.copyConfigs1(serviceConfigList, newVersionId, oldGrayId2NewGrayId, path);
		}
		return serviceVersion;
	}

	/**
	 * 复制服务
	 *
	 * @param body
	 * @return
	 */
	@Transactional
	public com.iflytek.ccr.polaris.cynosure.domain.Service copyAndAddService(CopyServiceRequestBody body) {

		AddServiceRequestBody addServiceRequestBody = new AddServiceRequestBody(body.getServiceName(), body.getDesc(), body.getClusterId());

		// (一)新增服务
		com.iflytek.ccr.polaris.cynosure.domain.Service service = this.serviceConditionImpl.add(addServiceRequestBody);
		String newServiceId = service.getId();

		// (二)根据被复制的服务的id查询tb_service_version表，如果该服务下存在版本，需要复制
		String oldServiceId = body.getOldServiceId();
		List<String> serviceIds = new ArrayList<>();
		serviceIds.add(oldServiceId);
		List<ServiceVersion> serviceVersionList = this.serviceVersionConditionImpl.findList(serviceIds);
		if (null != serviceVersionList && !serviceVersionList.isEmpty()) {
			for (ServiceVersion serviceVersion : serviceVersionList) {
				CopyServiceVersionRequestBody body1 = new CopyServiceVersionRequestBody(newServiceId, serviceVersion.getId(), serviceVersion.getVersion(), serviceVersion.getDescription());
				copyAndAddVersion(body1, serviceVersion.getUpdateTime());
			}
		}
		return service;
	}

	/**
	 * 复制集群
	 *
	 * @param body
	 * @return
	 */
	@Transactional
	public Cluster copyAndAddCluster(CopyClusterRequestBody body) {
		AddClusterRequestBody addClusterRequestBody = new AddClusterRequestBody(body.getClusterName(), body.getDesc(), body.getProjectId());

		// (一)新增集群
		Cluster cluster = this.clusterConditionImpl.add(addClusterRequestBody);
		String newClusterId = cluster.getId();

		// (二)根据被复制的集群id查询tb_service表，若该集群下存在服务，则复制服务
		String oldClusterId = body.getOldClusterId();
		List<String> clusterIds = new ArrayList<>();
		clusterIds.add(oldClusterId);
		List<com.iflytek.ccr.polaris.cynosure.domain.Service> serviceList = this.serviceConditionImpl.findList(clusterIds);
		if (null != serviceList && !serviceList.isEmpty()) {
			for (com.iflytek.ccr.polaris.cynosure.domain.Service service : serviceList) {
				CopyServiceRequestBody body1 = new CopyServiceRequestBody(newClusterId, service.getId(), service.getName(), service.getDescription());
				copyAndAddService(body1);
			}
		}

		return cluster;
	}
}
