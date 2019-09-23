package com.iflytek.ccr.polaris.cynosure.companionservice.domain.response;

/**
 * 服务数据-响应
 *
 * @author sctang2
 * @create 2017-12-07 15:03
 **/
public class ServiceDataNewResponse {
    private Sdk  sdk;
    private User user;

    public void setSdk(Sdk sdk) {
        this.sdk = sdk;
    }

    public Sdk getSdk() {
        return sdk;
    }

    public void setUser(User user) {
        this.user = user;
    }

    public User getUser() {
        return user;
    }
}
