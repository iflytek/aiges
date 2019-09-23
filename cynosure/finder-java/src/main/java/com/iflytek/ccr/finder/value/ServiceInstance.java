package com.iflytek.ccr.finder.value;

import java.io.Serializable;

/**
 * 服务实例对象
 */
public class ServiceInstance implements Serializable {

    /**
     * 服务实例地址，eg：1.2.3.4:8080
     */
    private String addr;

    /**
     * json配置信息，这个字段供服务订阅者使用
     */
    private String jsonConfig;

    public String getJsonConfig() {
        return jsonConfig;
    }

    public void setJsonConfig(String jsonConfig) {
        this.jsonConfig = jsonConfig;
    }

    public String getAddr() {
        return addr;
    }

    public void setAddr(String addr) {
        this.addr = addr;
    }


    @Override
    public String toString() {
        return "ServiceInstance{" +
                "addr='" + addr + '\'' +
                ", jsonConfig='" + jsonConfig + '\'' +
                '}';
    }
}
