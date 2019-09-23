package com.iflytek.ccr.polaris.cynosure.domain;

import java.util.Map;

/**
 * 服务提供者配置实例
 * create by ygli3
 */
public class ServiceProviderInstanceConf {

    private String addr;

    private Map<String, Object> user;

    private Map<String, Object> sdk;

    public String getAddr() {
        return addr;
    }

    public void setAddr(String addr) {
        this.addr = addr;
    }

    public Map<String, Object> getUser() {
        return user;
    }

    public void setUser(Map<String, Object> user) {
        this.user = user;
    }

    public Map<String, Object> getSdk() {
        return sdk;
    }

    public void setSdk(Map<String, Object> sdk) {
        this.sdk = sdk;
    }
}
