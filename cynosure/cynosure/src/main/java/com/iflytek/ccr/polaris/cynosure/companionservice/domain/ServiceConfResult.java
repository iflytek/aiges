package com.iflytek.ccr.polaris.cynosure.companionservice.domain;

/**
 * 服务配置结果
 *
 * @author sctang2
 * @create 2017-12-11 15:30
 **/
public class ServiceConfResult {
    //负载均衡
    private String lb_mode;

    public String getLb_mode() {
        return lb_mode;
    }

    public void setLb_mode(String lb_mode) {
        this.lb_mode = lb_mode;
    }
}
