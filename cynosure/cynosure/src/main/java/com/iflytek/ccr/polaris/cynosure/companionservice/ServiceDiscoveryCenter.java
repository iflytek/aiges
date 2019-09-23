package com.iflytek.ccr.polaris.cynosure.companionservice;

import com.alibaba.fastjson.JSON;
import com.alibaba.fastjson.JSONArray;
import com.alibaba.fastjson.JSONObject;
import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.cynosure.companionservice.domain.*;
import com.iflytek.ccr.polaris.cynosure.companionservice.domain.response.ServiceNewResponse;
import com.iflytek.ccr.polaris.cynosure.domain.Region;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceProviderInstanceConf;
import com.iflytek.ccr.polaris.cynosure.network.CustomHttpParams;
import com.iflytek.ccr.polaris.cynosure.request.servicediscovery.RouteRule;
import com.iflytek.ccr.polaris.cynosure.util.MD5Util;
import com.iflytek.ccr.polaris.cynosure.util.PropUtil;
import com.iflytek.ccr.polaris.cynosure.util.SnowflakeIdWorker;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.List;

/**
 * 服务发现中心
 *
 * @author sctang2
 * @create 2017-12-10 20:08
 **/
@Service
public class ServiceDiscoveryCenter extends BaseCenter {
	private final EasyLogger logger = EasyLoggerFactory.getInstance(ServiceDiscoveryCenter.class);
	@Autowired
	private PropUtil propUtil;

	@Autowired
	private CompanionCacheCenter companionCacheCenter;

	/**
	 * 获取服务path
	 *
	 * @param project
	 * @param group
	 * @return
	 */
	public String getServicePath(String project, String group) {
		String projectJointGroup = project + group;
		String projectJointGroupMD5 = MD5Util.getMD5(projectJointGroup.getBytes());
		return propUtil.SERVICE_PATH + projectJointGroupMD5;
	}

	/**
	 * 获取服务配置path
	 *
	 * @param project
	 * @param group
	 * @param service
	 * @return
	 */
	public String getConfigNewPath(String project, String group, String service, String apiVersion) {
		String servicePath = this.getServicePath(project, group);
		return servicePath + "/" + service + "/" + apiVersion;
	}

	/**
	 * 获取服务提供端path
	 *
	 * @param project
	 * @param group
	 * @param service
	 * @return
	 */
	public String getServiceProviderPath(String project, String group, String service, String apiVersion) {
		String servicePath = this.getServicePath(project, group);
		return servicePath + "/" + service + "/" + apiVersion + "/provider";
	}

	/**
	 * 获取服务提供端ip path
	 *
	 * @param project
	 * @param group
	 * @param service
	 * @param addr
	 * @return
	 */
	public String getServiceProviderAddrPath(String project, String group, String service, String version, String addr) {
		String providerPath = this.getServiceProviderPath(project, group, service, version);
		return providerPath + "/" + addr;
	}

	/**
	 * 获取服务消费者path
	 *
	 * @param project
	 * @param group
	 * @param service
	 * @return
	 */
	public String getServiceConsumerPath(String project, String group, String service, String apiVersion) {
		String servicePath = this.getServicePath(project, group);
		return servicePath + "/" + service + "/" + apiVersion + "/consumer";
	}

	/**
	 * 分页查询服务下的版本
	 *
	 * @param path
	 * @param regions
	 * @param serviceName
	 * @param isPage
	 * @param startIndex
	 * @param endIndex
	 * @return
	 */
	public ServiceResult findServicesByPaging(String path, List<Region> regions, String serviceName, int isPage, int startIndex, int endIndex) {
		// 查询服务下的版本
		ServiceResult serviceResult = this.companionCacheCenter.findServices(path, regions, serviceName);
		return serviceResult;
	}

	/**
	 * 同步提供方
	 *
	 * @param path
	 * @param region
	 * @param isProvider
	 * @return
	 */
	public ServiceResult syncProviders(String path, Region region, boolean isProvider) {
		String type = isProvider ? "provider" : "consumer";
		return this.companionCacheCenter.syncProviders(path, region, type);
	}

	/**
	 * 查询服务配置信息
	 *
	 * @param path
	 * @param region
	 * @return
	 */
	public ServiceConfRuleResult findServiceConf(String path, Region region) {
		// 创建
		String querys = "?path=" + path;

		// get请求
		Result cacheCenterResult = this.get(QUERY_SERVICE_CONFIG_PATH_URL, querys, region);

		// 解析服务配置参数
		ServiceConfRuleResult result = this.parseConf(cacheCenterResult);
		return result;
	}

	/**
	 * 查询提供端、消费端列表
	 *
	 * @param path
	 * @param region
	 * @param isProvider
	 * @param startIndex
	 * @param endIndex
	 * @return
	 */
	public ServiceResult findProviderConsumersByPaging(String path, Region region, boolean isProvider, int startIndex, int endIndex, int isPage) {
		// 查询提供端、消费端列表
		String type = isProvider ? "provider" : "consumer";

		ServiceResult serviceResult = this.companionCacheCenter.findProviderConsumers(path, region, type);
		int totalCount = serviceResult.getTotalCount();
		if (0 == totalCount) {
			return serviceResult;
		}
		ServiceResult newServiceResult = new ServiceResult();
		if ("provider".equals(type)) {
			List<ServiceProviderInstanceConf> serviceProviderResults = serviceResult.getResults();
			if (startIndex >= totalCount) {
				startIndex = 0;
			}
			if (endIndex >= totalCount) {
				endIndex = totalCount;
			}
			List<ServiceProviderInstanceConf> newServiceProviderResults = serviceProviderResults;
			// 判断是否分页
			if (isPage == 1) {
				newServiceProviderResults = serviceProviderResults.subList(startIndex, endIndex);
			}

			newServiceResult.setTotalCount(totalCount);
			newServiceResult.setResults(newServiceProviderResults);
			return newServiceResult;
		} else {
			List<ServiceProviderConsumerResult> serviceConsumerResults = serviceResult.getResults();
			if (startIndex >= totalCount) {
				startIndex = 0;
			}
			if (endIndex >= totalCount) {
				endIndex = totalCount;
			}
			List<ServiceProviderConsumerResult> newServiceConsumerResults = serviceConsumerResults;
			// 判断是否分页
			if (isPage == 1) {
				newServiceConsumerResults = serviceConsumerResults.subList(startIndex, endIndex);
			}

			newServiceResult.setTotalCount(totalCount);
			newServiceResult.setResults(newServiceConsumerResults);
			return newServiceResult;
		}

	}

	/**
	 * 通过一对一推送
	 *
	 * @param path
	 * @param bt
	 * @param region
	 * @param isTemporaryNode
	 * @return
	 */
	public PushResult pushByOneToOne(String path, byte[] bt, Region region, boolean isTemporaryNode) {
		String pushId = SnowflakeIdWorker.getId();
		// 创建
		CustomHttpParams params = this.create(path, pushId, bt);

		// 批量post请求
		List<Region> regionList = new ArrayList<>();
		regionList.add(region);
		String url;
		if (isTemporaryNode) {
			url = PUSH_INSTANCE_CONFIG_URL;
		} else {
			url = PUSH_CLUSTER_CONFIG_URL;
		}
		List<Result> cacheCenterResults = this.batchPost(url, params, regionList);

		// 解析推送参数
		List<PushDetailResult> cacheCenterPushDetailResults = this.parsePush(cacheCenterResults);
		PushResult result = new PushResult();
		result.setPushId(pushId);
		result.setResult(JSONArray.toJSONString(cacheCenterPushDetailResults));
		return result;
	}

	/**
	 * 通过一对一推送
	 *
	 * @param path
	 * @param bt
	 * @param region
	 * @param isTemporaryNode
	 * @return
	 */
	public PushResult pushNewByOneToOne(String path, byte[] bt, Region region, JSONArray sdk, JSONObject user, boolean isTemporaryNode) {
		String pushId = SnowflakeIdWorker.getId();
		// 创建
		CustomHttpParams params = this.create(path, pushId, sdk, user, bt);

		// 批量post请求
		List<Region> regionList = new ArrayList<>();
		regionList.add(region);
		String url;
		if (isTemporaryNode) {
			url = PUSH_INSTANCE_CONFIG_URL;
		} else {
			url = PUSH_SERVICE_CLUSTER_CONFIG_URL;
		}
		List<Result> cacheCenterResults = this.batchPost(url, params, regionList);

		// 解析推送参数
		List<PushDetailResult> cacheCenterPushDetailResults = this.parsePush(cacheCenterResults);
		PushResult result = new PushResult();
		result.setPushId(pushId);
		result.setResult(JSONArray.toJSONString(cacheCenterPushDetailResults));
		return result;
	}

	/**
	 * 解析服务配置参数(服务发现新版本)
	 *
	 * @param cacheCenterResult
	 * @return
	 */
	private ServiceConfRuleResult parseConf(Result cacheCenterResult) {
		ServiceConfRuleResult result = new ServiceConfRuleResult();
		ServiceNewResponse cacheServiceResponse = JSON.parseObject(cacheCenterResult.getResult(), ServiceNewResponse.class);
		if (null == cacheServiceResponse || 0 != cacheServiceResponse.getRet() || null == cacheServiceResponse.getData()) {
			return result;
		}
		String sdkData = cacheServiceResponse.getData().getSdk().getData();
		if (null != cacheServiceResponse.getData().getSdk().getPushId() && null != cacheServiceResponse.getData().getSdk().getPushId()) {
			List<RouteRule> sdk = JSONObject.parseArray(sdkData, RouteRule.class);
			result.setSdk(sdk);
		}
		String userData = cacheServiceResponse.getData().getUser().getData();
		if (null != cacheServiceResponse.getData().getUser().getPushId() && null != cacheServiceResponse.getData().getUser().getPushId()) {
			result.setUser(JSON.parseObject(userData));
		}
		result.setName(cacheCenterResult.getName());
		return result;
	}
}
