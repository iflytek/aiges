package com.iflytek.ccr.polaris.cynosure.request.servicediscovery;

import java.io.Serializable;
import java.util.List;

/**
 * Created by DELL-5490 on 2018/7/16.
 */

public class RouteRule implements Serializable {
    private static final long serialVersionUID = -4440989051339165137L;

    //路由规则id
    private String id;

    //路由规则名称参数
    private String name;

    //服务消费者
    private List<String> consumer;

    //服务提供者
    private List<String> provider;

    //仅供上述订阅者获取，Y是，N否
    private String only;

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public List<String> getConsumer() {
        return consumer;
    }

    public void setConsumer(List<String> consumer) {
        this.consumer = consumer;
    }

    public List<String> getProvider() {
        return provider;
    }

    public void setProvider(List<String> provider) {
        this.provider = provider;
    }

    public String getOnly() {
        return only;
    }

    public void setOnly(String only) {
        this.only = only;
    }

    @Override
    public String toString() {
        return "RouteRule{" +
                "routeRuleId='" + id + '\'' +
                ", name='" + name + '\'' +
                ", consumer='" + consumer + '\'' +
                ", provider='" + provider + '\'' +
                ", only='" + only + '\'' +
                '}';
    }
}

