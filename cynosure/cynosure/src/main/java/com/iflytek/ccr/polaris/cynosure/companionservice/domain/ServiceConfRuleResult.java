package com.iflytek.ccr.polaris.cynosure.companionservice.domain;

import com.alibaba.fastjson.JSONObject;
import com.iflytek.ccr.polaris.cynosure.request.servicediscovery.RouteRule;

import java.util.List;

/**
 * 服务配置结果
 *
 * @author sctang2
 * @create 2017-12-11 15:30
 **/
public class ServiceConfRuleResult {

    //区域名称
    private String name;

    //sdk
    private List<RouteRule> sdk;

    //user
    private JSONObject user;

    public List<RouteRule> getSdk() {
        return sdk;
    }

    public void setSdk(List<RouteRule> sdk) {
        this.sdk = sdk;
    }

    public JSONObject getUser() {
        return user;
    }

    public void setUser(JSONObject user) {
        this.user = user;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }
}
