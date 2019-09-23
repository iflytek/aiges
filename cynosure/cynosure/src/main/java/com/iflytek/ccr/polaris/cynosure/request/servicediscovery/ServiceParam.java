package com.iflytek.ccr.polaris.cynosure.request.servicediscovery;


import java.io.Serializable;

/**
 * 自定义规则设置
 * Created by DELL-5490 on 2018/5/17.
 */
public class ServiceParam implements Serializable {
    private static final long serialVersionUID = -4440989051339265137L;

    //自定义参数
    private String key;

    //自定义参数
    private String val;

    public ServiceParam() {
    }

    public ServiceParam(String key, String val) {
        this.key = key;
        this.val = val;
    }

    public String getKey() {
        return key;
    }

    public void setKey(String key) {
        this.key = key;
    }

    public String getVal() {
        return val;
    }

    public void setVal(String val) {
        this.val = val;
    }

    @Override
    public String toString() {
        return "ServiceParam{" +
                "key='" + key + '\'' +
                ", val='" + val + '\'' +
                '}';
    }
}
