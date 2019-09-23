package com.iflytek.ccr.polaris.cynosure.domain;

/**
 * 负载均衡模型
 *
 * @author sctang2
 * @create 2017-12-07 11:04
 **/
public class LoadBalance {
    //名称
    private String name;

    //简称
    private String abbr;

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name == null ? null : name.trim();
    }

    public String getAbbr() {
        return abbr;
    }

    public void setAbbr(String abbr) {
        this.abbr = abbr == null ? null : abbr.trim();
    }
}
