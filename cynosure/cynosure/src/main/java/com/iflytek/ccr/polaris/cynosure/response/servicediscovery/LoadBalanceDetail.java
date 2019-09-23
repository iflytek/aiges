package com.iflytek.ccr.polaris.cynosure.response.servicediscovery;

import java.io.Serializable;

/**
 * 负载均衡明细-响应
 *
 * @author sctang2
 * @create 2017-12-07 11:15
 **/
public class LoadBalanceDetail implements Serializable {
    private static final long serialVersionUID = -668128974591572570L;

    //负载均衡名称
    private String name;

    //负载均衡缩写
    private String abbr;

    //是否选中
    private boolean selected;

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }

    public String getAbbr() {
        return abbr;
    }

    public void setAbbr(String abbr) {
        this.abbr = abbr;
    }

    public boolean isSelected() {
        return selected;
    }

    public void setSelected(boolean selected) {
        this.selected = selected;
    }
}
