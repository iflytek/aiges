package com.iflytek.ccr.polaris.cynosure.service.impl;

import com.alibaba.fastjson.JSON;
import com.alibaba.fastjson.JSONArray;
import com.iflytek.ccr.polaris.cynosure.companionservice.ServiceDiscoveryCenter;
import com.iflytek.ccr.polaris.cynosure.companionservice.domain.PushDetailResult;
import com.iflytek.ccr.polaris.cynosure.companionservice.domain.PushResult;
import com.iflytek.ccr.polaris.cynosure.companionservice.domain.Result;
import com.iflytek.ccr.polaris.cynosure.companionservice.domain.response.ServiceResponse;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IRegionCondition;
import com.iflytek.ccr.polaris.cynosure.dbcondition.IServiceDiscoveryPushHistoryCondition;
import com.iflytek.ccr.polaris.cynosure.domain.Region;
import com.iflytek.ccr.polaris.cynosure.errorcode.SystemErrCode;
import com.iflytek.ccr.polaris.cynosure.exception.GlobalExceptionUtil;
import com.iflytek.ccr.polaris.cynosure.request.servicediscovery.EditServiceDiscoveryRequestBody;
import com.iflytek.ccr.polaris.cynosure.request.servicediscovery.RouteRule;
import com.iflytek.ccr.polaris.cynosure.request.servicediscovery.ServiceParam;
import com.iflytek.ccr.polaris.cynosure.response.Response;
import com.iflytek.ccr.polaris.cynosure.service.IMyServiceDiscovery;
import com.iflytek.ccr.polaris.cynosure.util.SnowflakeIdWorker;
import org.apache.http.HttpEntity;
import org.apache.http.NameValuePair;
import org.apache.http.client.fluent.Request;
import org.apache.http.message.BasicNameValuePair;
import org.apache.http.util.EntityUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.nio.charset.Charset;
import java.util.*;

@Service
public class MyServiceDiscovery implements IMyServiceDiscovery {
    public static final String PUSH_SERVICE_CLUSTER_CONFIG_URL = "/service/push_cluster_config";
    @Autowired
    private IRegionCondition regionConditionImpl;
    @Autowired
    private ServiceDiscoveryCenter serviceDiscoveryCenter;

    @Autowired
    private IServiceDiscoveryPushHistoryCondition serviceDiscoveryPushHistoryConditionImpl;

    @Override
    public Response<String> edit1(EditServiceDiscoveryRequestBody body) {

        //检查路由规则是不是每个提供者和消费之只能在一个规则中出现
        List<RouteRule> routeRules1 = body.getRouteRules();
        if (routeRules1!=null&&!routeRules1.isEmpty()) {
            List<String> allProvider = new ArrayList<>();
            List<String> allConsumer = new ArrayList<>();
            for (RouteRule routeRul : routeRules1) {
                allProvider.addAll(routeRul.getProvider());
                allConsumer.addAll(routeRul.getConsumer());
            }
            int providerSize = allProvider.size();
            long providerCount = allProvider.stream().distinct().count();
            boolean isRepeat = providerCount < providerSize;
            if (isRepeat) {
                return new Response<>(SystemErrCode.ERRCODE_SERVICE_DISCOVERY_PARAMS_REPEAT, SystemErrCode.ERRMSG_SERVICE_DISCOVERY_ROULE_PROVIDER_REPEAT);
            }
            int consumerSize = allConsumer.size();
            long consumerCount = allConsumer.stream().distinct().count();
            boolean isRepeat1 = consumerCount < consumerSize;
            if (isRepeat1) {
                return new Response<>(SystemErrCode.ERRCODE_SERVICE_DISCOVERY_PARAMS_REPEAT, SystemErrCode.ERRMSG_SERVICE_DISCOVERY_ROULE_CONSUMER_REPEAT);
            }
        }

        //获得推送区域
        String regionName = body.getRegion();
        Region region = regionConditionImpl.findByName(regionName);
        String path = this.serviceDiscoveryCenter.getConfigNewPath(body.getProject(), body.getCluster(), body.getService(), body.getApiVersion());
        //构造user参数
        Map<String, Object> user = new HashMap<>();
        user.put("loadbalance", body.getLoadbalance());
        List<ServiceParam> params = body.getParams();
        if (params!=null&&!params.isEmpty())
        for (ServiceParam param: params) {
            user.put(param.getKey(), param.getVal());
        }

        //校验自定义路由规则
        List<RouteRule> routeRules = new ArrayList<>();
        if (null != body.getRouteRules() && !body.getRouteRules().isEmpty()) {
            for (RouteRule routeRule : body.getRouteRules()) {
                routeRule.setId(SnowflakeIdWorker.getId());
                routeRule.setName(routeRule.getName());
                if (routeRule.getOnly().equals("true")) {
                    routeRule.setOnly("Y");
                } else {
                    routeRule.setOnly("N");
                }
                routeRules.add(routeRule);
            }

        }

        //推送
        String pushId = SnowflakeIdWorker.getId();
        String url = region.getPushUrl() + PUSH_SERVICE_CLUSTER_CONFIG_URL;
        List <NameValuePair> formParams = new ArrayList<>();
        formParams.add(new BasicNameValuePair("pushId", pushId));
        formParams.add(new BasicNameValuePair("path", path));
        formParams.add(new BasicNameValuePair("user", JSONArray.toJSONString(user)));
        formParams.add(new BasicNameValuePair("sdk", JSONArray.toJSONString(routeRules)));
        String response = null;
        try{
            HttpEntity result = Request.Post(url) .socketTimeout(5000).connectTimeout(5000)
                     .bodyForm(formParams, Charset.forName("utf-8"))
                    .execute().returnResponse().getEntity();
            response = EntityUtils.toString(result, "utf-8");
        }catch (Exception e){
            e.printStackTrace();
        }

        List<Result> cacheCenterResults = new ArrayList<>();
        Result result1 = new Result();
        result1.setName(region.getName());
        result1.setResult(response);
        cacheCenterResults.add(result1);


        //解析推送参数
        List<PushDetailResult> cacheCenterPushDetailResults = this.parsePush(cacheCenterResults);
        PushResult pushResult = new PushResult();
        pushResult.setPushId(pushId);
        pushResult.setResult(JSONArray.toJSONString(cacheCenterPushDetailResults));

        //新增服务发现推送历史
        this.serviceDiscoveryPushHistoryConditionImpl.add(body.getProject(), body.getCluster(), body.getService(), body.getApiVersion(), pushResult);
        return new Response<>(null);
    }

    /**
     *  解析推送参数
     *
     * @param cacheCenterResults
     * @return
     */
    protected List<PushDetailResult> parsePush(List<Result> cacheCenterResults) {
        List<PushDetailResult> results = new ArrayList<>();
        cacheCenterResults.forEach(x -> {
            PushDetailResult result = new PushDetailResult();
            String name = x.getName();
            result.setName(name);
            try {
                ServiceResponse cacheServiceResponse = JSON.parseObject(x.getResult(), ServiceResponse.class);
                if (null == cacheServiceResponse) {
                    result.setSuccessed(-1);
                } else {
                    int ret = cacheServiceResponse.getRet();
                    if (0 == ret) {
                        result.setSuccessed(1);
                    } else {
                        result.setSuccessed(-1);
                    }
                }
            } catch (Exception ex) {
                GlobalExceptionUtil.log(ex);
                result.setSuccessed(-1);
            }
            results.add(result);
        });
        return results;
    }


}
