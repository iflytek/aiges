package com.iflytek.ccr.polaris.cynosure.service.impl;

import com.alibaba.fastjson.JSON;
import com.alibaba.fastjson.JSONArray;
import com.alibaba.fastjson.JSONObject;
import com.iflytek.ccr.log.EasyLogger;
import com.iflytek.ccr.log.EasyLoggerFactory;
import com.iflytek.ccr.polaris.cynosure.base.BaseService;
import com.iflytek.ccr.polaris.cynosure.companionservice.ServiceDiscoveryCenter;
import com.iflytek.ccr.polaris.cynosure.companionservice.domain.*;
import com.iflytek.ccr.polaris.cynosure.customdomain.SearchCondition;
import com.iflytek.ccr.polaris.cynosure.dbcondition.ILoadbalanceCondition;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IRegionCondition;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IServiceDiscoveryPushFeedbackCondition;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IServiceDiscoveryPushHistoryCondition;
import com.iflytek.ccr.polaris.cynosure.dbtransactional.QuickStartTransactional;
import com.iflytek.ccr.polaris.cynosure.dbtransactional.ServiceDiscoveryApiVersionTransactional;
import com.iflytek.ccr.polaris.cynosure.domain.LoadBalance;
import com.iflytek.ccr.polaris.cynosure.domain.Region;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceDiscoveryPushHistory;
import com.iflytek.ccr.polaris.cynosure.domain.ServiceProviderInstanceConf;
import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.exception.GlobalExceptionUtil;
import com.iflytek.ccr.polaris.cynosure.request.servicediscovery.*;
import com.iflytek.ccr.polaris.cynosure.response.QueryPagingListResponseBody;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.response.servicediscovery.*;
import com.iflytek.ccr.polaris.cynosure.service.ILastestSearchService;
import com.iflytek.ccr.polaris.cynosure.service.IServiceDiscovery;
import com.iflytek.ccr.polaris.cynosure.util.PagingUtil;
import com.iflytek.ccr.polaris.cynosure.util.SnowflakeIdWorker;
import org.apache.commons.lang3.StringEscapeUtils;
import org.apache.commons.lang3.StringUtils;
import org.apache.http.client.fluent.Request;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.io.UnsupportedEncodingException;
import java.util.*;

/**
 * 服务发现业务逻辑接口实现
 *
 * @author sctang2
 * @create 2017-12-05 15:38
 **/
@Service
public class ServiceDiscoveryImpl extends BaseService implements IServiceDiscovery {
    private static final EasyLogger logger = EasyLoggerFactory.getInstance(ServiceDiscoveryImpl.class);

    @Autowired
    private QuickStartTransactional quickStartTransactional;

    @Autowired
    private ILastestSearchService lastestSearchServiceImpl;

    @Autowired
    private IRegionCondition regionConditionImpl;

    @Autowired
    private ILoadbalanceCondition loadbalanceConditionImpl;

    @Autowired
    private IServiceDiscoveryPushHistoryCondition serviceDiscoveryPushHistoryConditionImpl;

    @Autowired
    private IServiceDiscoveryPushFeedbackCondition serviceDiscoveryPushFeedbackConditionImpl;

    @Autowired
    private ServiceDiscoveryCenter serviceDiscoveryCenter;

    @Autowired
    private ServiceDiscoveryApiVersionTransactional serviceDiscoveryApiVersionTransactional;

    @Override
    public Response<AddApiVersionResponseBody> add(AddServiceDiscoveryRequestBody body) {
        return this.serviceDiscoveryApiVersionTransactional.addApiVersion(body, false);
    }

    @Override
    public Response<QueryPagingListResponseBody> findLastestList(QueryServiceDiscoveryListRequestBody body) {
        QueryPagingListResponseBody result;

        //查询区域列表
        List<Region> regionList = this.regionConditionImpl.findList(null);
        if (null == regionList || regionList.isEmpty()) {
            //result = PagingUtil.createResult(body, 0);
            return new Response<>(SystemErrCode.ERRCODE_REGION_NOT_EXISTS, SystemErrCode.ERRMSG_REGION_NOT_EXISTS);
        }

        //查询最近的搜索
        String projectName = body.getProject();
        String clusterName = body.getCluster();
        String serviceName = body.getService();
        SearchCondition searchCondition = this.lastestSearchServiceImpl.find(projectName, clusterName, serviceName);
        projectName = searchCondition.getProject();
        clusterName = searchCondition.getCluster();
        serviceName = searchCondition.getService();
        if (StringUtils.isBlank(projectName) || StringUtils.isBlank(clusterName) || StringUtils.isBlank(serviceName)) {
            result = PagingUtil.createResult(body, 0);
            return new Response<>(null);
        }

        //获取path
        String path = this.serviceDiscoveryCenter.getServicePath(projectName, clusterName);

        //查询服务下的版本
        int startIndex = PagingUtil.getStartIndex(body);
        int endIndex = PagingUtil.getEndIndex(body);
        int isPage = body.getIsPage();
        ServiceResult serviceResult = this.serviceDiscoveryCenter.findServicesByPaging(path, regionList, serviceName, isPage, startIndex, endIndex);
        int totalCount = serviceResult.getTotalCount();
        result = PagingUtil.createResult(body, totalCount);

        //保存最近的搜索
        String condition = this.lastestSearchServiceImpl.saveLastestSearch(projectName, clusterName, serviceName);
        result.setCondition(condition);
        if (0 == totalCount) {
            return new Response<>(null);
        }

        //创建服务发现列表
        List<ServiceDetailNewResult> serviceDetailResults = serviceResult.getResults();
        if (null != serviceDetailResults && !serviceDetailResults.isEmpty()) {
            result.setList(serviceDetailResults);
        }
        return new Response<>(result);
    }

    @Override
    public Response<QueryPagingListResponseBody> findLastestList1(QueryServiceDiscoveryListRequestBody body) {
        QueryPagingListResponseBody result;

        //查询区域列表
        List<Region> regionList = this.regionConditionImpl.findList(null);
        if (null == regionList || regionList.isEmpty()) {
            //result = PagingUtil.createResult(body, 0);
            return new Response<>(SystemErrCode.ERRCODE_REGION_NOT_EXISTS, SystemErrCode.ERRMSG_REGION_NOT_EXISTS);
        }

        //查询最近的搜索
        String projectName = body.getProject();
        String clusterName = body.getCluster();
        String serviceName = body.getService();
        SearchCondition searchCondition = this.lastestSearchServiceImpl.find(projectName, clusterName, serviceName);
        projectName = searchCondition.getProject();
        clusterName = searchCondition.getCluster();
        serviceName = searchCondition.getService();
        if (StringUtils.isBlank(projectName) || StringUtils.isBlank(clusterName) || StringUtils.isBlank(serviceName)) {
            result = PagingUtil.createResult(body, 0);
            return new Response<>(null);
        }

        //获取path
        String path = this.serviceDiscoveryCenter.getServicePath(projectName, clusterName);

        //查询服务下的版本
        int startIndex = PagingUtil.getStartIndex(body);
        int endIndex = PagingUtil.getEndIndex(body);
        int isPage = body.getIsPage();

        ServiceResult serviceResult = this.serviceDiscoveryCenter.findServicesByPaging(path, regionList, serviceName, isPage, startIndex, endIndex);
        List failureArea = serviceResult.getFailureArea();
        int failureCount = serviceResult.getFailureCount();
        List results = serviceResult.getResults();

        int totalCount = serviceResult.getTotalCount();
        result = PagingUtil.createResult(body, totalCount);

        //保存最近的搜索
        String condition = this.lastestSearchServiceImpl.saveLastestSearch(projectName, clusterName, serviceName);
        result.setCondition(condition);
        if (0 == totalCount) {
            return new Response<>(null);
        }

        //创建服务发现列表
        List<ServiceDetailNewResult> serviceDetailResults = results;
        if (null != serviceDetailResults && !serviceDetailResults.isEmpty()) {
            result.setList(serviceDetailResults);
        }
        Response<QueryPagingListResponseBody> queryPagingListResponseBodyResponse = new Response<>(result);
        if (failureCount!=0){
            queryPagingListResponseBodyResponse.setMessage("有"+failureCount+"个区域查询失败："+JSONArray.toJSONString(failureArea)+",请检查失败区域的url或者网络是否正常！");
            queryPagingListResponseBodyResponse.setCode(failureCount);
        }
        return queryPagingListResponseBodyResponse;
    }

    @Override
    public Response<ServiceDiscoveryResponseBody> find(QueryServiceDiscoveryDetailRequestBody body) {
        //通过区域名称查询区域信息
        String regionName = body.getRegion();
        Region region = this.regionConditionImpl.findByName(regionName);
        if (null == region) {
            return new Response<>(SystemErrCode.ERRCODE_REGION_NOT_EXISTS, SystemErrCode.ERRMSG_REGION_NOT_EXISTS);
        }

        //获取path
        String project = body.getProject();
        String cluster = body.getCluster();
        String service = body.getService();
        String apiVersion = body.getApiVersion();
        String path = this.serviceDiscoveryCenter.getConfigNewPath(project, cluster, service, apiVersion);

        //服务发现结果
        ServiceDiscoveryResponseBody result = new ServiceDiscoveryResponseBody();

        //版本名称
        result.setApiVersion(apiVersion);

        //查询服务配置信息:从缓存中查询
        ServiceConfRuleResult cacheCenterServiceConfResult = this.serviceDiscoveryCenter.findServiceConf(path, region);

        //拼装返回结果
        if (null != cacheCenterServiceConfResult && null != cacheCenterServiceConfResult.getUser()) {
            if (null != cacheCenterServiceConfResult.getUser() && !cacheCenterServiceConfResult.getUser().isEmpty()) {
                String loadBalance = cacheCenterServiceConfResult.getUser().getString("loadbalance");
                result.setLoadbalance(loadBalance);

                //去掉中间的负载均衡参数，留下自定义参数
                String params = cacheCenterServiceConfResult.getUser().fluentRemove("loadbalance").toJSONString();
                if ("{}".equals(params)){
                    params = null;
                }
                result.setParams(params);
            }
        }

        List<RouteRule> routeRule = cacheCenterServiceConfResult.getSdk();
        if (routeRule==null){
            routeRule = new ArrayList<>();
        }
        result.setRouteRules(routeRule);
        result.setRegion(cacheCenterServiceConfResult.getName());
        logger.info("----------------------------");
        logger.info(JSONArray.toJSONString(result));
        logger.info("----------------------------");
        return new Response<>(result);
    }

    @Override
    public Response<String> edit(EditServiceDiscoveryRequestBody body) {
        //通过区域名称查询区域信息
        String regionName = body.getRegion();
        Region region = this.regionConditionImpl.findByName(regionName);
        if (null == region) {
            return new Response<>(SystemErrCode.ERRCODE_REGION_NOT_EXISTS, SystemErrCode.ERRMSG_REGION_NOT_EXISTS);
        }

        //获取path
        String project = body.getProject();
        String cluster = body.getCluster();
        String service = body.getService();
        String apiVersion = body.getApiVersion();
        String path = this.serviceDiscoveryCenter.getConfigNewPath(project, cluster, service, apiVersion);
        //获取配置信息
        byte[] bt = {};

        //拼接参数
        String loadbalance = body.getLoadbalance();
        JSONObject User = new JSONObject();
        User.put("loadbalance", loadbalance);

        //校验自定义规则
        if (null != body.getParams() && !body.getParams().isEmpty()) {
            List<ServiceParam> params = body.getParams();
            for (ServiceParam serviceParam : params) {
                String key = serviceParam.getKey();
                String val = serviceParam.getVal();
                try {
                    User.put(new String(key.getBytes(), "UTF-8"), new String(val.getBytes(), "UTF-8"));
                } catch (UnsupportedEncodingException e) {
                    GlobalExceptionUtil.log(e);
                }
            }
            logger.info(User.toString());
        }

        //检查路由规则是不是每个提供者和消费之只能在一个规则中出现
        List<RouteRule> routeRules1 = body.getRouteRules();
        List<String> allProvider = new ArrayList<>();
        List<String> allConsumer = new ArrayList<>();
        for (RouteRule routeRul: routeRules1) {
            allProvider.addAll(routeRul.getProvider());
            allConsumer.addAll(routeRul.getConsumer());
        }
        int providerSize = allProvider.size();
        long providerCount = allProvider.stream().distinct().count();
        boolean isRepeat = providerCount<providerSize;
        if (isRepeat) {
            return new Response<>(SystemErrCode.ERRCODE_SERVICE_DISCOVERY_PARAMS_REPEAT, SystemErrCode.ERRMSG_SERVICE_DISCOVERY_ROULE_PROVIDER_REPEAT);
        }
        int consumerSize = allConsumer.size();
        long consumerCount = allConsumer.stream().distinct().count();
        boolean isRepeat1 = consumerCount<consumerSize;
        if (isRepeat1) {
            return new Response<>(SystemErrCode.ERRCODE_SERVICE_DISCOVERY_PARAMS_REPEAT, SystemErrCode.ERRMSG_SERVICE_DISCOVERY_ROULE_CONSUMER_REPEAT);
        }

        //校验自定义路由规则
        List<RouteRule> routeRules = new ArrayList<>();
        if (null != body.getRouteRules() && !body.getRouteRules().isEmpty()) {
            for (RouteRule routeRule : body.getRouteRules()) {
                routeRule.setId(SnowflakeIdWorker.getId());
                try {
                    byte[] bs = routeRule.getName().getBytes();
                    routeRule.setName(new String(bs, "UTF-8"));
                    logger.info("nameFormat" + "---------" + routeRule.getName());
                    logger.info("nameFormatBytes" + "---------" + Arrays.toString(routeRule.getName().getBytes()));
                } catch ( UnsupportedEncodingException e) {
                    GlobalExceptionUtil.log(e);
                }
                if (routeRule.getOnly().equals("true")) {
                    routeRule.setOnly("Y");
                } else {
                    routeRule.setOnly("N");
                }
                routeRules.add(routeRule);
            }

        }
        JSONArray routes = JSONArray.parseArray(JSON.toJSONString(routeRules));

        //通过一对一推送
        PushResult cacheCenterPushResult = this.serviceDiscoveryCenter.pushNewByOneToOne(path, bt, region, routes, User, false);

        //新增服务发现推送历史
        this.serviceDiscoveryPushHistoryConditionImpl.add(project, cluster, service, apiVersion, cacheCenterPushResult);
        return new Response<>(null);
    }


    @Override
    public Response<QueryPagingListResponseBody> provider(QueryServiceDiscoveryDetailRequestBody body) {
        QueryPagingListResponseBody result;
        //通过区域名称查询区域信息
        String regionName = body.getRegion();
        Region region = this.regionConditionImpl.findByName(regionName);
        if (null == region) {
            return new Response<>(SystemErrCode.ERRCODE_REGION_NOT_EXISTS, SystemErrCode.ERRMSG_REGION_NOT_EXISTS);
        }

        //获取服务生产者path
        String project = body.getProject();
        String cluster = body.getCluster();
        String service = body.getService();
        String apiVersion = body.getApiVersion();
        String path = this.serviceDiscoveryCenter.getServiceProviderPath(project, cluster, service, apiVersion);

        //查询服务提供端
        int startIndex = PagingUtil.getStartIndex(body);
        int endIndex = PagingUtil.getEndIndex(body);
        int isPage = body.getIsPage();
        ServiceResult serviceResult = this.serviceDiscoveryCenter.findProviderConsumersByPaging(path, region, true, startIndex, endIndex, isPage);
        int totalCount = serviceResult.getTotalCount();
        result = PagingUtil.createResult(body, totalCount);

        //创建提供端
        List<ServiceProviderInstanceConf> serviceProviderResults = serviceResult.getResults();
        if (null != serviceProviderResults && !serviceProviderResults.isEmpty()) {
            List<QueryServiceDiscoveryProviderResponseBody> list = this.createProvider(serviceProviderResults);
            result.setList(list);
        }
        return new Response<>(result);
    }

    @Override
    public Response<String> editProvider(EditServiceDiscoveryProviderRequestBody body) {
        //通过区域名称查询区域信息
        String regionName = body.getRegion();
        Region region = this.regionConditionImpl.findByName(regionName);
        if (null == region) {
            return new Response<>(SystemErrCode.ERRCODE_REGION_NOT_EXISTS, SystemErrCode.ERRMSG_REGION_NOT_EXISTS);
        }

        //获取服务生产者path
        String project = body.getProject();
        String cluster = body.getCluster();
        String service = body.getService();
        String apiVersion = body.getApiVersion();
        String addr = body.getAddr();
        String path = this.serviceDiscoveryCenter.getServiceProviderAddrPath(project, cluster, service, apiVersion, addr);

        //获取提供者配置信息
        byte[] bt = this.getProviderConf(body);
        //通过一对一推送
        PushResult cacheCenterPushResult = this.serviceDiscoveryCenter.pushByOneToOne(path, bt, region, true);

        //同步提供方
        String newPath = this.serviceDiscoveryCenter.getServiceProviderPath(project, cluster, service, apiVersion);
        this.serviceDiscoveryCenter.syncProviders(newPath, region, true);

        //新增服务发现推送历史
        this.serviceDiscoveryPushHistoryConditionImpl.add(project, cluster, service, apiVersion, cacheCenterPushResult);
        return new Response<>(null);
    }

    @Override
    public Response<QueryPagingListResponseBody> consumer(QueryServiceDiscoveryDetailRequestBody body) {
        QueryPagingListResponseBody result;
        //通过区域名称查询区域信息
        String regionName = body.getRegion();
        Region region = this.regionConditionImpl.findByName(regionName);
        if (null == region) {
            return new Response<>(SystemErrCode.ERRCODE_REGION_NOT_EXISTS, SystemErrCode.ERRMSG_REGION_NOT_EXISTS);
        }

        //获取服务消费者path
        String project = body.getProject();
        String cluster = body.getCluster();
        String service = body.getService();
        String apiVersion = body.getApiVersion();
        String path = this.serviceDiscoveryCenter.getServiceConsumerPath(project, cluster, service, apiVersion);

        //查询服务消费端
        int startIndex = PagingUtil.getStartIndex(body);
        int endIndex = PagingUtil.getEndIndex(body);
        int isPage = body.getIsPage();
        ServiceResult serviceResult = this.serviceDiscoveryCenter.findProviderConsumersByPaging(path, region, false, startIndex, endIndex, isPage);
        int totalCount = serviceResult.getTotalCount();
        result = PagingUtil.createResult(body, totalCount);

        //创建消费端
        List<ServiceProviderConsumerResult> consumerAddrs = serviceResult.getResults();
        if (null != consumerAddrs && !consumerAddrs.isEmpty()) {
            List<QueryServiceDiscoveryConsumerResponseBody> consumerList = this.createConsumer(consumerAddrs);
            result.setList(consumerList);
        }
        return new Response<>(result);
    }

    /**
     * 获取配置信息
     *
     * @param body
     * @return
     */
    private byte[] getConf(EditServiceDiscoveryRequestBody body) {
        String loadbalance = body.getLoadbalance();
        List<ServiceParam> params = body.getParams();
        ServiceConfResult content = new ServiceConfResult();
        content.setLb_mode(loadbalance);
        byte[] bt = new byte[0];
        try {
            bt = JSONObject.toJSONString(content).getBytes("UTF-8");
        } catch (UnsupportedEncodingException e) {
            GlobalExceptionUtil.log(e);
        }
        return bt;
    }

    /**
     * 获取配置规则信息
     *
     * @param body
     * @return
     */
    private byte[] getConfRule(EditServiceDiscoveryRequestBody body) {
        ServiceConfRuleResult confRuleResult = new ServiceConfRuleResult();
        String loadbalance = body.getLoadbalance();
        JSONObject User = new JSONObject();
        User.put("loadbalance", loadbalance);

        //校验自定义规则
        if (null != body.getParams() && !body.getParams().isEmpty()) {
            List<ServiceParam> params = body.getParams();
            for (ServiceParam serviceParam : params) {
                String key = serviceParam.getKey();
                String val = serviceParam.getVal();
                User.put(key, val);
            }
            logger.info(User.toString());
        }

        //校验自定义路由规则
        if (null != body.getRouteRules() && !body.getRouteRules().isEmpty()) {
            List<RouteRule> routeRules = new ArrayList<>();
            for (RouteRule routeRule : body.getRouteRules()) {
                routeRule.setId(SnowflakeIdWorker.getId());
                routeRules.add(routeRule);
            }
            confRuleResult.setSdk(routeRules);
        }
        confRuleResult.setUser(User);
        logger.info(confRuleResult.toString());

        byte[] bt = new byte[0];
        try {
            bt = StringEscapeUtils.unescapeJson(JSON.toJSONString(confRuleResult)).getBytes("UTF-8");
        } catch (UnsupportedEncodingException e) {
            GlobalExceptionUtil.log(e);
        }
        return bt;
    }

    /**
     * 获取提供者配置信息
     *
     * @param body
     * @return
     */
    private byte[] getProviderConf(EditServiceDiscoveryProviderRequestBody body) {
        int weight = body.getWeight();
        boolean isvalid = body.isValid();
        List<ServiceParam> params = body.getParams();
        //组装数据
        Map<String, Object> user = new HashMap<>();
        Map<String, Object> sdk = new HashMap<>();
        user.put("weight", weight);
        for (ServiceParam param: params) {
            user.put(param.getKey(), param.getVal());
        }
        sdk.put("is_valid", isvalid);
        ServiceProviderInstanceConf serviceProviderInstanceConf = new ServiceProviderInstanceConf();
        serviceProviderInstanceConf.setUser(user);
        serviceProviderInstanceConf.setSdk(sdk);
        byte[] bt = new byte[0];
        try {
            bt = JSONObject.toJSONString(serviceProviderInstanceConf).getBytes("UTF-8");
        } catch (UnsupportedEncodingException e) {
            GlobalExceptionUtil.log(e);
        }
        return bt;
    }

    /**
     * 创建负载均衡列表
     *
     * @param loadBalanceList
     * @param cacheCenterServiceConfResult
     * @return
     */
    private List<LoadBalanceDetail> createLoadBalanceList(List<LoadBalance> loadBalanceList, ServiceConfRuleResult cacheCenterServiceConfResult) {
        List<LoadBalanceDetail> list = new ArrayList<>();
        String lbMode = cacheCenterServiceConfResult.getUser().getString("loadbalance");
        boolean isSelect = false;
        for (LoadBalance loadBalance : loadBalanceList) {
            String abbr = loadBalance.getAbbr();
            String name = loadBalance.getName();
            LoadBalanceDetail loadBalanceDetail = new LoadBalanceDetail();
            loadBalanceDetail.setName(name);
            loadBalanceDetail.setAbbr(abbr);
            if (isSelect) {
                loadBalanceDetail.setSelected(false);
            } else {
                if (StringUtils.isEmpty(lbMode)) {
                    isSelect = true;
                    loadBalanceDetail.setSelected(true);
                } else {
                    if (abbr.toLowerCase().equals(lbMode.toLowerCase())) {
                        isSelect = true;
                        loadBalanceDetail.setSelected(true);
                    } else {
                        loadBalanceDetail.setSelected(false);
                    }
                }
            }
            list.add(loadBalanceDetail);
        }
        return list;
    }

    /**
     * 创建服务发现列表
     *
     * @param cacheCenterServiceResults
     * @return
     */
    private List<QueryServiceDiscoveryResponseBody> createServiceDiscoveryList(List<ServiceDetailResult> cacheCenterServiceResults) {
        List<QueryServiceDiscoveryResponseBody> serviceDiscoveryList = new ArrayList<>();
        QueryServiceDiscoveryResponseBody serviceDiscovery;
        for (ServiceDetailResult cacheCenterServiceResult : cacheCenterServiceResults) {
            serviceDiscovery = new QueryServiceDiscoveryResponseBody();
            serviceDiscovery.setRegion(cacheCenterServiceResult.getName());
            serviceDiscovery.setApiVersion(cacheCenterServiceResult.getApiVersion());
            serviceDiscoveryList.add(serviceDiscovery);
        }
        return serviceDiscoveryList;
    }

    @Override
    public Response<String> feedback(ServiceDiscoveryFeedBackRequestBody body) {
        String pushId = body.getPushId();
        ServiceDiscoveryPushHistory serviceDiscoveryPushHistory = this.serviceDiscoveryPushHistoryConditionImpl.findById(pushId);
        if (null == serviceDiscoveryPushHistory) {
            //不存在该轨迹
            return new Response<>(SystemErrCode.ERRCODE_TRACK_NOT_EXISTS, SystemErrCode.ERRMSG_TRACK_NOT_EXISTS);
        }

        //新增服务发现推送反馈
        this.serviceDiscoveryPushFeedbackConditionImpl.add(body);
        return new Response<>(null);
    }

    @Override
    public Response<List<LoadBalance>> findBalanceList() {
        //查询负载均衡列表
        List<LoadBalance> loadBalanceList = this.loadbalanceConditionImpl.findList(null);
        if (null == loadBalanceList || loadBalanceList.isEmpty()) {
            return new Response<>(null);
        }
        return new Response<>(loadBalanceList);
    }

    /**
     * 创建提供端
     *
     * @param  cacheCenterServiceProviderResults
     * @return
     */
    private List<QueryServiceDiscoveryProviderResponseBody> createProvider(List<ServiceProviderInstanceConf> cacheCenterServiceProviderResults) {
        List<QueryServiceDiscoveryProviderResponseBody> results = new ArrayList<>();
        for (ServiceProviderInstanceConf x: cacheCenterServiceProviderResults) {
            QueryServiceDiscoveryProviderResponseBody result = new QueryServiceDiscoveryProviderResponseBody();
            //放入地址
            result.setAddr(x.getAddr());
            Map<String, Object> user = x.getUser();
            Integer weight = (Integer)user.get("weight");
            //放入权值
            result.setWeight(weight.intValue());
            user.remove("weight");
            Map<String, Object> sdk = x.getSdk();
            //放入是否启用
            Boolean is_valid = (Boolean)sdk.get("is_valid");
            result.setValid(is_valid.booleanValue());
            //放入自定义K_V
            List<ServiceParam> user1 = new ArrayList<>();
            Iterator<Map.Entry<String, Object>> iter = user.entrySet().iterator();
            while (iter.hasNext()) {
                Map.Entry entry = (Map.Entry) iter.next();
                user1.add(new ServiceParam((String)entry.getKey(), (String) entry.getValue()));
            }
            result.setUser(user1);
            results.add(result);
        }
        return results;
    }

    /**
     *  创建消费端
     *
     * @param cacheCenterServiceProviderConsumerResults
     * @return
     */
    private List<QueryServiceDiscoveryConsumerResponseBody> createConsumer(List<ServiceProviderConsumerResult> cacheCenterServiceProviderConsumerResults) {
        List<QueryServiceDiscoveryConsumerResponseBody> results = new ArrayList<>();
        cacheCenterServiceProviderConsumerResults.forEach(x -> {
            QueryServiceDiscoveryConsumerResponseBody result = new QueryServiceDiscoveryConsumerResponseBody();
            result.setAddr(x.getAddr());
            results.add(result);
        });
        return results;
    }
}
