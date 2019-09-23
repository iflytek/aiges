package com.iflytek.ccr.finder.value;

import org.codehaus.jackson.annotate.JsonProperty;

import java.util.List;

/**
 * 服务发现路由对象
 */
public class ServiceRouteValue {

    /**
     * 路由规则id
     */
    @JsonProperty("id")
    private String id;

    /**
     *  规则对应的服务消费者
     */
    @JsonProperty("consumer")
    private List<String> consumer;

    /**
     *  规则对应的服务提供者
     */
    @JsonProperty("provider")
    private List<String> provider;

    /**
     *  该条规则对应的服务提供者是否只用于该规则,Y则对规则之外的服务消费者不可见
     */
    @JsonProperty("only")
    private String only;

    public List<String> getConsumer() {
        return consumer;
    }

    public void setConsumer(List<String> consumer) {
        this.consumer = consumer;
    }

    public String getOnly() {
        return only;
    }

    public void setOnly(String only) {
        this.only = only;
    }

    public List<String> getProvider() {
        return provider;
    }

    public void setProvider(List<String> provider) {
        this.provider = provider;
    }

    public String getId() {
        return id;
    }

    public void setId(String id) {
        this.id = id;
    }

    @Override
    public String toString() {
        return "ServiceRouteValue{" +
                "id='" + id + '\'' +
                ", consumer=" + consumer +
                ", provider=" + provider +
                ", only='" + only + '\'' +
                '}';
    }
}
