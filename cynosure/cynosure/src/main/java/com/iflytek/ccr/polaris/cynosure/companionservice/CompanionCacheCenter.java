package com.iflytek.ccr.polaris.cynosure.companionservice;

import com.alibaba.fastjson.JSON;
import com.alibaba.fastjson.JSONArray;
import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.cynosure.companionservice.domain.Result;
import com.iflytek.ccr.polaris.cynosure.companionservice.domain.ServiceDetailNewResult;
import com.iflytek.ccr.polaris.cynosure.companionservice.domain.ServiceProviderConsumerResult;
import com.iflytek.ccr.polaris.cynosure.companionservice.domain.ServiceResult;
import com.iflytek.ccr.polaris.cynosure.companionservice.domain.response.ServiceDataDetailResponse;
import com.iflytek.ccr.polaris.cynosure.companionservice.domain.response.ServiceDataNewDetailResponse;
import com.iflytek.ccr.polaris.cynosure.companionservice.domain.response.ServiceNewApiVersionResponse;
import com.iflytek.ccr.polaris.cynosure.companionservice.domain.response.ServiceResponse;
import com.iflytek.ccr.polaris.cynosure.domain.Region;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceProviderInstanceConf;
import com.iflytek.ccr.polaris.cynosure.util.HttpUtil;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.CachePut;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.stereotype.Service;

import java.util.*;

/**
 * companion缓存中心
 *
 * @author  sctang2
 * @create 2018-02-06 12:04
 **/
@Service
public class CompanionCacheCenter extends BaseCenter {
    private static final EasyLogger logger = EasyLoggerFactory.getInstance(CompanionCacheCenter.class);
    /**
     * 查询集群下的服务
     *
     * @param path
     * @param regions
     * @return
     */
    //缓存位置名称value,不能为空,缓存的key默认为空，支持spEL，触发条件，默认为空，全部都加入缓存
    @Cacheable(value = "services", key = "'services:cluster_'.concat(#path).concat(#serviceName)")
    public ServiceResult findServices(String path, List<Region> regions, String serviceName) {
        //创建（查询）
        String querys = "?path=" + path;

        //批量get请求
//        List<Result> cacheCenterResults = this.batchGet(QUERY_SERVICE_PATH_URL, querys, regions);

        //使用新的请求方法
        List<Result> cacheCenterResults = new ArrayList<>();
        for (Region region: regions) {
            Result result = new Result();
            String regionName = region.getName();
            String host_Port = region.getPushUrl();
            String url = host_Port+QUERY_SERVICE_PATH_URL+querys;
            String dataStr = HttpUtil.myHttpGet(url);
            result.setName(regionName);
            result.setResult(dataStr);
            cacheCenterResults.add(result);
        }
        //解析服务参数
        Map<String, Object> resultMap = this.parseService(cacheCenterResults, serviceName);
        List<ServiceDetailNewResult> serviceDetailResults = (List<ServiceDetailNewResult>)resultMap.get("results");
        int failureCount = (int)resultMap.get("failureCount");
        List<String> failureArea = (List<String>)resultMap.get("failureArea");

        ServiceResult results = new ServiceResult();
        results.setResults(serviceDetailResults);
        results.setFailureCount(failureCount);
        results.setFailureArea(failureArea);

        if (null == serviceDetailResults || serviceDetailResults.isEmpty()) {
            results.setTotalCount(0);
        } else {
            results.setTotalCount(serviceDetailResults.size());
            logger.info(cacheCenterResults.toString());
        }
        return results;
    }

    /**
     * 同步提供方
     *
     * @param path
     * @param region
     * @param type
     * @return
     */
    //value缓存的名称，在spring配置文件中定义
    @CacheEvict(value = "services", key = "'region_'.concat(#region.name).concat('_services_').concat(#type).concat(':cluster_').concat(#path)")
    public ServiceResult syncProviders(String path, Region region, String type) {
        return null;
    }

    /**
     * 查询提供端、消费端列表
     *
     * @param path
     * @param region
     * @param type, key = "'region_'.concat(#region.name).concat('_services_').concat(#type).concat(':cluster_').concat(#path)"
     * @return
     */
//    @Cacheable(value = "services", key = "'region_'.concat(#region.name).concat('_services_').concat(#type).concat(':cluster_').concat(#path)")
    public ServiceResult findProviderConsumers(String path, Region region, String type) {
        //创建
        String querys = "?path=" + path;
        //get请求
        Result cacheCenterResult;
        if ("provider".equals(type.toLowerCase())) {
            cacheCenterResult = this.get(QUERY_SERVICE_PROVIDER_URL, querys, region);
            //解析提供端参数
            List<ServiceProviderInstanceConf> serviceProviderInstanceConfResults = this.parseProviderConsumer(cacheCenterResult);
            ServiceResult results = new ServiceResult();
            results.setResults(serviceProviderInstanceConfResults);
            if (null == serviceProviderInstanceConfResults || serviceProviderInstanceConfResults.isEmpty()) {
                results.setTotalCount(0);
            } else {
                results.setTotalCount(serviceProviderInstanceConfResults.size());
            }
            return results;
        } else {
            cacheCenterResult = this.get(QUERY_SERVICE_CONSUMER_URL, querys, region);
            //消费端参数
            List<ServiceProviderConsumerResult> serviceConsumerResults = this.parseConsumer(cacheCenterResult);
            ServiceResult results = new ServiceResult();
            results.setResults(serviceConsumerResults);
            if (null == serviceConsumerResults || serviceConsumerResults.isEmpty()) {
                results.setTotalCount(0);
            } else {
                results.setTotalCount(serviceConsumerResults.size());
            }
            return results;
        }


    }

    /**
     * 查询消费端列表
     *
     * @param path
     * @param region
     * @param type
     * @return
     */
    @Cacheable(value = "configs", key = "'region_'.concat(#region.name).concat('_configs_').concat(#type).concat(':cluster_').concat(#path)")
    public ServiceResult findConfigsConsumers(String path, Region region, String type) {
        //创建
        String querys = "?path=" + path;

        //get请求
        Result cacheCenterResult = null;
        if ("consumer".equals(type.toLowerCase())) {
            cacheCenterResult = this.get(Query_CONFIG_CONSUMER_URL, querys, region);
        }

        //解析提供端、消费端参数
        List<ServiceProviderConsumerResult> serviceProviderConsumerResults = this.parseConsumer(cacheCenterResult);
        ServiceResult results = new ServiceResult();
        results.setResults(serviceProviderConsumerResults);
        if (null == serviceProviderConsumerResults || serviceProviderConsumerResults.isEmpty()) {
            results.setTotalCount(0);
        } else {
            results.setTotalCount(serviceProviderConsumerResults.size());
        }
        return results;
    }

    /**
     * 解析服务参数
     *
     * @param cacheCenterResults
     * @return
     */
    private Map<String, Object> parseService(List<Result> cacheCenterResults, String serviceName) {
        List<ServiceDetailNewResult> results = new ArrayList<>();
        //如果要解析的结果为null或者为空，直接返回null
        if (cacheCenterResults==null||cacheCenterResults.isEmpty()){
            return null;
        }

        int failureCount = 0;
        List<String> failureArea = new ArrayList<>();

        //筛选符合我要的服务名称path
        for (Result result : cacheCenterResults) {
            String resultStr = result.getResult();
            String regionName = result.getName();
            if ("error".equals(resultStr)||"connection time out".equals(resultStr)){
                    failureCount++;
                    failureArea.add(regionName);
            }else {
                ServiceNewApiVersionResponse cacheServiceResponse;
                cacheServiceResponse = JSON.parseObject(resultStr, ServiceNewApiVersionResponse.class);
                if (null == cacheServiceResponse || 0 != cacheServiceResponse.getRet() || null == cacheServiceResponse.getData() || null == cacheServiceResponse.getData().getPathList() || cacheServiceResponse.getData().getPathList().isEmpty()) {
                    continue;
                }
                List<ServiceDataNewDetailResponse> pathList = cacheServiceResponse.getData().getPathList();
                for (ServiceDataNewDetailResponse detailResponse : pathList) {
                    ServiceDetailNewResult newResult = new ServiceDetailNewResult();
                    String name = result.getName();
                    newResult.setName(name);
                    newResult.setVersionList(detailResponse.getVersionList());
                    if (serviceName.equals(detailResponse.getPath())) {
                        results.add(newResult);
                    }
                }
            }
        }

        //将请求得到的正常结果、失败区域数目、失败区域名字放入到map中返回
        Map<String, Object> resultMap = new HashMap<>();
        resultMap.put("results", results);
        resultMap.put("failureCount", failureCount);
        resultMap.put("failureArea", failureArea);

        return resultMap;
    }

    /**
     * 解析提供端、消费端参数
     *{
     "name": "test",
     "result": "{\"ret\":0,\"msg\":\"success\",\"data\":{\"pathList\":[{\"childPath\":\"80.1.86.21:311\",\"pushId\":\"3787528462135197696\",\"data\":\"{\\\"sdk\\\":{\\\"is_valid\\\":true},\\\"user\\\":{\\\"key1\\\":\\\"value1\\\",\\\"key2\\\":\\\"value2\\\",\\\"weight\\\":10}}\"}]}}"
     }
     * @param cacheCenterResult
     * @return
     */
    private List<ServiceProviderInstanceConf> parseProviderConsumer(Result cacheCenterResult) {
        List<ServiceProviderInstanceConf> results = new ArrayList<>();
        ServiceResponse cacheServiceResponse = JSON.parseObject(cacheCenterResult.getResult(), ServiceResponse.class);
        if (null == cacheServiceResponse || 0 != cacheServiceResponse.getRet() || null == cacheServiceResponse.getData() || null == cacheServiceResponse.getData().getPathList() || cacheServiceResponse.getData().getPathList().isEmpty()) {
            return new ArrayList<>();
        }
        List<ServiceDataDetailResponse> pathList = cacheServiceResponse.getData().getPathList();
        pathList.forEach(x -> {
            ServiceProviderInstanceConf result = JSON.parseObject(x.getData(), ServiceProviderInstanceConf.class);
//            ServiceProviderConsumerResult result = JSON.parseObject(x.getData(), ServiceProviderConsumerResult.class);
            if (null==result) {
//                result = new ServiceProviderConsumerResult();
                result = new ServiceProviderInstanceConf();
                Map<String, Object> user = new HashMap<>();
                user.put("weight", 1);
                result.setUser(user);
//                result.setIs_valid(true);
//                result.setWeight(100 / size);
                Map<String, Object> sdk = new HashMap<>();
                sdk.put("is_valid", true);
                result.setSdk(sdk);
            }
            result.setAddr(x.getChildPath());
            results.add(result);
        });
        //按照服务地址排序
        Collections.sort(results, Comparator.comparing(ServiceProviderInstanceConf::getAddr));
        return results;
    }

    /**
     * 解析消费端参数
     *
     * @param cacheCenterResult
     * @return
     */
    private List<ServiceProviderConsumerResult> parseConsumer(Result cacheCenterResult) {
        List<ServiceProviderConsumerResult> results = new ArrayList<>();
        ServiceResponse cacheServiceResponse = JSON.parseObject(cacheCenterResult.getResult(), ServiceResponse.class);
        if (null == cacheServiceResponse || 0 != cacheServiceResponse.getRet() || null == cacheServiceResponse.getData() || null == cacheServiceResponse.getData().getPathList() || cacheServiceResponse.getData().getPathList().isEmpty()) {
            return new ArrayList<>();
        }
        List<ServiceDataDetailResponse> pathList = cacheServiceResponse.getData().getPathList();
        int size = pathList.size();
        pathList.forEach(x -> {
            ServiceProviderConsumerResult result = JSON.parseObject(x.getData(), ServiceProviderConsumerResult.class);
            if (null == result) {
                result = new ServiceProviderConsumerResult();
                result.setIs_valid(true);
                result.setWeight(100 / size);
            }
            result.setAddr(x.getChildPath());
            results.add(result);
        });
        //按照服务地址排序
        Collections.sort(results, Comparator.comparing(ServiceProviderConsumerResult::getAddr));
        return results;
    }
}
